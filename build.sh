#!/bin/bash

# Monitor Agent Build Script

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 配置
BINARY_NAME="monitor-agent"
VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(date -u '+%Y-%m-%d_%H:%M:%S')
GO_VERSION=$(go version | awk '{print $3}')

# 输出信息
info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1"
    exit 1
}

# 检查Go环境
check_go() {
    if ! command -v go &> /dev/null; then
        error "Go is not installed. Please install Go first."
    fi
    info "Go version: $GO_VERSION"
}

# 清理旧文件
clean() {
    info "Cleaning old build files..."
    rm -f "$BINARY_NAME"
    rm -f "${BINARY_NAME}.pid"
    info "Clean completed"
}

# 下载依赖
deps() {
    info "Downloading dependencies..."
    go mod download
    go mod tidy
    info "Dependencies ready"
}

# 代码检查
lint() {
    info "Running code checks..."

    # go fmt
    if ! gofmt -l . | grep -q .; then
        info "✓ Code formatting OK"
    else
        warn "Code formatting issues found, auto-fixing..."
        gofmt -w .
    fi

    # go vet
    if go vet ./...; then
        info "✓ Go vet passed"
    else
        error "Go vet failed"
    fi
}

# 编译
build() {
    info "Building $BINARY_NAME..."

    # 构建标志
    LDFLAGS="-s -w"
    LDFLAGS="$LDFLAGS -X 'main.Version=$VERSION'"
    LDFLAGS="$LDFLAGS -X 'main.BuildTime=$BUILD_TIME'"

    # 编译
    go build -ldflags "$LDFLAGS" -o "$BINARY_NAME" .

    if [ -f "$BINARY_NAME" ]; then
        chmod +x "$BINARY_NAME"
        FILE_SIZE=$(du -h "$BINARY_NAME" | awk '{print $1}')
        info "✓ Build completed: $BINARY_NAME ($FILE_SIZE)"
        info "  Version: $VERSION"
        info "  Build Time: $BUILD_TIME"

        # 自动重启服务
        info "Auto-restarting service..."
        ./"$BINARY_NAME" restart
    else
        error "Build failed"
    fi
}

# 交叉编译
build_cross() {
    local os=$1
    local arch=$2
    local output="${BINARY_NAME}-${os}-${arch}"

    if [ "$os" = "windows" ]; then
        output="${output}.exe"
    fi

    info "Building for ${os}/${arch}..."

    LDFLAGS="-s -w -X 'main.Version=$VERSION' -X 'main.BuildTime=$BUILD_TIME'"

    GOOS=$os GOARCH=$arch go build -ldflags "$LDFLAGS" -o "$output" .

    if [ -f "$output" ]; then
        FILE_SIZE=$(du -h "$output" | awk '{print $1}')
        info "✓ Built: $output ($FILE_SIZE)"
    else
        error "Cross build failed for ${os}/${arch}"
    fi
}

# 多平台编译
build_all() {
    info "Building for multiple platforms..."

    # Linux
    build_cross linux amd64
    build_cross linux arm64

    # macOS
    build_cross darwin amd64
    build_cross darwin arm64

    # Windows
    build_cross windows amd64

    info "All platforms built successfully"
}

# 测试
test() {
    info "Running tests..."
    go test -v ./...
    info "Tests completed"
}

# 运行
run() {
    info "Starting $BINARY_NAME..."
    ./"$BINARY_NAME" "$@"
}

# 安装
install() {
    local install_path=${1:-/usr/local/bin}

    if [ ! -f "$BINARY_NAME" ]; then
        error "$BINARY_NAME not found. Run build first."
    fi

    info "Installing $BINARY_NAME to $install_path..."

    if [ ! -w "$install_path" ]; then
        sudo cp "$BINARY_NAME" "$install_path/"
    else
        cp "$BINARY_NAME" "$install_path/"
    fi

    info "✓ Installed to $install_path/$BINARY_NAME"
}

# 显示帮助
show_help() {
    cat << EOF
Monitor Agent Build Script

Usage: $0 [command]

Commands:
    build       Build the binary and auto-restart service (default)
    clean       Clean build artifacts
    deps        Download dependencies
    lint        Run code checks (fmt, vet)
    test        Run tests
    all         Clean + Deps + Lint + Build + Restart
    cross       Build for multiple platforms
    run         Build and run the application
    install     Install binary to system (default: /usr/local/bin)
    help        Show this help message

Examples:
    $0                      # Build and restart
    $0 all                  # Full build process
    $0 cross                # Cross-platform build
    $0 install              # Install to /usr/local/bin
    $0 install /opt/bin     # Install to custom path
    $0 run start            # Build and start daemon

EOF
}

# 主函数
main() {
    case "${1:-build}" in
        build)
            check_go
            build
            ;;
        clean)
            clean
            ;;
        deps)
            check_go
            deps
            ;;
        lint)
            check_go
            lint
            ;;
        test)
            check_go
            test
            ;;
        all)
            check_go
            clean
            deps
            lint
            build
            ;;
        cross)
            check_go
            deps
            build_all
            ;;
        run)
            check_go
            build
            shift
            run "$@"
            ;;
        install)
            shift
            install "$@"
            ;;
        help|--help|-h)
            show_help
            ;;
        *)
            error "Unknown command: $1. Use 'help' for usage."
            ;;
    esac
}

main "$@"