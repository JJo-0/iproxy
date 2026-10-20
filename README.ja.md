# iProxy

[한국어](README.md) · [English](README.en.md) · [简体中文](README.zh-CN.md) · **日本語**

iPhone の iSH 上で動作する小さな HTTP / SOCKS5 プロキシです。iPhone のインターネット共有に接続した Windows、Linux、Mac、iPad、iPhone から、**プロキシ設定に対応するアプリの TCP 通信**を転送します。

> **iOS 15・16・17・18 を対応目標としていますが、各バージョンでの実機検証はまだ完了していません。** VPN ではなく、すべてのアプリの通信を転送するものではありません。

## iPhone にインストール

脱獄や Mac/Xcode は不要です。[iSH Shell](https://apps.apple.com/app/ish-shell/id1436902243) をインストールしてください。iSH は x86 Linux 環境のため、iPhone でも **`linux-386` の実行ファイル**を使用します。`darwin-arm64` は Mac 用です。

1. 「設定 → インターネット共有」で**ほかの人の接続を許可**をオンにします。接続する端末をこの Wi-Fi に接続します。
2. [Releases](https://github.com/JJo-0/iproxy/releases) に公開済みのバージョンがあるか確認します。以下のダウンロードコマンドは、新しいビルドを含むリリースの公開後に使用できます。未公開の場合は、[Actions](https://github.com/JJo-0/iproxy/actions/workflows/ci.yml) の該当ブランチの成功した実行から `iproxy-binaries` アーティファクトをダウンロードして展開し、後述のファイルコピー手順でインストールできます（GitHub へのログインが必要）。ソースからのビルドも可能です。
3. iSH で次のコマンドを実行します。ダウンロードまたはチェックサムの確認に失敗すると、インストールは中止されます。

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

4. iSH を開いたまま実行します。ローカルネットワークへのアクセスを求められた場合は許可してください。

```sh
iproxy
```

標準ポートは HTTP **8080**、SOCKS5 **1080**、PAC **8090** です。`listening` と表示されれば待ち受けが開始されています。停止は **Ctrl+C** です。端末の再起動や iSH の終了後は、もう一度実行してください。

iPhone の既定アドレスは `172.20.10.1` です。接続端末の Wi-Fi 詳細で**ルーター／デフォルトゲートウェイ**のアドレスを確認してください。異なる場合は、次のコマンドと接続端末の設定のアドレスを置き換えます。

```sh
iproxy -a 172.20.10.1
```

`-a` は PAC で通知するアドレス、`-b` は実際に待ち受けるローカルアドレスです。既定の `-b 0.0.0.0` はすべての IPv4 インターフェースで待ち受けます。認証機能はないため、パスワードを設定した自分のインターネット共有で、信頼できる端末と使用してください。

## 接続端末の設定

まず接続端末のブラウザーで `http://172.20.10.1:8090/proxy` を開き、PAC のコードが表示またはダウンロードされることを確認します。これは iPhone への接続確認であり、インターネット通信の転送成功を意味するものではありません。

| 端末 | 推奨設定 |
| --- | --- |
| Windows 10/11 | 設定 → ネットワークとインターネット → プロキシ → セットアップスクリプトを使う → 下記の PAC URL を入力 |
| macOS | システム設定 → ネットワーク → 接続中の Wi-Fi の詳細 → プロキシ → 自動プロキシ構成 → PAC URL を入力（旧バージョン：システム環境設定 → ネットワーク → 詳細） |
| Linux | デスクトップのネットワーク設定 → プロキシ → 自動 → PAC URL を入力。デスクトップやアプリにより対応が異なるため、CLI では下記のコマンドを使用します。 |
| iPad / iPhone | 設定 → Wi-Fi → 共有元 iPhone のネットワークの横にある ⓘ → HTTP プロキシの「プロキシを構成」→ 自動 → PAC URL を入力 → 保存 |

```text
http://172.20.10.1:8090/proxy
```

自動設定が適用されないアプリでは、HTTP プロキシを手動で**サーバー `172.20.10.1`、ポート `8080`、認証オフ**に設定します。HTTPS も同じ HTTP プロキシの CONNECT トンネルを使うため、証明書のインストールは不要です。iOS/iPadOS の Wi-Fi プロキシ設定は、その設定に従うアプリにのみ適用されます。

CLI での確認（Windows では `curl.exe` を使用）：

```sh
curl --proxy http://172.20.10.1:8080 https://example.com
curl --proxy socks5h://172.20.10.1:1080 https://example.com
```

Linux/macOS のシェルで、対応するプログラムに設定する場合：

```sh
export http_proxy=http://172.20.10.1:8080
export https_proxy=http://172.20.10.1:8080
export all_proxy=socks5h://172.20.10.1:1080
```

使い終わったら、各端末のプロキシ設定を**オフ／なし**に戻してください。シェルでは `unset http_proxy https_proxy all_proxy` を実行します。プロキシ設定を残したまま共有接続を切ると、インターネットに接続できなくなる場合があります。

## 対応範囲と制限

- **Wi-Fi によるインターネット共有**を基本の接続方法としています。USB/Bluetooth は OS やドライバーに依存するため、別途検証が必要です。
- HTTP、HTTPS CONNECT、SOCKS5 TCP に対応します。UDP、QUIC/HTTP3、すべてのゲーム・メッセージアプリ・VPN の動作は保証しません。アプリ側で TCP へのフォールバックが必要になる場合があります。
- PAC URL を直接入力する方式です。WPAD による自動検出は提供しません。
- 既定の PAC は HTTP → SOCKS5 の順に試し、`DIRECT` を含みません。`-direct` を有効にすると、プロキシが失敗した際に直接接続へ切り替えられます。PAC を無視するアプリの通信を遮断する機能ではありません。
- iOS がバックグラウンドの iSH を一時停止または終了する場合があります。`-l` は `/dev/location` を読み取る実験的なオプションで、継続動作は保証しません。座標の保存や出力は行いません。位置情報デバイスが利用できなくても、プロキシは動作を続けます。

バージョン・端末ごとの結果は[互換性チェックリスト](docs/compatibility.ja.md)に記録します。問題がある場合は[トラブルシューティング](docs/troubleshooting.ja.md)を確認してください。

## オプション

| オプション | 既定値 | 説明 |
| --- | --- | --- |
| `-a` | `172.20.10.1` | PAC で通知する iPhone の IP |
| `-b` | `0.0.0.0` | 待ち受けるローカル IP |
| `-p` | `8080` | HTTP/HTTPS CONNECT ポート。`0` で無効 |
| `-s` | `1080` | SOCKS5 TCP ポート。`0` で無効 |
| `-d` | `8090` | PAC `/proxy` ポート。`0` で無効 |
| `-direct` | オフ | PAC の直接接続フォールバックを許可 |
| `-l` | オフ | iSH の位置情報デバイスの読み取りを試行 |
| `-v` | オフ | HTTP の詳細ログ（機密性のある URL を含む場合あり） |
| `-version` | — | ビルドバージョンを表示 |
| `-h` | — | ヘルプを表示 |

元のバージョンとは異なり、ポート指定なしでもサービスが起動します。一部の機能だけを使う場合は、不要なポートを `0` にしてください。

```sh
iproxy -p 8080 -s 0 -d 0
```

## ビルドとリリース

必要な Go バージョンは `go.mod` に記載しています。Linux/macOS、または Windows の WSL でのビルドを推奨します。iSH 内での Go ビルドは遅い場合や、互換性の問題で失敗する場合があります。

```sh
git clone https://github.com/JJo-0/iproxy.git
cd iproxy
go test -race ./...
go vet ./...
./scripts/build.sh dev
```

`dist/iproxy-linux-386` を iPhone の「ファイル」アプリ内の iSH フォルダーにコピーし、iSH でそのディレクトリに移動してインストールします。

```sh
install -m 755 ./iproxy-linux-386 /usr/local/bin/iproxy
iproxy
```

ビルドでは iSH 用の Linux 386、Linux amd64/arm64、Mac amd64/arm64、Windows amd64/arm64 の実行ファイルと SHA-256 一覧を生成します。接続端末への iProxy のインストールは不要で、プロキシ設定のみ必要です。`v*` タグを push すると、GitHub Actions がテスト・ビルド後に**ドラフトリリース**を作成します。実機での結果とファイルを確認してから公開してください。

## 参考資料

- [iSH 公式リポジトリ：x86 Linux エミュレーション](https://github.com/ish-app/ish)
- [Apple：インターネット共有の設定と接続](https://support.apple.com/en-us/111785)
- [Apple：Mac のプロキシ設定を変更する](https://support.apple.com/guide/mac-help/change-proxy-settings-on-mac-mchlp2591/mac)

これらの資料はプラットフォームの設定方法を説明するもので、本プロジェクトの実機検証に代わるものではありません。
