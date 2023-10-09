FROM golang:1.21-bullseye as build

WORKDIR /
COPY ./apps/server ./
COPY ./libs ./libs

# Set the release variable at build time
ARG RELEASE_ARG
ENV RELEASE=$RELEASE_ARG

RUN	CGO_ENABLED=0 go build -ldflags="-s -w -X main.release=$RELEASE" -o bin/server cmd/main.go

FROM gcr.io/distroless/static-debian11

COPY --from=build /bin/server /bin/server

ENTRYPOINT ["/bin/server"]
