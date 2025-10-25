# Workload Test Server

負荷試験用のテストサーバー。Go標準パッケージのみで実装されています。

## 概要

このサーバーは、負荷試験ツール (mycli) のテスト用に設計されています。
レスポンス遅延、エラー率、ランダム遅延などを設定可能で、様々な負荷試験シナリオをシミュレートできます。

## 使い方

### 基本的な起動

```bash
go run main.go
```

デフォルトでポート8080で起動し、10msの遅延でレスポンスを返します。

### オプション

| フラグ | デフォルト | 説明 |
|--------|-----------|------|
| `-port` | `8080` | リスニングポート |
| `-delay` | `10` | レスポンス遅延 (ミリ秒) |
| `-error-rate` | `0.0` | エラー率 (0.0 ～ 1.0) |
| `-random-delay` | `false` | ランダム遅延を有効化 (±50%の変動) |

### 起動例

```bash
# ポート9000で起動
go run main.go -port 9000

# 50msの遅延で起動
go run main.go -delay 50

# 10%のエラー率で起動
go run main.go -error-rate 0.1

# ランダム遅延を有効化
go run main.go -random-delay

# すべてのオプションを指定
go run main.go -port 9000 -delay 30 -error-rate 0.05 -random-delay
```

## エンドポイント

### `GET /`

通常のエンドポイント。設定された遅延とエラー率でレスポンスを返します。

**レスポンス例:**
```json
{
  "status": "ok",
  "timestamp": "2024-01-01T12:00:00Z",
  "delay_ms": 10
}
```

### `GET /health`

ヘルスチェックエンドポイント。遅延なしで即座にレスポンスを返します。

**レスポンス例:**
```json
{
  "status": "healthy"
}
```

### `GET /stats`

サーバー統計情報を返します。

**レスポンス例:**
```json
{
  "total_requests": 1234,
  "uptime": "1h23m45s",
  "start_time": "2024-01-01T10:00:00Z"
}
```

### `GET /slow`

意図的に遅いエンドポイント。常に500msの遅延でレスポンスを返します。

**レスポンス例:**
```json
{
  "status": "ok",
  "timestamp": "2024-01-01T12:00:00Z",
  "delay_ms": 500
}
```

### `GET /error`

常にエラーを返すエンドポイント。ステータスコード500を返します。

## 使用例

### 1. 基本的な負荷試験

```bash
# サーバー起動
go run main.go

# 別のターミナルで負荷試験実行
cd ..
./mycli run --rps 10 --concurrency 2 -o html > report.html
```

### 2. 高負荷シナリオ

```bash
# サーバー起動 (低遅延)
go run main.go -delay 5

# 高RPSで負荷試験
cd ..
./mycli run --rps 100 --concurrency 10 -o json > report.json
```

### 3. エラーハンドリングテスト

```bash
# サーバー起動 (20%エラー率)
go run main.go -error-rate 0.2

# 負荷試験実行
cd ..
./mycli run --rps 20 --concurrency 5 -o html > report.html
```

### 4. 遅延バリエーションテスト

```bash
# サーバー起動 (ランダム遅延)
go run main.go -delay 50 -random-delay

# 負荷試験実行
cd ..
./mycli run --rps 15 --concurrency 3 -o html > report.html
```

### 5. 遅いエンドポイントのテスト

設定ファイルを作成:

```yaml
# slow-test.yaml
loadtest:
  domain: "http://localhost:8080"
  endpoint: "/slow"
  rps: 5
  concurrency: 2
  duration: 10
  output: "html"
```

実行:

```bash
# サーバー起動
go run main.go

# 負荷試験実行
cd ..
./mycli run --config slow-test.yaml > report.html
```

### 6. 複数エンドポイントの負荷試験

設定ファイルを作成:

```yaml
# multi-endpoint-test.yaml
loadtest:
  domain: "http://localhost:8080"
  endpoints:
    - path: "/"
      weight: 1.0
    - path: "/health"
      weight: 0.5
    - path: "/slow"
      weight: 0.2
  rps: 20
  concurrency: 5
  duration: 10
  output: "html"
```

実行:

```bash
# サーバー起動
go run main.go

# 負荷試験実行
cd ..
./mycli run --config multi-endpoint-test.yaml > report.html
```

レポートには各エンドポイントへのリクエスト分散状況が表示されます。
重みに応じて、`/` が最も多く、`/health` が中程度、`/slow` が最も少ないリクエスト数になります。

## 動作確認

サーバーが正常に起動しているか確認:

```bash
# ヘルスチェック
curl http://localhost:8080/health

# 統計情報確認
curl http://localhost:8080/stats

# 通常エンドポイント
curl http://localhost:8080/
```

## ビルド

バイナリとしてビルドする場合:

```bash
go build -o test-server main.go
./test-server -port 8080 -delay 10
```

## 注意事項

- このサーバーは負荷試験用のテストサーバーです。本番環境では使用しないでください。
- 高負荷時はシステムリソースを消費する可能性があります。
- `-error-rate` を1.0に設定すると、すべてのリクエストがエラーを返します。
