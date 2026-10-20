package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/armon/go-socks5"
	"github.com/elazarl/goproxy"
)

var version = "dev"

type config struct {
	address, bind             string
	socks, http, pac          int
	location, verbose, direct bool
}

func parseConfig(args []string, out io.Writer) (config, error) {
	c := config{}
	f := flag.NewFlagSet("iproxy", flag.ContinueOnError)
	f.SetOutput(out)
	f.StringVar(&c.address, "a", "172.20.10.1", "Phone IP advertised in PAC (check the client's gateway)")
	f.StringVar(&c.bind, "b", "0.0.0.0", "Local IP to listen on")
	f.IntVar(&c.http, "p", 8080, "HTTP/HTTPS CONNECT port; 0 disables")
	f.IntVar(&c.socks, "s", 1080, "SOCKS5 TCP port; 0 disables")
	f.IntVar(&c.pac, "d", 8090, "PAC HTTP port (/proxy); 0 disables")
	f.BoolVar(&c.location, "l", false, "Read iSH /dev/location (experimental background aid)")
	f.BoolVar(&c.verbose, "v", false, "Verbose HTTP proxy logging (may contain URLs)")
	f.BoolVar(&c.direct, "direct", false, "Allow PAC clients to fall back to a direct connection")
	showVersion := f.Bool("version", false, "Print version")
	if err := f.Parse(args); err != nil {
		return c, err
	}
	if *showVersion {
		fmt.Fprintln(out, version)
		return c, flag.ErrHelp
	}
	if f.NArg() != 0 {
		return c, fmt.Errorf("unexpected arguments: %s", strings.Join(f.Args(), " "))
	}
	return c, c.validate()
}

func (c config) validate() error {
	if net.ParseIP(c.address) == nil || net.ParseIP(c.address).IsUnspecified() {
		return errors.New("-a must be a reachable phone IP address")
	}
	if net.ParseIP(c.bind) == nil {
		return errors.New("-b must be a local IP address")
	}
	used := map[int]bool{}
	for _, port := range []int{c.http, c.socks, c.pac} {
		if port < 0 || port > 65535 {
			return errors.New("ports must be between 0 and 65535")
		}
		if port == 0 {
			continue
		}
		if used[port] {
			return fmt.Errorf("port %d is used more than once", port)
		}
		used[port] = true
	}
	if c.http == 0 && c.socks == 0 {
		return errors.New("enable at least one proxy with -p or -s")
	}
	return nil
}

func endpoint(host string, port int) string { return net.JoinHostPort(host, strconv.Itoa(port)) }

func pacHandler(c config) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/proxy", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			w.Header().Set("Allow", "GET, HEAD")
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		routes := []string{}
		// HTTP CONNECT has broader system PAC support than SOCKS5.
		if c.http != 0 {
			routes = append(routes, "PROXY "+endpoint(c.address, c.http))
		}
		if c.socks != 0 {
			routes = append(routes, "SOCKS5 "+endpoint(c.address, c.socks))
		}
		if c.direct {
			routes = append(routes, "DIRECT")
		}
		w.Header().Set("Content-Type", "application/x-ns-proxy-autoconfig")
		w.Header().Set("Cache-Control", "no-store")
		fmt.Fprintf(w, "function FindProxyForURL(url, host) { return %q; }\n", strings.Join(routes, "; "))
	})
	return mux
}

func fetchLocation(ctx context.Context, path string) error {
	reader, err := os.Open(path)
	if err != nil {
		return err
	}
	defer reader.Close()
	done := make(chan struct{})
	defer close(done)
	go func() {
		select {
		case <-ctx.Done():
			reader.Close()
		case <-done:
		}
	}()
	// Discard coordinates: location access is only an optional background aid.
	_, err = io.Copy(io.Discard, reader)
	if ctx.Err() != nil {
		return nil
	}
	return err
}

func run(ctx context.Context, c config) error {
	type service struct {
		name  string
		port  int
		serve func(net.Listener) error
	}
	services := []service{}
	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = c.verbose
	// Use the phone's connection directly, even if its shell has proxy variables.
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.Proxy = nil
	proxy.Tr = transport
	defer transport.CloseIdleConnections()
	proxy.ConnectDialWithReq = func(req *http.Request, network, address string) (net.Conn, error) {
		return (&net.Dialer{Timeout: 30 * time.Second}).DialContext(req.Context(), network, address)
	}
	for _, spec := range []struct {
		name    string
		port    int
		handler http.Handler
	}{
		{"HTTP", c.http, proxy}, {"PAC", c.pac, pacHandler(c)},
	} {
		if spec.port == 0 {
			continue
		}
		server := &http.Server{Handler: spec.handler, ReadHeaderTimeout: 15 * time.Second}
		defer server.Close()
		services = append(services, service{spec.name, spec.port, server.Serve})
	}
	if c.socks != 0 {
		server, err := socks5.New(&socks5.Config{Dial: func(ctx context.Context, network, addr string) (net.Conn, error) {
			return (&net.Dialer{Timeout: 30 * time.Second}).DialContext(ctx, network, addr)
		}})
		if err != nil {
			return err
		}
		services = append(services, service{"SOCKS5", c.socks, server.Serve})
	}
	listeners := []net.Listener{}
	// Bind every port before starting any service, so startup cannot partially succeed.
	for _, s := range services {
		listener, err := net.Listen("tcp", endpoint(c.bind, s.port))
		if err != nil {
			return fmt.Errorf("%s: %w", s.name, err)
		}
		defer listener.Close()
		listeners = append(listeners, listener)
	}
	errs := make(chan error, len(services))
	for i, s := range services {
		log.Printf("%s listening on %s", s.name, listeners[i].Addr())
		go func(s service, listener net.Listener) { errs <- s.serve(listener) }(s, listeners[i])
	}
	if c.pac != 0 {
		log.Printf("PAC URL: http://%s/proxy", endpoint(c.address, c.pac))
	}
	if c.location {
		locationCtx, cancel := context.WithCancel(ctx)
		defer cancel()
		go func() {
			if err := fetchLocation(locationCtx, "/dev/location"); err != nil {
				log.Printf("location unavailable (proxy stays running): %v", err)
			}
		}()
	}
	select {
	case <-ctx.Done():
		return nil
	case err := <-errs:
		return err
	}
}

func main() {
	c, err := parseConfig(os.Args[1:], os.Stdout)
	if errors.Is(err, flag.ErrHelp) {
		return
	}
	if err != nil {
		log.Print(err)
		os.Exit(2)
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	if err := run(ctx, c); err != nil {
		log.Print(err)
		os.Exit(1)
	}
}
