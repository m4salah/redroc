proto-download:
	protoc                                      \
			--go_out=.                          \
			--go_opt=paths=source_relative      \
			--go-grpc_out=.                     \
			--go-grpc_opt=paths=source_relative \
			grpc/protos/download.proto

proto-upload:
	protoc                                      \
			--go_out=.                          \
			--go_opt=paths=source_relative      \
			--go-grpc_out=.                     \
			--go-grpc_opt=paths=source_relative \
			grpc/protos/upload.proto

proto-search:
	protoc                                      \
			--go_out=.                          \
			--go_opt=paths=source_relative      \
			--go-grpc_out=.                     \
			--go-grpc_opt=paths=source_relative \
			grpc/protos/search.proto

proto: proto-download proto-upload proto-search

build-download:
	go build -o bin/download grpc/services/download/main.go

build-upload:
	go build -o bin/upload grpc/services/upload/main.go

build-search:
	go build -o bin/search grpc/services/search/main.go
	
build-server:
	go build -o bin/server restful/cmd/server/main.go