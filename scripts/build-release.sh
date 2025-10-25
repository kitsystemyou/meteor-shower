#!/bin/bash

set -e

# バージョン情報
VERSION=${1:-"dev"}
GIT_COMMIT=$(git rev-parse HEAD)
BUILD_DATE=$(date -u +%Y-%m-%dT%H:%M:%SZ)

echo "Building meteor-shower ${VERSION}"
echo "Git commit: ${GIT_COMMIT}"
echo "Build date: ${BUILD_DATE}"
echo ""

# ビルド用のディレクトリを作成
mkdir -p dist

# ビルド関数
build() {
    local GOOS=$1
    local GOARCH=$2
    local OUTPUT=$3
    
    echo "Building for ${GOOS}/${GOARCH}..."
    
    GOOS=${GOOS} GOARCH=${GOARCH} go build \
        -ldflags="-X 'github.com/kitsystemyou/meteor-shower/internal/cli.Version=${VERSION}' \
                  -X 'github.com/kitsystemyou/meteor-shower/internal/cli.GitCommit=${GIT_COMMIT}' \
                  -X 'github.com/kitsystemyou/meteor-shower/internal/cli.BuildDate=${BUILD_DATE}'" \
        -o ${OUTPUT} ./cmd/meteor-shower
    
    echo "✓ Built ${OUTPUT}"
}

# 各プラットフォーム向けにビルド
build linux amd64 dist/meteor-shower-linux-amd64
build linux arm64 dist/meteor-shower-linux-arm64
build darwin amd64 dist/meteor-shower-darwin-amd64
build darwin arm64 dist/meteor-shower-darwin-arm64
build windows amd64 dist/meteor-shower-windows-amd64.exe

echo ""
echo "Creating archives..."

cd dist

# tar.gz 形式（Linux/macOS）
tar -czf meteor-shower-linux-amd64.tar.gz meteor-shower-linux-amd64
tar -czf meteor-shower-linux-arm64.tar.gz meteor-shower-linux-arm64
tar -czf meteor-shower-darwin-amd64.tar.gz meteor-shower-darwin-amd64
tar -czf meteor-shower-darwin-arm64.tar.gz meteor-shower-darwin-arm64

# zip 形式（Windows）
if command -v zip &> /dev/null; then
    zip meteor-shower-windows-amd64.zip meteor-shower-windows-amd64.exe
else
    echo "Warning: zip command not found, skipping Windows archive"
fi

# チェックサムを生成
echo ""
echo "Generating checksums..."
sha256sum meteor-shower-* > checksums.txt

cd ..

echo ""
echo "✓ Build complete!"
echo ""
echo "Files in dist/:"
ls -lh dist/

echo ""
echo "To create a GitHub release, run:"
echo "  gh release create ${VERSION} --title \"${VERSION}\" --notes \"Release ${VERSION}\" dist/*.tar.gz dist/*.zip dist/checksums.txt"
