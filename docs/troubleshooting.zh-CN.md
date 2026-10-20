# 故障排查

[한국어](troubleshooting.md) · [English](troubleshooting.en.md) · **简体中文** · [日本語](troubleshooting.ja.md) · [README](../README.zh-CN.md)

| 现象 | 检查方法 |
| --- | --- |
| 下载返回 404 | 新版本可能尚未发布，或文件名不同。请检查 Releases 的文件列表，或按 README 从源码构建。 |
| `Exec format error` | iSH 需要 `iproxy-linux-386`，不能使用 Mac/ARM 可执行文件。 |
| `Permission denied` | 执行 `chmod +x`，或使用 README 中的 `install -m 755` 命令。 |
| 启动时出现端口错误 | 检查是否已有 iProxy 在运行。如需更改端口，也要同步修改连接设备的设置。 |
| 无法打开 PAC URL | 检查热点连接、实际网关 IP、iSH 是否正在运行，以及本地网络权限。 |
| 能看到 PAC 代码，但无法访问网页 | 使用手动 HTTP 代理和 curl 测试，区分 PAC 配置问题与代理转发问题。也请检查 iPhone 本身的蜂窝网络连接。 |
| 只有浏览器能使用 | 检查其他应用是否支持系统代理设置。本工具不是支持 UDP 和所有应用的 VPN。 |
| 锁屏后断开 | 重新打开 iSH，检查进程是否仍在运行。对于 iOS 暂停应用，没有保证有效的解决办法。`-l` 是可选的实验功能。 |
| `/dev/location` 错误 | 当前 iSH 环境或版本可能不提供位置设备。请去掉 `-l` 后运行。 |
| 停止代理后无法上网 | 移除系统 PAC/手动代理设置，以及 shell 中的代理环境变量。 |

`-v` 日志可能包含请求 URL。分享日志前请删除个人信息，并注明复现问题所需的设备、操作系统、应用版本和启动选项。
