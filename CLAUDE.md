# meteor-shower 開発ガイドライン

このドキュメントは、meteor-showerプロジェクトの開発方針、アーキテクチャ、およびコーディング規約をまとめたものです。

## プロジェクト概要

meteor-showerは、Go標準パッケージのみで構築された負荷試験ツールです。シンプルさと移植性を重視し、外部依存を最小限に抑えています。

## 開発方針

### 1. 依存関係の制限

**原則**: Go標準パッケージのみを使用する

**例外**: 
- `gopkg.in/yaml.v3` - YAML設定ファイルの読み書きに使用
  - 標準パッケージにYAMLサポートがないため、唯一の外部依存として許可

**理由**:
- シンプルさの維持
- セキュリティリスクの最小化
- ビルドの高速化
- 長期的なメンテナンス性の向上

### 2. エラーハンドリング

**原則**: すべてのエラーは適切に処理し、ユーザーに明確なメッセージを提供する

**実装パターン**:
```go
// main.go - エントリーポイント
func main() {
    app := cli.New(os.Args[1:])
    if err := app.Run(); err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }
}

// 各コマンド - エラーを返す
func (c *CLI) someCommand(args []string) error {
    if err := doSomething(); err != nil {
        return fmt.Errorf("failed to do something: %w", err)
    }
    return nil
}
```

**重要事項**:
- エラーメッセージは必ず `os.Stderr` に出力
- エラー発生時は終了コード `1` で終了
- エラーメッセージは具体的で実行可能な情報を含める
- エラーラッピングには `%w` を使用してコンテキストを保持

### 3. コマンドライン引数とフラグ

**原則**: 標準パッケージ `flag` を使用し、一貫性のあるインターフェースを提供

**実装パターン**:
```go
fs := flag.NewFlagSet("command", flag.ContinueOnError)
fs.StringVar(&output, "o", "config.yaml", "output filename")
fs.BoolVar(&force, "f", false, "force overwrite")
```

**規約**:
- 短いフラグ名（`-o`, `-f`）と長いフラグ名（`--output`, `--force`）の両方をサポート
- デフォルト値は常に指定
- ヘルプメッセージは簡潔で明確に

### 4. 設定ファイル

**形式**: YAML

**構造**:
```yaml
domain: example.com
endpoints:
  - path: /api/users
    weight: 0.7
  - path: /api/posts
    weight: 0.3
rps: 100
concurrency: 10
duration: 60
output: html
```

**重要事項**:
- `weight` が指定されていない場合は均等分散（1.0）
- すべてのフィールドにデフォルト値を設定
- 設定ファイルとコマンドライン引数の両方をサポート（引数が優先）

## ディレクトリ構造

```
meteor-shower/
├── cmd/
│   └── meteor-shower/
│       └── main.go              # エントリーポイント
├── internal/
│   ├── cli/
│   │   ├── cli.go               # CLIルーティング
│   │   ├── run.go               # 負荷試験実行
│   │   ├── config_init.go       # 設定ファイル生成
│   │   ├── view_help.go         # ヘルプ表示
│   │   └── view_version.go      # バージョン表示
│   ├── config/
│   │   └── config.go            # 設定構造体とロード処理
│   └── report/
│       ├── report.go            # レポート構造体
│       ├── html.go              # HTMLレポート生成
│       └── json.go              # JSONレポート生成
├── scripts/
│   └── build-release.sh         # クロスコンパイルスクリプト
├── workload_test_server/
│   └── main.go                  # テスト用サーバー
├── .github/
│   └── workflows/
│       └── release.yml          # リリース自動化
├── go.mod
├── go.sum
├── README.md
└── CLAUDE.md                    # このファイル
```

### ディレクトリの役割

#### `cmd/meteor-shower/`
- アプリケーションのエントリーポイント
- 最小限のコード（エラーハンドリングのみ）
- ビジネスロジックは含めない

#### `internal/cli/`
- CLIコマンドの実装
- コマンドルーティング
- フラグパース
- 各コマンドは独立したファイルに分離

#### `internal/config/`
- 設定ファイルの読み込み
- 設定構造体の定義
- デフォルト値の管理

#### `internal/report/`
- レポート生成ロジック
- HTML/JSON形式のサポート
- テンプレートの管理

#### `scripts/`
- ビルドスクリプト
- クロスコンパイル設定
- リリース準備

#### `workload_test_server/`
- 開発・テスト用のHTTPサーバー
- 遅延やエラー率を設定可能
- 本番コードには含まれない

## コーディング規約

### 1. パッケージ構成

**原則**: 機能ごとにパッケージを分離

```go
// 良い例
package config

type LoadTestConfig struct {
    Domain      string
    Endpoints   []Endpoint
    // ...
}

func Load(filename string) (*LoadTestConfig, error) {
    // ...
}
```

```go
// 悪い例 - すべてを main パッケージに詰め込む
package main

type Config struct { /* ... */ }
func LoadConfig() { /* ... */ }
func RunTest() { /* ... */ }
func GenerateReport() { /* ... */ }
```

### 2. エクスポート

**原則**: 必要最小限の型と関数のみをエクスポート

- パッケージ外から使用される型・関数のみ大文字で開始
- 内部実装の詳細は小文字で開始（非エクスポート）

### 3. エラーメッセージ

**原則**: 具体的で実行可能な情報を提供

