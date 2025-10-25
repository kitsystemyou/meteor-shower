# meteor-shower

Go標準パッケージのみで構築された負荷試験ツール。

## 前提条件

- Go 1.21 以上
- Git (オプション: ソースからビルドする場合)

## インストール

### ソースからビルド

```bash
git clone https://github.com/example/meteor-shower.git
cd meteor-shower
go build -o meteor-shower ./cmd/meteor-shower
```

### Go install を使用

```bash
go install github.com/example/meteor-shower/cmd/meteor-shower@latest
```

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
go build -ldflags="-X 'github.com/example/meteor-shower/internal/cli.Version=1.0.0' \
                    -X 'github.com/example/meteor-shower/internal/cli.GitCommit=$(git rev-parse HEAD)' \
                    -X 'github.com/example/meteor-shower/internal/cli.BuildDate=$(date -u +%Y-%m-%dT%H:%M:%SZ)'" \
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
go install github.com/example/meteor-shower/cmd/meteor-shower@latest
```

## ライセンス

[MIT License](https://www.tldrlegal.com/license/mit-license)

## 貢献

プルリクエストを歓迎します。大きな変更の場合は、まずissueを開いて変更内容を議論してください。
