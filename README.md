# iProxy

아이폰의 iSH에서 실행하는 작은 HTTP / SOCKS5 프록시입니다. 아이폰 개인용 핫스팟에 연결한 Windows, Linux, Mac, iPad, iPhone에서 **프록시 설정을 지원하는 앱의 TCP 통신**을 전달합니다.

> **iOS 15·16·17·18을 지원 목표로 정비 중이며, 각 버전의 실기기 검증은 아직 완료되지 않았습니다.** VPN이나 모든 앱의 트래픽을 전달하는 도구가 아닙니다. 통신사의 테더링 허용, 한도 우회, 무제한 사용은 보장하지 않습니다.

## 아이폰에서 설치

탈옥이나 Mac/Xcode는 필요하지 않습니다. [iSH Shell](https://apps.apple.com/app/ish-shell/id1436902243)을 설치하세요. iSH는 x86 Linux 환경이므로 아이폰에서도 **`linux-386` 파일**을 사용합니다. `darwin-arm64`는 Mac용입니다.

1. 설정 → 개인용 핫스팟 → **다른 사람의 연결 허용**을 켭니다. 다른 기기를 이 핫스팟 Wi-Fi에 연결합니다.
2. [Releases](https://github.com/JJo-0/iproxy/releases)에서 게시된 버전이 있는지 확인합니다. 아래 명령은 새 빌드의 릴리스가 게시된 뒤 사용할 수 있습니다. 아직 없다면 아래 소스 빌드 방법을 사용하세요.
3. iSH에서 다음 명령을 실행합니다. 다운로드 또는 체크섬 확인이 실패하면 설치가 중단됩니다.

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

4. iSH를 열어 둔 채 실행합니다. 로컬 네트워크 접근 권한을 묻는 경우 허용하세요.

```sh
iproxy
```

기본 포트는 HTTP **8080**, SOCKS5 **1080**, PAC **8090**입니다. 화면에 `listening` 메시지가 나와야 정상적으로 시작한 것입니다. 중지는 **Ctrl+C**입니다. 기기 재부팅이나 iSH 종료 후에는 다시 실행해야 합니다.

기본 아이폰 주소는 `172.20.10.1`입니다. 연결 기기의 Wi-Fi 정보에서 **라우터/기본 게이트웨이** 주소를 확인하세요. 다르면 다음처럼 바꿉니다.

```sh
iproxy -a 172.20.10.1
```

`-a`는 PAC에 알려 주는 주소이고, `-b`는 실제 수신 인터페이스입니다. 기본 `-b 0.0.0.0`은 모든 IPv4 인터페이스에서 수신합니다. 인증 기능이 없으므로 비밀번호가 있는 본인 핫스팟에서 신뢰하는 기기와 사용하세요.

## 연결 기기 설정

먼저 연결 기기 브라우저에서 `http://172.20.10.1:8090/proxy`를 열어 PAC 코드가 표시되거나 다운로드되는지 확인합니다. 이것은 아이폰까지의 연결 확인이며, 인터넷 전달 성공을 의미하지는 않습니다.

| 연결 기기 | 권장 설정 |
| --- | --- |
| Windows 10/11 | 설정 → 네트워크 및 인터넷 → 프록시 → 설정 스크립트 사용 → 아래 PAC URL 입력 |
| macOS | 시스템 설정 → 네트워크 → 연결된 Wi-Fi의 세부사항 → 프록시 → 자동 프록시 구성 → PAC URL 입력 (구버전은 시스템 환경설정 → 네트워크 → 고급) |
| Linux | 데스크톱 네트워크 설정의 프록시 → 자동 → PAC URL 입력. 데스크톱/앱별 지원 차이가 있으므로 CLI는 아래 명령 사용 |
| iPad / iPhone | 설정 → Wi-Fi → 아이폰 핫스팟 옆 ⓘ → HTTP 프록시 구성 → 자동 → PAC URL 입력 → 저장 |

```text
http://172.20.10.1:8090/proxy
```

자동 설정이 적용되지 않는 앱은 수동 HTTP 프록시를 **서버 `172.20.10.1`, 포트 `8080`, 인증 끔**으로 지정합니다. HTTPS도 동일한 HTTP 프록시의 CONNECT 터널을 사용하며 인증서 설치는 필요하지 않습니다. iOS/iPadOS의 Wi-Fi 프록시 설정은 이를 따르는 앱에만 적용됩니다.

CLI 확인 (Windows에서는 `curl.exe` 사용):

```sh
curl --proxy http://172.20.10.1:8080 https://example.com
curl --proxy socks5h://172.20.10.1:1080 https://example.com
```

Linux/macOS 셸에서 지원하는 프로그램에 적용:

```sh
export http_proxy=http://172.20.10.1:8080
export https_proxy=http://172.20.10.1:8080
export all_proxy=socks5h://172.20.10.1:1080
```

사용을 끝내면 각 기기의 프록시 설정을 **끔/없음**으로 되돌리세요. 셸에서는 `unset http_proxy https_proxy all_proxy`를 실행합니다. 프록시를 켜 둔 채 핫스팟을 끊으면 인터넷이 안 되는 것처럼 보일 수 있습니다.

## 지원 범위와 통신사

- **Wi-Fi 핫스팟 연결**을 기본 경로로 안내합니다. USB/Bluetooth는 OS·드라이버에 따라 다르며 별도 실험 대상입니다.
- HTTP, HTTPS CONNECT, SOCKS5 TCP를 제공합니다. UDP, QUIC/HTTP3, 게임·메신저·VPN 등 모든 앱의 동작은 보장하지 않습니다. 앱의 TCP 대체 경로가 필요할 수 있습니다.
- PAC URL을 직접 입력하는 방식입니다. WPAD 자동 검색 기능은 제공하지 않습니다.
- 기본 PAC는 HTTP → SOCKS5 순서이며 `DIRECT`를 포함하지 않습니다. `-direct`를 켜면 프록시 실패 시 일반 연결로 전환할 수 있습니다. PAC 자체를 무시하는 앱까지 막는 기능은 아닙니다.
- KT·SKT·LG U+ 및 알뜰폰은 **검증 전**입니다. `sky`가 SKT인지 스카이라이프 등 다른 사업자인지에 따라 따로 기록해야 합니다. 통신사명만으로 판단할 수 없으며 요금제·테더링 잔여량·iOS·기종을 함께 확인해야 합니다.
- iOS가 iSH를 백그라운드에서 정지하거나 종료할 수 있습니다. `-l`은 `/dev/location`을 읽는 실험 옵션이며 지속 실행을 보장하지 않습니다. 좌표는 저장·출력하지 않습니다. 위치 장치가 없어도 프록시는 계속 실행됩니다.

버전별·기기별·통신사별 결과는 [호환성 검증표](docs/compatibility.md)에 기록합니다. 문제가 있으면 [문제 해결](docs/troubleshooting.md)을 먼저 확인하세요.

## 옵션

| 옵션 | 기본값 | 설명 |
| --- | --- | --- |
| `-a` | `172.20.10.1` | PAC에 광고할 아이폰 IP |
| `-b` | `0.0.0.0` | 수신할 로컬 IP |
| `-p` | `8080` | HTTP/HTTPS CONNECT 포트, `0`이면 끔 |
| `-s` | `1080` | SOCKS5 TCP 포트, `0`이면 끔 |
| `-d` | `8090` | PAC `/proxy` 포트, `0`이면 끔 |
| `-direct` | 꺼짐 | PAC의 일반 연결 fallback 허용 |
| `-l` | 꺼짐 | iSH 위치 장치 읽기 시도 |
| `-v` | 꺼짐 | HTTP 상세 로그 (URL 등 민감 정보 포함 가능) |
| `-version` | — | 빌드 버전 표시 |
| `-h` | — | 도움말 |

기존 버전과 달리 포트 옵션 없이 실행해도 서비스가 켜집니다. 일부 기능만 쓰려면 나머지 포트를 `0`으로 지정하세요.

```sh
iproxy -p 8080 -s 0 -d 0
```

## 빌드와 배포

Go 버전은 `go.mod`에 명시되어 있습니다. Linux/macOS 또는 Windows의 WSL에서 빌드하는 것을 권장합니다. iSH 내부의 Go 빌드는 느리거나 호환성 문제로 실패할 수 있습니다.

```sh
git clone https://github.com/JJo-0/iproxy.git
cd iproxy
go test -race ./...
go vet ./...
./scripts/build.sh dev
```

`dist/iproxy-linux-386`을 아이폰 파일 앱의 iSH 폴더로 복사하고 iSH에서 해당 경로로 이동한 뒤 설치합니다.

```sh
install -m 755 ./iproxy-linux-386 /usr/local/bin/iproxy
iproxy
```

빌드는 iSH용 Linux 386, Linux amd64/arm64, Mac amd64/arm64, Windows amd64/arm64와 SHA-256 목록을 생성합니다. 연결 기기에는 iProxy 설치가 필요 없으며 OS 프록시 설정만 하면 됩니다. 태그 `v*`를 push하면 GitHub Actions가 테스트·빌드 후 **초안 릴리스**를 만듭니다. 실기기 결과와 파일을 검토한 다음 게시하세요.

## 근거 자료

- [iSH 공식 저장소: x86 Linux 에뮬레이션](https://github.com/ish-app/ish)
- [Apple: 개인용 핫스팟 설정 및 연결](https://support.apple.com/en-us/111785)
- [Apple: Mac 프록시 설정](https://support.apple.com/guide/mac-help/change-proxy-settings-on-mac-mchlp2591/mac)

이 자료는 플랫폼의 설정 방법을 설명하며, 이 프로젝트나 특정 요금제의 실기기 검증을 대신하지 않습니다.