```go
// 良い例
return fmt.Errorf("file %s already exists. Use -f to overwrite", filename)

// 悪い例
return fmt.Errorf("file exists")
```

### 4. コメント

**原則**: 「なぜ」を説明し、「何を」は説明しない

```go
// 良い例
// Use weighted random selection to distribute requests according to endpoint weights
endpoint := selectEndpoint(config.Endpoints)

// 悪い例
// Select an endpoint
endpoint := selectEndpoint(config.Endpoints)
```

**エクスポートされた型・関数**:
- 必ずドキュメントコメントを記述
- パッケージ名で始める

```go
// LoadTestConfig represents the configuration for a load test.
type LoadTestConfig struct {
    // ...
}

// Load reads and parses a YAML configuration file.
func Load(filename string) (*LoadTestConfig, error) {
    // ...
}
```

### 5. テスト

**原則**: 標準の `testing` パッケージを使用

```go
func TestConfigLoad(t *testing.T) {
    config, err := Load("testdata/config.yaml")
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if config.Domain != "example.com" {
        t.Errorf("expected domain 'example.com', got '%s'", config.Domain)
    }
}
```

## 重み付けエンドポイント選択

### アルゴリズム

重み付けランダム選択を使用して、設定された重みに応じてエンドポイントを分散:

```go
func selectEndpoint(endpoints []config.Endpoint) string {
    if len(endpoints) == 1 {
        return endpoints[0].Path
    }

    totalWeight := 0.0
    for _, ep := range endpoints {
        weight := ep.Weight
        if weight == 0 {
            weight = 1.0
        }
        totalWeight += weight
    }

    r := rand.Float64() * totalWeight
    cumulative := 0.0
    for _, ep := range endpoints {
        weight := ep.Weight
        if weight == 0 {
            weight = 1.0
        }
        cumulative += weight
        if r <= cumulative {
            return ep.Path
        }
    }

    return endpoints[len(endpoints)-1].Path
}
```

### 重要事項

- `weight` が `0` または未指定の場合は `1.0` として扱う
- 累積重みを使用してランダム選択
- 最後のエンドポイントをフォールバックとして使用

## ビルドとリリース

### ローカルビルド

```bash
go build -o meteor-shower ./cmd/meteor-shower
```

### クロスコンパイル

```bash
./scripts/build-release.sh
```

対応プラットフォーム:
- Linux (amd64, arm64)
- macOS (amd64, arm64)
- Windows (amd64, arm64)

### リリースプロセス

1. バージョンタグを作成: `git tag v1.0.0`
2. タグをプッシュ: `git push origin v1.0.0`
3. GitHub Actionsが自動的にビルドとリリースを実行

## レポート生成

### HTML レポート

- `html/template` パッケージを使用
- インラインCSS（外部依存なし）
- レスポンシブデザイン
- 統計情報の視覚化

### JSON レポート

- `encoding/json` パッケージを使用
- 構造化データ
- 他のツールとの連携が容易

## 開発時の注意事項

### 1. 標準パッケージの活用

以下の標準パッケージを積極的に活用:

- `flag` - コマンドライン引数
- `net/http` - HTTPクライアント
- `html/template` - HTMLレポート生成
- `encoding/json` - JSONレポート生成
- `time` - タイミング制御
- `sync` - 並行処理
- `os` - ファイル操作
- `fmt` - フォーマット出力

### 2. 並行処理

- `sync.WaitGroup` でゴルーチンの完了を待機
- チャネルで結果を収集
- `sync.Mutex` で共有データを保護

```go
var wg sync.WaitGroup
resultsChan := make(chan Result, totalRequests)

for i := 0; i < concurrency; i++ {
    wg.Add(1)
    go func() {
        defer wg.Done()
        // 処理
        resultsChan <- result
    }()
}

wg.Wait()
close(resultsChan)
```

### 3. パフォーマンス

- 不要なメモリアロケーションを避ける
- バッファ付きチャネルを使用
- プリアロケーションを活用

```go
// 良い例
results := make([]Result, 0, totalRequests)

// 悪い例
var results []Result
```

### 4. エラーハンドリングのベストプラクティス

- エラーは無視しない
- エラーをラップしてコンテキストを追加
- センチネルエラーは使用しない（標準パッケージのみの方針に従う）

```go
if err := doSomething(); err != nil {
    return fmt.Errorf("failed to do something: %w", err)
}
```

## 今後の拡張

### 追加を検討する機能

- [ ] プロメテウス形式のメトリクス出力
- [ ] リアルタイム進捗表示
- [ ] カスタムHTTPヘッダーのサポート
- [ ] リクエストボディのテンプレート
- [ ] 認証サポート（Basic, Bearer）

### 拡張時の原則

1. 標準パッケージのみを使用する原則を維持
2. 既存のコード構造に従う
3. 後方互換性を保つ
4. ドキュメントを更新する

## まとめ

meteor-showerは、シンプルさと実用性のバランスを重視したツールです。開発時は以下の原則を常に意識してください:

1. **標準パッケージのみ** - 外部依存を最小限に
2. **明確なエラーメッセージ** - ユーザーが問題を理解できるように
3. **一貫性のあるコード** - 既存のパターンに従う
4. **適切なテスト** - 機能の正確性を保証
5. **ドキュメントの更新** - 変更を必ず文書化

これらの原則に従うことで、長期的にメンテナンス可能で信頼性の高いツールを維持できます。
