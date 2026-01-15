# SkillSync Makefile

# 变量
BINARY_NAME := skillsync
# 版本号：优先使用精确 tag，否则使用 dev+commit
VERSION := $(shell git describe --tags --exact-match 2>/dev/null || echo "dev")
BUILD_DIR := build
GO := go

# Git 信息
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')

# 编译标志
LDFLAGS := -ldflags "-s -w \
	-X github.com/AlfonsSkills/SkillSync/cmd.Version=$(VERSION) \
	-X github.com/AlfonsSkills/SkillSync/cmd.GitCommit=$(GIT_COMMIT) \
	-X github.com/AlfonsSkills/SkillSync/cmd.BuildTime=$(BUILD_TIME)"

# 目标平台
PLATFORMS := darwin/amd64 darwin/arm64 linux/amd64 linux/arm64 windows/amd64

.PHONY: all build clean test lint help cross

# 默认目标
all: build

# 编译当前平台到 build 目录
build:
	@echo "🔨 Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) .
	@echo "✅ Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

# 清理构建产物
clean:
	@echo "🧹 Cleaning..."
	rm -rf $(BUILD_DIR)
	@echo "✅ Clean complete"

# 运行测试
test:
	@echo "🧪 Running tests..."
	$(GO) test -v ./...

# 代码检查
lint:
	@echo "🔍 Running linter..."
	golangci-lint run ./...

# 跨平台编译
cross:
	@echo "🌍 Cross-compiling for multiple platforms..."
	@mkdir -p $(BUILD_DIR)
	@for platform in $(PLATFORMS); do \
		GOOS=$${platform%/*} GOARCH=$${platform#*/} \
		$(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-$${platform%/*}-$${platform#*/}$(if $(findstring windows,$${platform%/*}),.exe,) . ; \
		echo "  ✓ $(BUILD_DIR)/$(BINARY_NAME)-$${platform%/*}-$${platform#*/}" ; \
	done
	@echo "✅ Cross-compile complete"

# 快速运行
run: build
	$(BUILD_DIR)/$(BINARY_NAME) --help

# 帮助信息
help:
	@echo "SkillSync Makefile"
	@echo ""
	@echo "Usage:"
	@echo "  make build  - Build for current platform (output: build/)"
	@echo "  make clean  - Clean build artifacts"
	@echo "  make test   - Run tests"
	@echo "  make lint   - Run linter"
	@echo "  make cross  - Cross-compile for all platforms"
	@echo "  make run    - Build and run"
	@echo "  make help   - Show this help"
