# meteor-shower

Go標準パッケージのみで構築された負荷試験ツール。

## 前提条件

- Go 1.21 以上
- Git (オプション: ソースからビルドする場合)

## インストール

### ソースからビルド

```bash
git clone https://github.com/kitsystemyou/meteor-shower.git
cd meteor-shower
go build -o meteor-shower ./cmd/meteor-shower
```

### Go install を使用

```bash
go install github.com/kitsystemyou/meteor-shower/cmd/meteor-shower@latest
```

インストール後、`$GOPATH/bin/meteor-shower` として利用可能になります。

## 使い方

### 基本的な使い方

```bash
# デフォルト設定で負荷試験を実行
meteor-shower run

# RPSを指定して実行
meteor-shower run --rps 100

# 並列数を指定して実行
meteor-shower run --concurrency 10

# JSON形式で結果を出力
meteor-shower run -o json

# すべてのパラメータを指定
meteor-shower run --rps 50 --concurrency 5 -o html
```

### グローバルフラグ

すべてのコマンドで使用できるフラグ:

| フラグ | 短縮形 | デフォルト | 説明 |
|--------|--------|------------|------|
| `--config` | - | `./config.yaml` | 設定ファイルのパス |
| `--rps` | - | 設定ファイル参照 | 秒間リクエスト数 (設定ファイルより優先) |
| `--concurrency` | - | 設定ファイル参照 | 並列クライアント数 (設定ファイルより優先) |
| `--output` | `-o` | 設定ファイル参照 | 出力形式: html, json (設定ファイルより優先) |

### サブコマンド

#### `run` - 負荷試験を実行

指定されたエンドポイントに対して負荷試験を実行します。

```bash
meteor-shower run [flags]
```

**フラグ:**
- `--rps int`: 秒間リクエスト数 (設定ファイルより優先)
- `--concurrency int`: 並列クライアント数 (設定ファイルより優先)
- `-o, --output string`: 出力形式 (html, json)

**例:**

```bash
# デフォルト設定で実行
meteor-shower run

# RPSを指定
meteor-shower run --rps 100

# 並列数を指定
meteor-shower run --concurrency 10

# JSON形式で出力
meteor-shower run -o json

# すべてのパラメータを指定
meteor-shower run --rps 50 --concurrency 5 -o html

# カスタム設定ファイルを使用
meteor-shower run --config /path/to/config.yaml
```

#### `version` - バージョン情報を表示

CLIツールのバージョン情報を表示します。

```bash
meteor-shower version
```

#### `help` - ヘルプを表示

コマンドのヘルプ情報を表示します。

```bash
# 全体のヘルプ
meteor-shower help

# 特定のコマンドのヘルプ
meteor-shower help run
meteor-shower run --help
```

## 設定ファイル

`config.yaml` ファイルで設定を管理できます。

### 設定ファイルの場所

以下の順序で設定ファイルを検索します:

1. `--config` フラグで指定されたパス
2. カレントディレクトリの `config.yaml`
3. `$HOME/.meteor-shower/config.yaml`

### 設定例

#### 単一エンドポイント

```yaml
loadtest:
  # ターゲットドメイン
  domain: "http://localhost:8080"
  
  # エンドポイント
  endpoints:
    - path: "/"
      weight: 1.0
  
  # 秒間リクエスト数
  rps: 10
  
  # 並列クライアント数
  concurrency: 1
  
  # テスト実行時間 (秒)
  duration: 10
  
  # 出力形式 (html または json)
  output: "html"
```

#### 複数エンドポイント (重み付き)

```yaml
loadtest:
  # ターゲットドメイン
  domain: "http://localhost:8080"
  
  # 複数エンドポイントと重み
  # 重みに応じてリクエストが分散されます
  endpoints:
    - path: "/"
      weight: 1.0      # 最も高い頻度
    - path: "/health"
      weight: 0.5      # 中程度の頻度
    - path: "/slow"
      weight: 0.2      # 低い頻度
  
  # 秒間リクエスト数
  rps: 10
  
  # 並列クライアント数
  concurrency: 1
  
  # テスト実行時間 (秒)
  duration: 10
  
  # 出力形式 (html または json)
  output: "html"
```

### 設定項目

| 項目 | 型 | デフォルト | 説明 |
|------|-----|-----------|------|
| `loadtest.domain` | string | `"http://localhost:8080"` | ターゲットドメイン |
| `loadtest.endpoints` | array | `[{path: "/", weight: 1.0}]` | エンドポイント設定 (必須) |
| `loadtest.endpoints[].path` | string | - | エンドポイントのパス |
| `loadtest.endpoints[].weight` | float | `1.0` | リクエスト分散の重み |
| `loadtest.rps` | int | `10` | 秒間リクエスト数 |
| `loadtest.concurrency` | int | `1` | 並列クライアント数 |
| `loadtest.duration` | int | `10` | テスト実行時間 (秒) |
| `loadtest.output` | string | `"html"` | 出力形式 (html, json) |

