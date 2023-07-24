proto:
	protoc ./cmd/api/*.proto --go-grpc_out=pkg --go_out=pkg