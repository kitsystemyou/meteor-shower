# mycli

A sample CLI tool built with Go that demonstrates best practices for command-line applications.

## 前提条件

- Go 1.21 以上
- Git (オプション: ソースからビルドする場合)

## インストール

### ソースからビルド

```bash
git clone https://github.com/example/mycli.git
cd mycli
go build -o mycli ./cmd/mycli
```

### Go install を使用

```bash
go install github.com/example/mycli/cmd/mycli@latest
```

## 使い方

### 基本的な使い方

```bash
# デフォルト設定で実行
mycli run

# 名前を指定して実行
mycli run Alice

# カスタムメッセージを指定
mycli run Bob --message "Good morning"
```

### グローバルフラグ

すべてのコマンドで使用できるフラグ:

| フラグ | 短縮形 | デフォルト | 説明 |
|--------|--------|------------|------|
| `--config` | - | `./config.yaml` | 設定ファイルのパス |
| `--output` | `-o` | `text` | 出力形式 (text, json) |
| `--verbose` | `-v` | `false` | 詳細な出力を表示 |

### サブコマンド

#### `run` - メインロジックを実行

アプリケーションのメインロジックを実行します。

```bash
mycli run [name] [flags]
```

**引数:**
- `name` (オプション): 挨拶に使用する名前

**フラグ:**
- `-m, --message string`: カスタムメッセージを指定

**例:**

```bash
# 基本的な実行
mycli run

# 名前を指定
mycli run Alice

# JSON形式で出力
mycli run Bob --output json

# 詳細モードで実行
mycli run Charlie --verbose

# カスタムメッセージ
mycli run Dave --message "Good evening"

# 設定ファイルを指定
mycli run --config /path/to/config.yaml
```

#### `version` - バージョン情報を表示

CLIツールのバージョン情報を表示します。

```bash
mycli version
```

#### `help` - ヘルプを表示

コマンドのヘルプ情報を表示します。

```bash
# 全体のヘルプ
mycli help

# 特定のコマンドのヘルプ
mycli help run
mycli run --help
```

## 設定ファイル

`config.yaml` ファイルで設定を管理できます。

### 設定ファイルの場所

以下の順序で設定ファイルを検索します:

1. `--config` フラグで指定されたパス
2. カレントディレクトリの `config.yaml`
3. `$HOME/.mycli/config.yaml`

### 設定例

```yaml
app:
  # 挨拶に使用する名前
  name: "World"
  
  # メッセージのプレフィックス
  message: "Hello"
  
  # タイムアウト (秒)
  timeout: 30
  
  # デバッグモードを有効化
  debug: false
```

### 設定項目

| 項目 | 型 | デフォルト | 説明 |
|------|-----|-----------|------|
| `app.name` | string | `"World"` | 挨拶に使用するデフォルトの名前 |
| `app.message` | string | `"Hello"` | メッセージのプレフィックス |
| `app.timeout` | int | `30` | タイムアウト時間 (秒) |
| `app.debug` | bool | `false` | デバッグモードの有効化 |

## 環境変数

設定は環境変数でも上書きできます:

```bash
export APP_NAME="Alice"
export APP_MESSAGE="Hi"
export APP_TIMEOUT=60
export APP_DEBUG=true

mycli run
```

## 開発

### プロジェクト構造

```
.
├── cmd/
│   └── mycli/          # メインエントリーポイント
│       └── main.go
├── internal/
│   ├── cmd/            # コマンド実装
│   │   ├── root.go     # ルートコマンド
│   │   ├── run.go      # runサブコマンド
│   │   └── version.go  # versionサブコマンド
│   └── config/         # 設定管理
│       └── config.go
├── config.yaml         # 設定ファイル例
├── go.mod
├── go.sum
└── README.md
```

### ビルド

```bash
# 開発用ビルド
go build -o mycli ./cmd/mycli

# リリース用ビルド (バージョン情報を埋め込み)
go build -ldflags="-X 'github.com/example/mycli/internal/cmd.Version=1.0.0' \
                    -X 'github.com/example/mycli/internal/cmd.GitCommit=$(git rev-parse HEAD)' \
                    -X 'github.com/example/mycli/internal/cmd.BuildDate=$(date -u +%Y-%m-%dT%H:%M:%SZ)'" \
         -o mycli ./cmd/mycli
```

### テスト

```bash
# すべてのテストを実行
go test ./...

# カバレッジ付きでテスト
go test -cover ./...
```

### 依存関係

このプロジェクトは以下のライブラリを使用しています:

- [cobra](https://github.com/spf13/cobra) - CLIフレームワーク
- [viper](https://github.com/spf13/viper) - 設定管理

## パッケージ公開

### GitHub Releases

1. タグを作成してプッシュ:
```bash
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0
```

2. GitHub Releasesでバイナリを公開

### Go Modules

Go modulesとして公開する場合、ユーザーは以下のコマンドでインストールできます:

```bash
go install github.com/example/mycli/cmd/mycli@latest
```

## ライセンス

[MIT License](https://www.tldrlegal.com/license/mit-license)

## 貢献

プルリクエストを歓迎します。大きな変更の場合は、まずissueを開いて変更内容を議論してください。
