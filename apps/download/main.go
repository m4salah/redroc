package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net"

	pb "github.com/m4salah/redroc/libs/proto"
	"github.com/m4salah/redroc/libs/storage"
	"github.com/m4salah/redroc/libs/util"
	"google.golang.org/grpc"
)

// release is set through the linker at build time, generally from a git sha.
// Used for logging and error reporting.
var release string

var (
	env           = flag.String("env", "development", "Env")
	listenPort    = flag.Int("listen_port", 8080, "start server on this port")
	storageDryRun = flag.Bool("storage_dry_run", false, "disable storage bucket reads")
)

type DownloadServiceRPC struct {
	pb.UnimplementedDownloadPhotoServer
	DB storage.ObjectDB
}

func (d *DownloadServiceRPC) Download(ctx context.Context, request *pb.DownloadPhotoRequest) (*pb.DownloadPhotoResponse, error) {
	if *storageDryRun {
		return &pb.DownloadPhotoResponse{}, nil
	}
	image, err := d.DB.Get(ctx, request.ImgName)
	if err != nil {
		return nil, err
	}
	return &pb.DownloadPhotoResponse{ImgBlob: image}, nil
}

type Config struct {
	EncryptionKey string `mapstructure:"ENCRYPTION_KEY"`
	StorageBucket string `mapstructure:"STORAGE_BUCKET"`
}

func main() {
	flag.Parse()

	config := util.LoadConfig(Config{})

	// load env variables
	util.InitializeSlog(*env, release)

	bucket, err := storage.NewBuckets(storage.NewBucketsOptions{BucketName: config.StorageBucket})
	if err != nil {
		fmt.Println("Error initializing Bucket", err)
		return
	}
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", *listenPort))
	if err != nil {
		slog.Error("failed to listen", slog.Int("port", *listenPort), slog.String("error", err.Error()))
		return
	}
	grpcServer := grpc.NewServer()
	pb.RegisterDownloadPhotoServer(grpcServer, &DownloadServiceRPC{DB: bucket})

	slog.Info("starting GRPC server", slog.Int("port", *listenPort))
	if err := grpcServer.Serve(listener); err != nil {
		slog.Error("Failed to serve", slog.Int("port", *listenPort))
	}
}
