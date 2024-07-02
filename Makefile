# サーバを起動
start-dev:
	@echo "Starting development server..."
	@go run cmd/server/main.go

# protoからgoファイルを生成
buf-generate:
	@buf generate