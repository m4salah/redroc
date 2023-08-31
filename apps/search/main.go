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
	env                 = flag.String("env", "development", "Env")
	firestoreLatestPath = flag.String("firestore_latest_path", "latest", "path for storing latest images")
	listenPort          = flag.Int("listen_port", 8080, "start server on this port")
	storageDryRun       = flag.Bool("storage_dry_run", false, "disable storage bucket writes")
	thumbnailCount      = flag.Int("thumbnail_count", 30, "number of thumbnails to return")
	thumbnailPrefix     = flag.String("thumbnail_prefix", "/download/thumbnail_", "name prefix to use for storing thumbnails")
)

type SearchServiceRPC struct {
	pb.UnimplementedGetThumbnailServer
	Log        *zap.Logger
	MetadataDB storage.MetadataDB
}

func (s *SearchServiceRPC) GetThumbnail(ctx context.Context, request *pb.GetThumbnailImagesRequest) (*pb.GetThumbnailImagesResponse, error) {
	if *storageDryRun {
		return &pb.GetThumbnailImagesResponse{}, nil
	}
	urls, err := s.MetadataDB.GetThumbnails(ctx, *thumbnailCount, request.SearchKeyword)
	if err != nil {
		return nil, err
	}
	return &pb.GetThumbnailImagesResponse{StorageUrl: urls}, nil
}

type Config struct {
	FilestoreProject string `mapstructure:"FILESTORE_PROJECT"`
}

func main() {
	flag.Parse()
	logger, err := util.CreateLogger(*env, release)

	if err != nil {
		fmt.Println("Error setting up the logger:", err)
		return
	}

	// load env variables
	config := util.LoadConfig(Config{})

	filestore, err := storage.NewFilestore(storage.NewFilestoreOptions{ProjectID: config.FilestoreProject,
		Log:             logger,
		FilestoreLatest: *firestoreLatestPath,
		ThumbnailPerfix: *thumbnailPrefix,
	})

	if err != nil {
		fmt.Println("Error initilizing filestore", err)
		return
	}
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", *listenPort))
	if err != nil {
		logger.Fatal("failed to listen", zap.Int("port", *listenPort), zap.Error(err))
		return
	}
	grpcServer := grpc.NewServer()
	pb.RegisterGetThumbnailServer(grpcServer, &SearchServiceRPC{MetadataDB: filestore, Log: logger})

	logger.Info("starting GRPC server", zap.Int("port", *listenPort))
	if err := grpcServer.Serve(listener); err != nil {
		logger.Fatal("Failed to serve", zap.Int("port", *listenPort))
	}
}
