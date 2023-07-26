proto-download:
	protoc                                      \
			--go_out=.                          \
			--go_opt=paths=source_relative      \
			--go-grpc_out=.                     \
			--go-grpc_opt=paths=source_relative \
			grpc/protos/download.proto

build-download:
	go build -o bin/download grpc/services/download/main.go

build-server:
	go build -o bin/server restful/cmd/server/main.go