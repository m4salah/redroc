# load .env file to get access to GOOGLE_PROJECT_ID
include .env

RELEASE:=$(shell git rev-parse --short HEAD)

proto-download:
	protoc                                      \
			--go_out=.                          \
			--go_opt=paths=source_relative      \
			--go-grpc_out=.                     \
			--go-grpc_opt=paths=source_relative \
			libs/proto/download.proto

proto-upload:
	protoc                                      \
			--go_out=.                          \
			--go_opt=paths=source_relative      \
			--go-grpc_out=.                     \
			--go-grpc_opt=paths=source_relative \
			libs/proto/upload.proto

proto-search:
	protoc                                      \
			--go_out=.                          \
			--go_opt=paths=source_relative      \
			--go-grpc_out=.                     \
			--go-grpc_opt=paths=source_relative \
			libs/proto/search.proto

proto: proto-download proto-upload proto-search

build-download:
	go build -ldflags "-X main.release=$(RELEASE)" -o bin/download apps/download/main.go

run-download: build-download
	cd apps/download && ../../bin/download -listen_port 8081

build-upload:
	go build -ldflags "-X main.release=$(RELEASE)" -o bin/upload apps/upload/main.go

run-upload: build-upload
	cd apps/upload && ../../bin/upload -listen_port 8082

build-search:
	go build -ldflags "-X main.release=$(RELEASE)" -o bin/search apps/search/main.go
	
run-search: build-search
	cd apps/search && ../../bin/search -storage_dry_run true

build-server:
	go build -ldflags "-X main.release=$(RELEASE)" -o bin/server apps/server/cmd/main.go

run-server: build-server
	cd apps/server && ../../bin/server -skip_gcloud_auth true

build-migrator:
	go build -ldflags "-X main.release=$(RELEASE)" -o bin/migrator apps/migrator/main.go

run-migrator-up: build-migrator
	cd apps/migrator && ../../bin/migrator up

run-migrator-down: build-migrator
	cd apps/migrator && ../../bin/migrator down

run-all: 
	make -j 5 run-download run-upload run-search run-server run-frontend

# docker command for server.
docker-build-server:
	docker build --build-arg RELEASE_ARG=$(RELEASE) -t redroc-server -f Dockerfile.server .

docker-run-server: docker-build-server
	docker run -p 8080:8080 redroc-server:latest

docker-tag-server: docker-build-server
	docker tag redroc-server gcr.io/$(GOOGLE_PROJECT_ID)/redroc-server

docker-push-server: docker-tag-server
	docker push gcr.io/$(GOOGLE_PROJECT_ID)/redroc-server

deploy-server: docker-push-server
	gcloud run deploy redroc-server \
  		--image gcr.io/$(GOOGLE_PROJECT_ID)/redroc-server \
		--platform managed \
		--region us-central1  \
		--allow-unauthenticated

# docker command for server-rs.
docker-build-server-rs:
	docker build -t redroc-server-rs -f Dockerfile.server-rs .

docker-run-server-rs: docker-build-server-rs
	docker run -p 8080:8080 redroc-server-rs:latest

docker-tag-server-rs: docker-build-server-rs
	docker tag redroc-server-rs gcr.io/$(GOOGLE_PROJECT_ID)/redroc-server-rs

docker-push-server-rs: docker-tag-server-rs
	docker push gcr.io/$(GOOGLE_PROJECT_ID)/redroc-server-rs

deploy-server-rs: docker-push-server-rs
	gcloud run deploy redroc-server-rs \
  		--image gcr.io/$(GOOGLE_PROJECT_ID)/redroc-server-rs \
		--platform managed \
		--region us-central1  \
		--allow-unauthenticated

# docker command for download.
docker-build-download:
	docker build --build-arg RELEASE_ARG=$(RELEASE) -t redroc-download -f Dockerfile.download-grpc .

docker-run-download: docker-build-download
	docker run -p 8080:8080 redroc-download:latest

docker-tag-download: docker-build-download
	docker tag redroc-download gcr.io/$(GOOGLE_PROJECT_ID)/redroc-download

docker-push-download: docker-tag-download
	docker push gcr.io/$(GOOGLE_PROJECT_ID)/redroc-download

deploy-download: docker-push-download
	gcloud run deploy redroc-download \
  		--image gcr.io/$(GOOGLE_PROJECT_ID)/redroc-download \
		--platform managed \
		--region us-central1  

# docker command for upload.
docker-build-upload:
	docker build --build-arg RELEASE_ARG=$(RELEASE) -t redroc-upload -f Dockerfile.upload-server .

docker-run-upload: docker-build-upload
	docker run -p 8080:8080 redroc-upload:latest

docker-tag-upload: docker-build-upload
	docker tag redroc-upload gcr.io/$(GOOGLE_PROJECT_ID)/redroc-upload

docker-push-upload: docker-tag-upload
	docker push gcr.io/$(GOOGLE_PROJECT_ID)/redroc-upload

deploy-upload: docker-push-upload
	gcloud config set project $(GOOGLE_PROJECT_ID) 
	gcloud run deploy redroc-upload \
  		--image gcr.io/$(GOOGLE_PROJECT_ID)/redroc-upload \
		--platform managed \
		--region us-central1  

# docker command for search.
docker-build-search:
	docker build --build-arg RELEASE_ARG=$(RELEASE) -t redroc-search -f Dockerfile.search-grpc .

docker-run-search: docker-build-search
	docker run -p 8080:8080 redroc-search:latest

docker-tag-search: docker-build-search
	docker tag redroc-search gcr.io/$(GOOGLE_PROJECT_ID)/redroc-search

docker-push-search: docker-tag-search
	docker push gcr.io/$(GOOGLE_PROJECT_ID)/redroc-search

deploy-search: docker-push-search
	gcloud run deploy redroc-search \
  		--image gcr.io/$(GOOGLE_PROJECT_ID)/redroc-search \
		--platform managed \
		--region us-central1  

# Deploy all services.
deploy-all: 
	make -j 4 deploy-server deploy-download deploy-upload deploy-search

# docker command for frontend.
run-frontend:
	cd apps/frontend && npm run dev

build-frontend:
	cd apps/frontend && npm run build

# Go commands.
mod-tidy:
	cd apps/download && go mod tidy
	cd apps/upload && go mod tidy
	cd apps/search && go mod tidy
	cd apps/server && go mod tidy
	cd apps/migrator && go mod tidy
	cd libs/util && go mod tidy
	cd libs/storage && go mod tidy

# rust commands
server-rs:
	cargo run --bin server