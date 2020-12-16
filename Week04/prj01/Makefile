server:
	go build -o server  ./cmd/service

client:
	go build -o client  ./cmd/client

pb:
	protoc api/hello.proto --go_out=plugins=grpc:api

wire:
	wire ./cmd/service/rpc.go

.PHONY:server client pb wire