#!/usr/bin/env bash
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o release/yudai-mac main.go
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o release/yudai-mac-arm64 main.go
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o release/yudai-linux main.go
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o release/yudai-linux-arm64 main.go
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o release/yudai-windows.exe main.go
CGO_ENABLED=0 GOOS=windows GOARCH=arm64 go build -ldflags="-s -w" -o release/yudai-windows-arm64.exe main.go