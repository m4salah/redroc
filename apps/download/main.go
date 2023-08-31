package main

import (
	"context"
	"flag"
	"fmt"
	"net"

	pb "github.com/m4salah/redroc/libs/proto"
	"github.com/m4salah/redroc/libs/storage"
	"github.com/m4salah/redroc/libs/util"
	"go.uber.org/zap"
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
	Log *zap.Logger
	DB  storage.ObjectDB
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
	logger, err := util.CreateLogger(*env, release)
	if err != nil {
		fmt.Println("Error setting up the logger:", err)
		return
	}
	bucket, err := storage.NewBuckets(storage.NewBucketsOptions{Log: logger, BucketName: config.StorageBucket})
	if err != nil {
		fmt.Println("Error initializing Bucket", err)
		return
	}
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", *listenPort))
	if err != nil {
		logger.Fatal("failed to listen", zap.Int("port", *listenPort), zap.Error(err))
		return
	}
	grpcServer := grpc.NewServer()
	pb.RegisterDownloadPhotoServer(grpcServer, &DownloadServiceRPC{DB: bucket, Log: logger})

	logger.Info("starting GRPC server", zap.Int("port", *listenPort))
	if err := grpcServer.Serve(listener); err != nil {
		logger.Fatal("Failed to serve", zap.Int("port", *listenPort))
	}
}
