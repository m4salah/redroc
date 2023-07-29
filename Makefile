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
	go build -o bin/server restful/cmd/main.go

docker-build-server:
	docker build -t redroc-server -f Dockerfile.server .

docker-run-server: docker-build-server
	docker run -p 8080:8080 redroc-server:latest

docker-build-download:
	docker build -t redroc-download -f Dockerfile.download-grpc .

docker-run-download: docker-build-download
	docker run -p 8080:8080 redroc-download:latest

docker-build-upload:
	docker build -t redroc-upload -f Dockerfile.upload-grpc .

docker-run-upload: docker-build-upload
	docker run -p 8080:8080 redroc-upload:latest

docker-build-search:
	docker build -t redroc-search -f Dockerfile.search-grpc .

docker-run-search: docker-build-search
	docker run -p 8080:8080 redroc-search:latest