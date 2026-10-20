# 兼容性验证表

[한국어](compatibility.md) · [English](compatibility.en.md) · **简体中文** · [日本語](compatibility.ja.md) · [README](../README.zh-CN.md)

“未验证”不代表失败，而是尚无可复现的真机测试结果。仅凭 Linux 自动测试和交叉编译，无法确认 iOS 兼容性。

| iPhone 系统 | iSH 安装与启动 | HTTP / HTTPS | SOCKS5 | PAC | 锁屏与后台 |
| --- | --- | --- | --- | --- | --- |
| iOS 15 | 未验证 | 未验证 | 未验证 | 未验证 | 未验证 |
| iOS 16 | 未验证 | 未验证 | 未验证 | 未验证 | 未验证 |
| iOS 17 | 未验证 | 未验证 | 未验证 | 未验证 | 未验证 |
| iOS 18 | 未验证 | 未验证 | 未验证 | 未验证 | 未验证 |

| 连接设备 | Wi-Fi + 手动 HTTP 代理 | Wi-Fi + PAC |
| --- | --- | --- |
| Windows 10/11 | 未验证 | 未验证 |
| Linux（注明发行版和桌面环境） | 未验证 | 未验证 |
| macOS | 未验证 | 未验证 |
| iPadOS | 未验证 | 未验证 |
| iOS | 未验证 | 未验证 |

## 验证步骤

1. 确认 iPhone 本身可以通过蜂窝网络访问互联网。
2. 关闭连接设备的代理，记录普通热点连接的结果。
3. 在 iSH 中运行 `iproxy`，然后在连接设备上打开 PAC URL。
4. 运行 `curl --proxy http://<phone-IP>:8080 http://example.com`，将 `<phone-IP>` 替换为 iPhone 的地址。
5. 通过同一个 HTTP 代理访问 `https://example.com`，验证 CONNECT。
6. 运行 `curl --proxy socks5h://<phone-IP>:1080 https://example.com`，验证 SOCKS5。
7. 配置系统 PAC 后，在浏览器中检查 HTTP/HTTPS。避免叠加其他代理或 VPN 设置。
8. 锁屏或切换应用后，分别在 1、5、15 分钟再次检查，并记录是否启用了 `-l`。
9. 移除代理设置，确认普通互联网连接恢复正常。

## 结果记录模板

```text
日期 / iProxy 版本：
iPhone 型号 / 完整 iOS 版本 / iSH 版本：
连接设备 / 操作系统 / 应用版本：
连接方式 / iPhone IP / 启动选项：
普通热点 / HTTP / HTTPS CONNECT / SOCKS5 / PAC 结果：
锁屏 1、5、15 分钟后的结果 / 是否启用 -l：
错误信息（按需隐藏 URL、IP 和个人信息）：
```
