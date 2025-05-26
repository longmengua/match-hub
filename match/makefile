# Makefile 中的目標如果和檔案名稱相同，就會被當作檔案依賴進行比對。
# .PHONY 宣告表示這是一個純命令執行目標，而非檔案產出，因此每次都會執行。
# 避免因檔案快取或更新時間問題而跳過執行。
.PHONY: dev init proto help

# 啟動應用程式
dev:  
	go run cmd/main.go

# 安裝依賴
init:  
	go mod tidy

# 建置
build:  
	go build -ldflags "-X main.gitCommit=$(git rev-parse --short HEAD)" -o match

# 容器化建置
docker: 
	docker build -t match:latest .

# 生成 Protobuf Go 代碼
proto:  
	for file in ../proto/*.proto; do \
		proto_name=$$(basename "$$file"); \
		protoc --proto_path=../proto \
			--go_out=./proto --go_opt=paths=source_relative --go_opt=M$$proto_name=match/proto \
			--go-grpc_out=./proto --go-grpc_opt=paths=source_relative --go-grpc_opt=M$$proto_name=match/proto \
			"$$file"; \
	done

# 顯示幫助訊息
help:  
	@echo ''
	@echo 'Usage:'
	@echo ' make [target]'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
	helpMessage = match(lastLine, /^# (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 2, RLENGTH); \
			printf "\033[36m%-22s\033[0m %s\n", helpCommand, helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help
