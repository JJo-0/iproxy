package main

import (
	"context"
	"crypto/tls"
	"errors"
	"flag"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"

	"golang.org/x/net/proxy"
)

func TestConfig(t *testing.T) {
	c, err := parseConfig([]string{"-a", "192.168.1.2", "-b", "127.0.0.1", "-p", "8888", "-s", "0", "-d", "9000", "-l"}, io.Discard)
	if err != nil || c.http != 8888 || c.socks != 0 || c.pac != 9000 || c.bind != "127.0.0.1" || !c.location {
		t.Fatalf("flags not applied: %+v, %v", c, err)
	}
	for _, args := range [][]string{{"-p", "-1"}, {"-s", "65536"}, {"-p", "1080"}, {"-p", "0", "-s", "0"}, {"-a", "x'; DIRECT"}, {"-a", "0.0.0.0"}, {"-b", "wrong"}, {"unexpected"}} {
		if _, err := parseConfig(args, io.Discard); err == nil {
			t.Errorf("accepted %v", args)
		}
	}
	for _, arg := range []string{"-h", "-version"} {
		if _, err := parseConfig([]string{arg}, io.Discard); !errors.Is(err, flag.ErrHelp) {
			t.Errorf("%s: %v", arg, err)
		}
	}
}

func TestPAC(t *testing.T) {
	c := config{address: "172.20.10.1", http: 8080, socks: 1080}
	w := httptest.NewRecorder()
	pacHandler(c).ServeHTTP(w, httptest.NewRequest("GET", "/proxy", nil))
	if w.Header().Get("Content-Type") != "application/x-ns-proxy-autoconfig" || !strings.Contains(w.Body.String(), "PROXY 172.20.10.1:8080; SOCKS5 172.20.10.1:1080") || strings.Contains(w.Body.String(), "DIRECT") {
		t.Fatal(w)
	}
	c.direct = true
	w = httptest.NewRecorder()
	pacHandler(c).ServeHTTP(w, httptest.NewRequest("GET", "/proxy", nil))
	if !strings.Contains(w.Body.String(), "; DIRECT") {
		t.Fatal(w.Body.String())
	}
	w = httptest.NewRecorder()
	pacHandler(c).ServeHTTP(w, httptest.NewRequest("GET", "/missing", nil))
	if w.Code != 404 {
		t.Fatal(w.Code)
	}
}

func TestMissingLocation(t *testing.T) {
	if err := fetchLocation(context.Background(), t.TempDir()+"/missing"); err == nil {
		t.Fatal("expected an error")
	}
}

func freePorts(t *testing.T) []int {
	t.Helper()
	ports := []int{}
	for i := 0; i < 3; i++ {
		l, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			t.Fatal(err)
		}
		defer l.Close()
		ports = append(ports, l.Addr().(*net.TCPAddr).Port)
	}
	return ports
}

func TestProxyTraffic(t *testing.T) {
	origin := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { io.WriteString(w, "through proxy") }))
	defer origin.Close()
	secure := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { io.WriteString(w, "through proxy") }))
	defer secure.Close()
	ports := freePorts(t)
	// Advertised IP deliberately differs: all listeners must honor -b.
	c := config{address: "192.0.2.1", bind: "127.0.0.1", http: ports[0], socks: ports[1], pac: ports[2]}
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- run(ctx, c) }()
	defer func() {
		cancel()
		select {
		case err := <-done:
			if err != nil {
				t.Error(err)
			}
		case <-time.After(3 * time.Second):
			t.Error("server failed to stop")
		}
	}()
	deadline := time.Now().Add(5 * time.Second)
	for {
		conn, err := net.DialTimeout("tcp", endpoint(c.bind, c.pac), 100*time.Millisecond)
		if err == nil {
			conn.Close()
			break
		}
		if time.Now().After(deadline) {
			t.Fatal(err)
		}
		time.Sleep(10 * time.Millisecond)
	}
	u, _ := url.Parse("http://" + endpoint(c.bind, c.http))
	transport := &http.Transport{Proxy: http.ProxyURL(u), TLSClientConfig: &tls.Config{InsecureSkipVerify: true}} // Test TLS server only.
	defer transport.CloseIdleConnections()
	client := &http.Client{Transport: transport, Timeout: 5 * time.Second}
	check := func(client *http.Client, address string) {
		t.Helper()
		res, err := client.Get(address)
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()
		body, err := io.ReadAll(res.Body)
		if err != nil || string(body) != "through proxy" {
			t.Fatalf("body=%q err=%v", body, err)
		}
	}
	check(client, origin.URL)
	check(client, secure.URL)
	dialer, err := proxy.SOCKS5("tcp", endpoint(c.bind, c.socks), nil, &net.Dialer{Timeout: 3 * time.Second})
	if err != nil {
		t.Fatal(err)
	}
	socksTransport := &http.Transport{DialContext: func(ctx context.Context, network, address string) (net.Conn, error) {
		return dialer.(proxy.ContextDialer).DialContext(ctx, network, address)
	}}
	defer socksTransport.CloseIdleConnections()
	check(&http.Client{Transport: socksTransport, Timeout: 5 * time.Second}, origin.URL)
}

func TestOccupiedPort(t *testing.T) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()
	_, port, _ := net.SplitHostPort(l.Addr().String())
	p, _ := strconv.Atoi(port)
	err = run(context.Background(), config{bind: "127.0.0.1", address: "127.0.0.1", http: p})
	if err == nil {
		t.Fatal("expected bind failure")
	}
}