## 出力形式

### HTML形式

HTMLレポートには以下の情報が含まれます:
- テスト設定 (URL, RPS, 並列数, 実行時間)
- サマリー (総リクエスト数, 成功/失敗数, 実際のRPS)
- レスポンスタイム統計 (最小/平均/中央値/95パーセンタイル/99パーセンタイル/最大)
- ステータスコード分布

```bash
meteor-shower run -o html > report.html
```

### JSON形式

JSON形式では、すべてのリクエスト結果を含む詳細なデータが出力されます:

```bash
meteor-shower run -o json > report.json
```

JSON出力例:
```json
{
  "url": "http://localhost:8080/",
  "rps": 10,
  "concurrency": 1,
  "duration": 10,
  "statistics": {
    "total_requests": 100,
    "success_requests": 100,
    "failed_requests": 0,
    "avg_duration_ms": 15,
    "p95_duration_ms": 25,
    "requests_per_sec": 10.5
  },
  "status_codes": {
    "200": 100
  }
}
```

## 開発

### プロジェクト構造

```
.
├── cmd/
│   └── meteor-shower/          # メインエントリーポイント
│       └── main.go
├── internal/
│   ├── cli/            # CLI実装
│   │   ├── cli.go      # メインCLIロジック
│   │   ├── view_run.go # 負荷試験実行
│   │   ├── view_version.go # バージョン表示
│   │   └── view_help.go    # ヘルプ表示
│   ├── config/         # 設定管理
│   │   └── config.go
│   └── report/         # レポート生成
│       ├── report.go   # 統計計算
│       ├── html.go     # HTMLレポート
│       └── json.go     # JSONレポート
├── config.yaml         # 設定ファイル例
├── go.mod
├── go.sum
└── README.md
```

### ビルド

```bash
# 開発用ビルド
go build -o meteor-shower ./cmd/meteor-shower

# リリース用ビルド (バージョン情報を埋め込み)
go build -ldflags="-X 'github.com/kitsystemyou/meteor-shower/internal/cli.Version=1.0.0' \
                    -X 'github.com/kitsystemyou/meteor-shower/internal/cli.GitCommit=$(git rev-parse HEAD)' \
                    -X 'github.com/kitsystemyou/meteor-shower/internal/cli.BuildDate=$(date -u +%Y-%m-%dT%H:%M:%SZ)'" \
         -o meteor-shower ./cmd/meteor-shower
```

### テスト

```bash
# すべてのテストを実行
go test ./...

# カバレッジ付きでテスト
go test -cover ./...
```

### 依存関係

このプロジェクトはGo標準パッケージのみを使用しています:

- `flag` - コマンドライン引数のパース
- `net/http` - HTTPリクエスト送信
- `html/template` - HTMLレポート生成
- `encoding/json` - JSONレポート生成
- `time` - タイミング制御と計測
- `sync` - 並行処理制御
- `gopkg.in/yaml.v3` - YAML設定ファイルのパース (準標準ライブラリ)

## パッケージ公開

### 前提条件

`go install` でインストール可能にするには、以下の構造が必要です:

```
リポジトリルート/
├── go.mod
├── cmd/
│   └── meteor-shower/
│       └── main.go
└── internal/
```

### GitHub へのプッシュ

1. 変更をコミット:
```bash
git add -A
git commit -m "Release v1.0.0"
```

2. メインブランチにプッシュ:
```bash
git push origin main
```

3. タグを作成してプッシュ:
```bash
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0
```

**重要**: タグをプッシュする前に、必ず最新のコードを `main` ブランチにプッシュしてください。

### トラブルシューティング

`go install` 実行時に以下のエラーが出る場合:
```
module github.com/kitsystemyou/meteor-shower@latest found (v1.0.1), but does not contain package github.com/kitsystemyou/meteor-shower/cmd/meteor-shower
```

**原因**: タグが古いコミットを指しているか、最新のコードがプッシュされていません。

**解決方法**:
1. 最新のコードをプッシュ:
   ```bash
   git push origin main
   ```

2. 既存のタグを削除して再作成:
   ```bash
   # ローカルのタグを削除
   git tag -d v1.0.1
   
   # リモートのタグを削除
   git push origin :refs/tags/v1.0.1
   
   # 新しいタグを作成
   git tag -a v1.0.1 -m "Release v1.0.1"
   
   # タグをプッシュ
   git push origin v1.0.1
   ```

3. Go のモジュールキャッシュをクリア:
   ```bash
   go clean -modcache
   ```

