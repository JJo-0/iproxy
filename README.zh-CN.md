# iProxy

[한국어](README.md) · [English](README.en.md) · **简体中文** · [日本語](README.ja.md)

在 iPhone 的 iSH 中运行的轻量级 HTTP / SOCKS5 代理。Windows、Linux、Mac、iPad 和 iPhone 连接到这台 iPhone 的个人热点后，可通过本工具转发**支持代理设置的应用的 TCP 流量**。

> **目标是支持 iOS 15、16、17 和 18，但各版本均尚未完成真机验证。** 本工具不是 VPN，也不会转发所有应用的流量。

## 在 iPhone 上安装

无需越狱，也不需要 Mac/Xcode。请安装 [iSH Shell](https://apps.apple.com/app/ish-shell/id1436902243)。iSH 使用 x86 Linux 环境，因此即使在 iPhone 上，也应选择 **`linux-386` 可执行文件**。`darwin-arm64` 是 Mac 版本。

1. 打开“设置 → 个人热点”，启用**允许其他人加入**。让其他设备通过 Wi-Fi 连接此热点。
2. 在 [Releases](https://github.com/JJo-0/iproxy/releases) 中检查是否已有发布版本。以下下载命令仅在包含新构建文件的版本发布后可用。在此之前，可从 [Actions](https://github.com/JJo-0/iproxy/actions/workflows/ci.yml) 中对应分支的成功运行记录下载 `iproxy-binaries` 构件并解压，再按下方的文件复制步骤安装（需要登录 GitHub）。也可从源码构建。
3. 在 iSH 中执行以下命令。下载或校验失败时，安装会停止。

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

4. 保持 iSH 打开并运行代理。如果出现本地网络访问提示，请允许访问。

```sh
iproxy
```

默认端口：HTTP **8080**、SOCKS5 **1080**、PAC **8090**。出现 `listening` 消息表示服务已开始监听。按 **Ctrl+C** 停止。手机重启或 iSH 关闭后，需要重新运行。

默认 iPhone 地址为 `172.20.10.1`。请在连接设备的 Wi-Fi 详情中检查**路由器/默认网关**地址。如果不同，请替换以下命令以及客户端设置中的地址：

```sh
iproxy -a 172.20.10.1
```

`-a` 指定 PAC 中提供给客户端的地址，`-b` 指定本地监听地址。默认 `-b 0.0.0.0` 会监听所有 IPv4 接口。代理不提供身份验证，请仅在自己设置了密码的热点上与可信设备使用。

## 配置连接设备

先在连接设备的浏览器中打开 `http://172.20.10.1:8090/proxy`，确认能显示或下载 PAC 代码。这只能确认设备能连接到 iPhone，并不代表互联网流量已成功转发。

| 设备 | 推荐设置 |
| --- | --- |
| Windows 10/11 | 设置 → 网络和 Internet → 代理 → 使用设置脚本 → 填入下方 PAC URL |
| macOS | 系统设置 → 网络 → 已连接 Wi-Fi 的详细信息 → 代理 → 自动代理配置 → 填入 PAC URL（旧版本：系统偏好设置 → 网络 → 高级） |
| Linux | 桌面网络设置 → 代理 → 自动 → 填入 PAC URL。不同桌面和应用的支持情况有差异，命令行工具可使用下方命令。 |
| iPad / iPhone | 设置 → Wi-Fi → iPhone 热点旁的 ⓘ → HTTP 代理中的“配置代理” → 自动 → 填入 PAC URL → 存储 |

```text
http://172.20.10.1:8090/proxy
```

若应用不采用自动配置，可手动设置 HTTP 代理：**服务器 `172.20.10.1`、端口 `8080`、关闭认证**。HTTPS 也使用同一 HTTP 代理的 CONNECT 隧道，无需安装证书。iOS/iPadOS 的 Wi-Fi 代理设置只对遵循该设置的应用生效。

命令行检查（Windows 请使用 `curl.exe`）：

```sh
curl --proxy http://172.20.10.1:8080 https://example.com
curl --proxy socks5h://172.20.10.1:1080 https://example.com
```

在 Linux/macOS shell 中为支持这些变量的程序设置代理：

```sh
export http_proxy=http://172.20.10.1:8080
export https_proxy=http://172.20.10.1:8080
export all_proxy=socks5h://172.20.10.1:1080
```

使用结束后，请将设备的代理设置恢复为**关闭/无**。在 shell 中运行 `unset http_proxy https_proxy all_proxy`。断开热点后若仍保留代理设置，可能会出现无法访问互联网的现象。

## 支持范围与限制

- 本文默认使用 **Wi-Fi 热点连接**。USB/Bluetooth 取决于操作系统和驱动，需要单独测试。
- 支持 HTTP、HTTPS CONNECT 和 SOCKS5 TCP。不保证 UDP、QUIC/HTTP3，以及所有游戏、聊天应用或 VPN 都能正常工作。应用可能需要回退到 TCP。
- 需要手动填写 PAC URL，不提供 WPAD 自动发现。
- 默认 PAC 先尝试 HTTP，再尝试 SOCKS5，不包含 `DIRECT`。启用 `-direct` 后，代理失败时可回退到直接连接。这不会阻止忽略 PAC 的应用直接联网。
- iOS 可能在后台暂停或终止 iSH。`-l` 是读取 `/dev/location` 的实验选项，不保证持续运行。坐标不会被保存或输出。位置设备不可用时，代理仍会继续运行。

请在[兼容性验证表](docs/compatibility.zh-CN.md)中记录各版本和设备的结果。如遇问题，请先查看[故障排查](docs/troubleshooting.zh-CN.md)。

## 选项

| 选项 | 默认值 | 说明 |
| --- | --- | --- |
| `-a` | `172.20.10.1` | PAC 中提供的 iPhone IP |
| `-b` | `0.0.0.0` | 本地监听 IP |
| `-p` | `8080` | HTTP/HTTPS CONNECT 端口，`0` 表示关闭 |
| `-s` | `1080` | SOCKS5 TCP 端口，`0` 表示关闭 |
| `-d` | `8090` | PAC `/proxy` 端口，`0` 表示关闭 |
| `-direct` | 关闭 | 允许 PAC 回退到直接连接 |
| `-l` | 关闭 | 尝试读取 iSH 位置设备 |
| `-v` | 关闭 | HTTP 详细日志（可能包含敏感 URL） |
| `-version` | — | 显示构建版本 |
| `-h` | — | 显示帮助 |

与原始版本不同，不指定端口选项也会启动服务。若只使用部分功能，请将其他端口设为 `0`。

```sh
iproxy -p 8080 -s 0 -d 0
```

## 构建与发布

所需 Go 版本见 `go.mod`。建议在 Linux/macOS 或 Windows 的 WSL 中构建。在 iSH 内构建 Go 程序可能很慢，也可能因兼容性问题失败。

```sh
git clone https://github.com/JJo-0/iproxy.git
cd iproxy
go test -race ./...
go vet ./...
./scripts/build.sh dev
```

将 `dist/iproxy-linux-386` 复制到 iPhone“文件”应用中的 iSH 文件夹。在 iSH 中进入该目录并安装：

```sh
install -m 755 ./iproxy-linux-386 /usr/local/bin/iproxy
iproxy
```

构建会生成用于 iSH 的 Linux 386、Linux amd64/arm64、Mac amd64/arm64、Windows amd64/arm64 文件及 SHA-256 校验列表。连接设备无需安装 iProxy，只需配置代理。推送 `v*` 标签后，GitHub Actions 会测试、构建并创建**草稿版本**。请检查文件和真机测试结果后再正式发布。

## 参考资料

- [iSH 官方仓库：x86 Linux 模拟](https://github.com/ish-app/ish)
- [Apple：设置和连接个人热点](https://support.apple.com/en-us/111785)
- [Apple：更改 Mac 代理设置](https://support.apple.com/guide/mac-help/change-proxy-settings-on-mac-mchlp2591/mac)

这些资料说明平台的设置方法，不能代替本项目的真机验证。
