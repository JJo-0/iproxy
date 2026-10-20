# Compatibility checklist

[한국어](compatibility.md) · **English** · [简体中文](compatibility.zh-CN.md) · [日本語](compatibility.ja.md) · [README](../README.en.md)

`Untested` does not mean failed. It means reproducible results from physical devices are not yet available. Linux automated tests and cross-compilation alone do not establish iOS compatibility.

| iPhone OS | iSH installation/startup | HTTP / HTTPS | SOCKS5 | PAC | Screen lock/background |
| --- | --- | --- | --- | --- | --- |
| iOS 15 | Untested | Untested | Untested | Untested | Untested |
| iOS 16 | Untested | Untested | Untested | Untested | Untested |
| iOS 17 | Untested | Untested | Untested | Untested | Untested |
| iOS 18 | Untested | Untested | Untested | Untested | Untested |

| Connected device | Wi-Fi + manual HTTP | Wi-Fi + PAC |
| --- | --- | --- |
| Windows 10/11 | Untested | Untested |
| Linux (specify distribution/desktop) | Untested | Untested |
| macOS | Untested | Untested |
| iPadOS | Untested | Untested |
| iOS | Untested | Untested |

## Test procedure

1. Confirm that cellular internet works on the iPhone itself.
2. Turn off the connected device's proxy and record normal hotspot connectivity.
3. Run `iproxy` in iSH and open the PAC URL on the connected device.
4. Run `curl --proxy http://<phone-IP>:8080 http://example.com`, replacing `<phone-IP>` with the phone's address.
5. Access `https://example.com` through the same HTTP proxy to test CONNECT.
6. Run `curl --proxy socks5h://<phone-IP>:1080 https://example.com` to test SOCKS5.
7. Configure PAC in the OS and check HTTP/HTTPS in a browser. Avoid overlapping proxy or VPN settings.
8. After locking the screen or switching apps, test again after 1, 5, and 15 minutes. Record whether `-l` was enabled.
9. Remove the proxy settings and confirm normal internet connectivity returns.

## Results template

```text
Date / iProxy version:
iPhone model / exact iOS version / iSH version:
Connected device / OS / app version:
Connection method / phone IP / command-line options:
Normal hotspot / HTTP / HTTPS CONNECT / SOCKS5 / PAC results:
Results after 1, 5, and 15 minutes of screen lock / whether -l was enabled:
Error messages (redact URLs, IPs, and personal information as needed):
```