4. 再度インストール:
   ```bash
   go install github.com/kitsystemyou/meteor-shower/cmd/meteor-shower@latest
   ```

### GitHub Releases でバイナリを公開

#### 1. クロスコンパイルでバイナリをビルド

**簡単な方法: ビルドスクリプトを使用**

```bash
# バージョンを指定してビルド
./scripts/build-release.sh v1.0.0

# または開発版としてビルド
./scripts/build-release.sh
```

**手動でビルドする場合:**

複数のプラットフォーム向けにバイナリをビルドします:

```bash
# バージョン情報を設定
VERSION="1.0.0"
GIT_COMMIT=$(git rev-parse HEAD)
BUILD_DATE=$(date -u +%Y-%m-%dT%H:%M:%SZ)

# ビルド用のディレクトリを作成
mkdir -p dist

# Linux (amd64)
GOOS=linux GOARCH=amd64 go build \
  -ldflags="-X 'github.com/kitsystemyou/meteor-shower/internal/cli.Version=${VERSION}' \
            -X 'github.com/kitsystemyou/meteor-shower/internal/cli.GitCommit=${GIT_COMMIT}' \
            -X 'github.com/kitsystemyou/meteor-shower/internal/cli.BuildDate=${BUILD_DATE}'" \
  -o dist/meteor-shower-linux-amd64 ./cmd/meteor-shower

# Linux (arm64)
GOOS=linux GOARCH=arm64 go build \
  -ldflags="-X 'github.com/kitsystemyou/meteor-shower/internal/cli.Version=${VERSION}' \
            -X 'github.com/kitsystemyou/meteor-shower/internal/cli.GitCommit=${GIT_COMMIT}' \
            -X 'github.com/kitsystemyou/meteor-shower/internal/cli.BuildDate=${BUILD_DATE}'" \
  -o dist/meteor-shower-linux-arm64 ./cmd/meteor-shower

# macOS (amd64)
GOOS=darwin GOARCH=amd64 go build \
  -ldflags="-X 'github.com/kitsystemyou/meteor-shower/internal/cli.Version=${VERSION}' \
            -X 'github.com/kitsystemyou/meteor-shower/internal/cli.GitCommit=${GIT_COMMIT}' \
            -X 'github.com/kitsystemyou/meteor-shower/internal/cli.BuildDate=${BUILD_DATE}'" \
  -o dist/meteor-shower-darwin-amd64 ./cmd/meteor-shower

# macOS (arm64 / Apple Silicon)
GOOS=darwin GOARCH=arm64 go build \
  -ldflags="-X 'github.com/kitsystemyou/meteor-shower/internal/cli.Version=${VERSION}' \
            -X 'github.com/kitsystemyou/meteor-shower/internal/cli.GitCommit=${GIT_COMMIT}' \
            -X 'github.com/kitsystemyou/meteor-shower/internal/cli.BuildDate=${BUILD_DATE}'" \
  -o dist/meteor-shower-darwin-arm64 ./cmd/meteor-shower

# Windows (amd64)
GOOS=windows GOARCH=amd64 go build \
  -ldflags="-X 'github.com/kitsystemyou/meteor-shower/internal/cli.Version=${VERSION}' \
            -X 'github.com/kitsystemyou/meteor-shower/internal/cli.GitCommit=${GIT_COMMIT}' \
            -X 'github.com/kitsystemyou/meteor-shower/internal/cli.BuildDate=${BUILD_DATE}'" \
  -o dist/meteor-shower-windows-amd64.exe ./cmd/meteor-shower
```

#### 2. チェックサムファイルを生成

```bash
cd dist
sha256sum meteor-shower-* > checksums.txt
cd ..
```

#### 3. アーカイブを作成（オプション）

```bash
cd dist

# tar.gz 形式（Linux/macOS）
tar -czf meteor-shower-linux-amd64.tar.gz meteor-shower-linux-amd64
tar -czf meteor-shower-linux-arm64.tar.gz meteor-shower-linux-arm64
tar -czf meteor-shower-darwin-amd64.tar.gz meteor-shower-darwin-amd64
tar -czf meteor-shower-darwin-arm64.tar.gz meteor-shower-darwin-arm64

# zip 形式（Windows）
zip meteor-shower-windows-amd64.zip meteor-shower-windows-amd64.exe

cd ..
```

#### 4. GitHub Release を作成

**方法1: GitHub Web UI を使用**

1. GitHub リポジトリページにアクセス
2. 右側の「Releases」をクリック
3. 「Create a new release」をクリック
4. タグを選択または新規作成（例: `v1.0.0`）
5. リリースタイトルを入力（例: `v1.0.0`）
6. リリースノートを記載:

