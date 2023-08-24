GOOGLE_PROJECT_ID?=carbon-relic-393513
RELEASE:=$(shell git rev-parse --short HEAD)

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
	go build -ldflags "-X main.release=$(RELEASE)" -o bin/download grpc/services/download/main.go

run-download: build-download
	./bin/download -listen_port 8081

build-upload:
	go build -ldflags "-X main.release=$(RELEASE)" -o bin/upload grpc/services/upload/main.go

run-upload: build-upload
	./bin/upload -listen_port 8082

build-search:
	go build -ldflags "-X main.release=$(RELEASE)" -o bin/search grpc/services/search/main.go
	
run-search: build-search
	./bin/search -listen_port 8083

build-server:
	go build -ldflags "-X main.release=$(RELEASE)" -o bin/server restful/cmd/main.go

run-server: build-server
	./bin/server -listen_port 8080               \
		-download_backend_address localhost:8081 \
		-upload_backend_address   localhost:8082 \
		-search_backend_address   localhost:8083 \
		-skip_gcloud_auth 		  true

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
	docker build --build-arg RELEASE_ARG=$(RELEASE) -t redroc-upload -f Dockerfile.upload-grpc .

docker-run-upload: docker-build-upload
	docker run -p 8080:8080 redroc-upload:latest

docker-tag-upload: docker-build-upload
	docker tag redroc-upload gcr.io/$(GOOGLE_PROJECT_ID)/redroc-upload

docker-push-upload: docker-tag-upload
	docker push gcr.io/$(GOOGLE_PROJECT_ID)/redroc-upload

deploy-upload: docker-push-upload
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
deploy-all: deploy-server deploy-download deploy-upload deploy-search deploy-frontend

# docker command for frontend.
run-frontend:
	cd frontend && npm run dev

build-frontend:
	cd frontend && npm run build
