package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/joho/godotenv"
	pb "github.com/m4salah/redroc/grpc/protos"
	"github.com/m4salah/redroc/grpc/storage"
	"github.com/m4salah/redroc/grpc/types"
	"github.com/m4salah/redroc/util"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// release is set through the linker at build time, generally from a git sha.
// Used for logging and error reporting.
var release string

var (
	env           = flag.String("env", "development", "Env")
	listenPort    = flag.Int("listen_port", 8080, "start server on this port")
	storageBucket = flag.String("storage_bucket", "sre-classroom-image-server_photos-3", "storage bucket to use for storing photos")
	storageDryRun = flag.Bool("storage_dry_run", false, "disable storage bucket reads")
)

type DownloadServiceRPC struct {
	pb.UnimplementedDownloadPhotoServer
	types.DownloadService
}

func (d *DownloadServiceRPC) Download(ctx context.Context, request *pb.DownloadPhotoRequest) (*pb.DownloadPhotoResponse, error) {
	if *storageDryRun {
		return &pb.DownloadPhotoResponse{}, nil
	}
	image, err := d.DownloadService.DB.Get(ctx, request.ImgName)
	if err != nil {
		return nil, err
	}
	return &pb.DownloadPhotoResponse{ImgBlob: image}, nil
}

func main() {
	flag.Parse()

	// load env variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	logger, err := util.CreateLogger(*env, release)
	if err != nil {
		fmt.Println("Error setting up the logger:", err)
		return
	}
	bucket, err := storage.NewBuckets(storage.NewBucketsOptions{Log: logger, BucketName: *storageBucket})
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
	downloadService := types.DownloadService{DB: bucket, Log: logger}
	pb.RegisterDownloadPhotoServer(grpcServer, &DownloadServiceRPC{DownloadService: downloadService})

	logger.Info("starting GRPC server", zap.Int("port", *listenPort))
	if err := grpcServer.Serve(listener); err != nil {
		logger.Fatal("Failed to serve", zap.Int("port", *listenPort))
	}
}
