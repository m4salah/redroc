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

var config Config

type DownloadServiceRPC struct {
	pb.UnimplementedDownloadPhotoServer
	DB storage.ObjectDB
}

func (d *DownloadServiceRPC) Download(ctx context.Context, request *pb.DownloadPhotoRequest) (*pb.DownloadPhotoResponse, error) {
	if *storageDryRun {
		return &pb.DownloadPhotoResponse{}, nil
	}
	// get the encrypted image from the datastore
	encryptedImage, err := d.DB.Get(ctx, request.ImgName)
	if err != nil {
		slog.Error("error while getting the image", "error", err)
		return nil, err
	}

	// decrypt the image
	decryptedImage, err := util.DecryptAES(encryptedImage, []byte(config.EncryptionKey))
	if err != nil {
		slog.Error("error while decrypting the image", "error", err)
		return nil, err
	}
	return &pb.DownloadPhotoResponse{ImgBlob: decryptedImage}, nil
}

type Config struct {
	EncryptionKey string `env:"ENCRYPTION_KEY,notEmpty"`
	StorageBucket string `env:"STORAGE_BUCKET,notEmpty"`
}

func main() {
	flag.Parse()

	// load env variables
	if err := util.LoadConfig(&config); err != nil {
		panic(err)
	}

	util.InitializeSlog(*env, release)

	bucket, err := storage.NewBuckets(storage.NewBucketsOptions{BucketName: config.StorageBucket})
	if err != nil {
		slog.Error("failed to initialize the bucket", slog.String("bucketName", config.StorageBucket))
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