```markdown
## 新機能
- 複数エンドポイントへの負荷試験対応
- 重み付けによるリクエスト分散
- HTML/JSON形式のレポート出力

## インストール方法

### Go install
```bash
go install github.com/kitsystemyou/meteor-shower/cmd/meteor-shower@v1.0.0
```

### バイナリダウンロード

お使いのプラットフォームに応じたバイナリをダウンロードしてください:

- **Linux (amd64)**: `meteor-shower-linux-amd64.tar.gz`
- **Linux (arm64)**: `meteor-shower-linux-arm64.tar.gz`
- **macOS (Intel)**: `meteor-shower-darwin-amd64.tar.gz`
- **macOS (Apple Silicon)**: `meteor-shower-darwin-arm64.tar.gz`
- **Windows**: `meteor-shower-windows-amd64.zip`

ダウンロード後、展開して実行権限を付与してください:

```bash
# Linux/macOS
tar -xzf meteor-shower-*.tar.gz
chmod +x meteor-shower-*
sudo mv meteor-shower-* /usr/local/bin/meteor-shower

# Windows
# ZIPを展開してPATHに追加
```

### チェックサム検証

```bash
sha256sum -c checksums.txt
```
```

7. 「Attach binaries」セクションにファイルをドラッグ&ドロップ:
   - `meteor-shower-linux-amd64.tar.gz`
   - `meteor-shower-linux-arm64.tar.gz`
   - `meteor-shower-darwin-amd64.tar.gz`
   - `meteor-shower-darwin-arm64.tar.gz`
   - `meteor-shower-windows-amd64.zip`
   - `checksums.txt`

8. 「Publish release」をクリック

**方法2: GitHub CLI を使用**

```bash
# GitHub CLI をインストール（未インストールの場合）
# macOS: brew install gh
# Linux: https://github.com/cli/cli/blob/trunk/docs/install_linux.md

# 認証
gh auth login

# リリースを作成
gh release create v1.0.0 \
  --title "v1.0.0" \
  --notes "Release v1.0.0" \
  dist/meteor-shower-linux-amd64.tar.gz \
  dist/meteor-shower-linux-arm64.tar.gz \
  dist/meteor-shower-darwin-amd64.tar.gz \
  dist/meteor-shower-darwin-arm64.tar.gz \
  dist/meteor-shower-windows-amd64.zip \
  dist/checksums.txt
```

#### 5. リリースの確認

1. GitHub リポジトリの Releases ページで確認
2. バイナリがダウンロード可能か確認
3. インストール手順をテスト:

```bash
# バイナリダウンロードのテスト
wget https://github.com/kitsystemyou/meteor-shower/releases/download/v1.0.0/meteor-shower-linux-amd64.tar.gz
tar -xzf meteor-shower-linux-amd64.tar.gz
./meteor-shower-linux-amd64 version
```

### 自動化（GitHub Actions）

このリポジトリには `.github/workflows/release.yml` が含まれており、タグをプッシュすると自動的にリリースが作成されます。

**使い方:**

1. タグを作成してプッシュ:
```bash
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0
```

2. GitHub Actions が自動的に:
   - 全プラットフォーム向けにビルド
   - アーカイブとチェックサムを生成
   - GitHub Release を作成
   - バイナリをアップロード

3. リリースページで確認:
   https://github.com/kitsystemyou/meteor-shower/releases

**ワークフローの内容:**

```yaml
name: Release

on:
  push:
    tags:
      - 'v*'

jobs:
  release:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      
      - name: Build binaries
        run: |
          VERSION=${GITHUB_REF#refs/tags/}
          GIT_COMMIT=$(git rev-parse HEAD)
          BUILD_DATE=$(date -u +%Y-%m-%dT%H:%M:%SZ)
          
          mkdir -p dist
          
          # Linux amd64
          GOOS=linux GOARCH=amd64 go build \
            -ldflags="-X 'github.com/kitsystemyou/meteor-shower/internal/cli.Version=${VERSION}' \
                      -X 'github.com/kitsystemyou/meteor-shower/internal/cli.GitCommit=${GIT_COMMIT}' \
                      -X 'github.com/kitsystemyou/meteor-shower/internal/cli.BuildDate=${BUILD_DATE}'" \
            -o dist/meteor-shower-linux-amd64 ./cmd/meteor-shower
          
          # その他のプラットフォームも同様にビルド...
      
      - name: Create Release
        uses: softprops/action-gh-release@v1
        with:
          files: dist/*
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

### Go Modules

Go modulesとして公開する場合、ユーザーは以下のコマンドでインストールできます:

```bash
go install github.com/kitsystemyou/meteor-shower/cmd/meteor-shower@latest
```

## ライセンス

[MIT License](https://www.tldrlegal.com/license/mit-license)

## 貢献

プルリクエストを歓迎します。大きな変更の場合は、まずissueを開いて変更内容を議論してください。
