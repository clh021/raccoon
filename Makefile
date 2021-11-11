.PHONY: generate test serve

generate:
	@go mod tidy
	@go generate ./...
	@echo "[OK] Generate all completed!"

security:
	@gosec ./...
	@echo "[OK] Go security check was completed!"

gitTime=$(shell date +00%y%m%d%H%M%S)
gitCID=$(shell git rev-parse HEAD)
# gitTime=$(git log -1 --format=%at | xargs -I{} date -d @{} +%Y%m%d_%H%M%S)
# 去除 符号信息 和 调试信息  -ldflags="-s -w"
build: generate
	@cd cmd;go build -ldflags "-s -w -X main.build=${gitTime}.${gitCID}" -o "../bin/raccoon"
	@echo "[OK] App binary was created!"

buildcross: generate
	@cd cmd;CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -ldflags "-s -w -X main.build=${gitTime}.${gitCID}" -o "../bin/raccoon.amd64"
	@cp ./bin/raccoon.amd64 ./bin/raccoon.x86_64
	@echo "[OK] App amd64 binary was created!"
	@cd cmd;CGO_ENABLED=0 GOARCH=arm64 GOOS=linux go build -ldflags "-s -w -X main.build=${gitTime}.${gitCID}" -o "../bin/raccoon.arm64"
	@cp ./bin/raccoon.arm64 ./bin/raccoon.aarch64
	@echo "[OK] App arm64 binary was created!"
	@cd cmd;CGO_ENABLED=0 GOARCH=mips64le GOOS=linux go build -ldflags "-s -w -X main.build=${gitTime}.${gitCID}" -o "../bin/raccoon.mips64le"
	@echo "[OK] App mips64le binary was created!"

# 另有 https://golang.org/doc/install/gccgo 压缩方式
# 使用 upx 压缩 体积 `pacman -S upx`
compress:
	@upx -9 ./bin/raccoon

run:
	@./bin/raccoon

test: 
	go test -v ./...

serve: build run
