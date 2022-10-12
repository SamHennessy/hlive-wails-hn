# Development Helpers

.PHONY: build
build:
	wails build

.PHONY: build-win
build-win:
	wails build -platform windows/amd64

# Bug: Browser will need to be manually reload after each rebuild
.PHONY: dev
dev:
	wails dev

.PHONY: run-mac
run-mac:
	"build/bin/HLive Hacker News.app/Contents/MacOS/HLive Hacker News"

# format code in a way the linter will be happy
.PHONY: format
format:
	go mod tidy
	goimports -l -w .
	gofumpt -l -w .
