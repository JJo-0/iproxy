# iProxy

[한국어](README.md) · **English** · [简体中文](README.zh-CN.md) · [日本語](README.ja.md)

A small HTTP / SOCKS5 proxy running in iSH on an iPhone. It forwards **TCP traffic from apps that support proxy settings** on Windows, Linux, Mac, iPad, and iPhone connected to the iPhone's Personal Hotspot.

> **iOS 15, 16, 17, and 18 are target versions; testing on physical devices is not yet complete for any of them.** This is not a VPN and does not forward traffic from every app.

## Install on your iPhone

No jailbreak or Mac/Xcode is required. Install [iSH Shell](https://apps.apple.com/app/ish-shell/id1436902243). iSH runs an x86 Linux environment, so use the **`linux-386` binary**, even on an iPhone. `darwin-arm64` is for Macs.

1. Open Settings → Personal Hotspot and enable **Allow Others to Join**. Connect your other device to this hotspot over Wi-Fi.
2. Check [Releases](https://github.com/JJo-0/iproxy/releases) for a published build. The download commands below work after a release containing the new binaries is published. Until then, download and extract the `iproxy-binaries` artifact from a successful run for the relevant branch in [Actions](https://github.com/JJo-0/iproxy/actions/workflows/ci.yml), then install it using the file-copy instructions below (GitHub sign-in required). You can also build from source.
3. Run these commands in iSH. Installation stops if downloading or checksum verification fails.

```sh
apk update && apk add ca-certificates curl
(
  set -eu
  mkdir -p /tmp/iproxy-download
  cd /tmp/iproxy-download
  base=https://github.com/JJo-0/iproxy/releases/latest/download
  curl -fL "$base/iproxy-linux-386" -o iproxy-linux-386
  curl -fL "$base/SHA256SUMS" -o SHA256SUMS
  awk '$2 == "iproxy-linux-386" { print }' SHA256SUMS > CHECKSUM
  test -s CHECKSUM
  sha256sum -c CHECKSUM
  install -m 755 iproxy-linux-386 /usr/local/bin/iproxy
)
```

4. Keep iSH open and start the proxy. Allow local network access if prompted.

```sh
iproxy
```

Default ports: HTTP **8080**, SOCKS5 **1080**, PAC **8090**. Look for the `listening` messages to confirm startup. Press **Ctrl+C** to stop. Start it again after restarting your phone or closing iSH.

The default phone address is `172.20.10.1`. Check the **router/default gateway** in the connected device's Wi-Fi details. If it differs, replace the address in this command and in the client settings below:

```sh
iproxy -a 172.20.10.1
```

`-a` sets the address advertised in PAC; `-b` sets the local listening address. The default `-b 0.0.0.0` listens on all IPv4 interfaces. There is no proxy authentication, so use your own password-protected hotspot with trusted devices.

## Configure connected devices

First open `http://172.20.10.1:8090/proxy` in the connected device's browser. The PAC code should appear or download. This checks connectivity to the phone, not successful internet forwarding.

| Device | Recommended settings |
| --- | --- |
| Windows 10/11 | Settings → Network & Internet → Proxy → Use setup script → enter the PAC URL below |
| macOS | System Settings → Network → connected Wi-Fi Details → Proxies → Automatic proxy configuration → enter PAC URL (older versions: System Preferences → Network → Advanced) |
| Linux | Desktop network settings → Proxy → Automatic → enter PAC URL. Support varies by desktop/app; use the commands below for CLI tools. |
| iPad / iPhone | Settings → Wi-Fi → ⓘ next to the iPhone hotspot → Configure Proxy under HTTP Proxy → Automatic → enter PAC URL → Save |

```text
http://172.20.10.1:8090/proxy
```

For apps that do not apply automatic configuration, set a manual HTTP proxy: **server `172.20.10.1`, port `8080`, authentication off**. HTTPS uses a CONNECT tunnel through the same HTTP proxy; no certificate installation is required. iOS/iPadOS Wi-Fi proxy settings apply only to apps that honor them.

Check from the command line (use `curl.exe` on Windows):

```sh
curl --proxy http://172.20.10.1:8080 https://example.com
curl --proxy socks5h://172.20.10.1:1080 https://example.com
```

For supported programs in a Linux/macOS shell:

```sh
export http_proxy=http://172.20.10.1:8080
export https_proxy=http://172.20.10.1:8080
export all_proxy=socks5h://172.20.10.1:1080
```

When finished, set the device's proxy settings back to **Off/None**. In the shell, run `unset http_proxy https_proxy all_proxy`. Leaving a proxy enabled after disconnecting the hotspot may make the internet appear unavailable.

## Scope and limitations

- **Wi-Fi hotspot connections** are the documented default. USB/Bluetooth depend on the OS and drivers and require separate testing.
- Supports HTTP, HTTPS CONNECT, and SOCKS5 TCP. UDP, QUIC/HTTP3, and all games, messaging apps, or VPNs are not guaranteed to work. An app may need a TCP fallback.
- Enter the PAC URL manually. WPAD automatic discovery is not provided.
- The default PAC tries HTTP, then SOCKS5, without `DIRECT`. Enable `-direct` to permit a direct connection when the proxy fails. This does not block apps that ignore PAC.
- iOS may suspend or terminate iSH in the background. `-l` experimentally reads `/dev/location` and does not guarantee continuous operation. Coordinates are neither stored nor printed. The proxy keeps running if the location device is unavailable.

Record device/version results in the [compatibility checklist](docs/compatibility.en.md). For problems, see [troubleshooting](docs/troubleshooting.en.md).

## Options

| Option | Default | Description |
| --- | --- | --- |
| `-a` | `172.20.10.1` | Phone IP advertised in PAC |
| `-b` | `0.0.0.0` | Local listening IP |
| `-p` | `8080` | HTTP/HTTPS CONNECT port; `0` disables |
| `-s` | `1080` | SOCKS5 TCP port; `0` disables |
| `-d` | `8090` | PAC `/proxy` port; `0` disables |
| `-direct` | Off | Allow direct fallback in PAC |
| `-l` | Off | Attempt to read the iSH location device |
| `-v` | Off | Verbose HTTP logs (may include sensitive URLs) |
| `-version` | — | Print build version |
| `-h` | — | Show help |

Unlike the original version, services start even without port options. Set unwanted ports to `0` to enable only selected services.

```sh
iproxy -p 8080 -s 0 -d 0
```

## Build and release

The required Go version is specified in `go.mod`. Build on Linux/macOS or WSL on Windows. Building inside iSH can be slow or fail because of compatibility issues.

```sh
git clone https://github.com/JJo-0/iproxy.git
cd iproxy
go test -race ./...
go vet ./...
./scripts/build.sh dev
```

Copy `dist/iproxy-linux-386` into the iSH folder in the iPhone Files app. In iSH, change to that directory and install:

```sh
install -m 755 ./iproxy-linux-386 /usr/local/bin/iproxy
iproxy
```

The build produces Linux 386 for iSH, Linux amd64/arm64, Mac amd64/arm64, Windows amd64/arm64, and a SHA-256 checksum list. Connected devices do not need iProxy installed; configure their proxy settings only. Pushing a `v*` tag triggers GitHub Actions to test, build, and create a **draft release**. Review the files and physical-device results before publishing.

## References

- [Official iSH repository: x86 Linux emulation](https://github.com/ish-app/ish)
- [Apple: Set up and connect to Personal Hotspot](https://support.apple.com/en-us/111785)
- [Apple: Change proxy settings on Mac](https://support.apple.com/guide/mac-help/change-proxy-settings-on-mac-mchlp2591/mac)

These sources explain platform settings; they do not replace physical-device testing of this project.
