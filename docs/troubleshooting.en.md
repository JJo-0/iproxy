# Troubleshooting

[한국어](troubleshooting.md) · **English** · [简体中文](troubleshooting.zh-CN.md) · [日本語](troubleshooting.ja.md) · [README](../README.en.md)

| Symptom | What to check |
| --- | --- |
| Download returns 404 | The new release may not be published yet, or the filename may differ. Check the release assets or use the source-build instructions in the README. |
| `Exec format error` | iSH needs `iproxy-linux-386`. Mac/ARM binaries will not work in iSH. |
| `Permission denied` | Use `chmod +x` or the README's `install -m 755` command. |
| Port error at startup | Check for another running iProxy instance. If you change ports, also update the client settings. |
| Cannot open PAC URL | Check hotspot connectivity, the actual gateway IP, whether iSH is running, and local network permission. |
| PAC code loads but websites fail | Test manual HTTP proxy settings and curl to distinguish PAC configuration issues from forwarding issues. Also check the phone's cellular internet connection. |
| Only the browser works | Check whether other apps support system proxy settings. This is not a VPN that supports UDP and every app. |
| Disconnects after screen lock | Reopen iSH and check whether the process is running. There is no guaranteed solution to iOS suspending an app. `-l` is an optional experiment. |
| `/dev/location` error | The location device may be unavailable in your iSH environment/version. Run without `-l`. |
| No internet after stopping the proxy | Remove the OS's PAC/manual proxy settings and the shell proxy environment variables. |

`-v` logs may include request URLs. Remove personal information before sharing logs. Include the device, OS, app version, and command-line options needed to reproduce the issue.
