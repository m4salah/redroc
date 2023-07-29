package main

import (
	"context"
	"flag"
	"fmt"
	"net"

	pb "github.com/m4salah/redroc/grpc/protos"
	"github.com/m4salah/redroc/grpc/storage"
	"github.com/m4salah/redroc/grpc/types"
	"github.com/m4salah/redroc/util"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var (
	env                 = flag.String("env", "development", "Env")
	firestoreLatestPath = flag.String("firestore_latest_path", "latest", "path for storing latest images")
	firestoreProject    = flag.String("firestore_project", "carbon-relic-393513", "firestore project to use for storing tags")
	listenPort          = flag.String("listen_port", ":8080", "start server on this port")
	storageDryRun       = flag.Bool("storage_dry_run", false, "disable storage bucket writes")
	thumbnailCount      = flag.Int("thumbnail_count", 10, "number of thumbnails to return")
	thumbnailPrefix     = flag.String("thumbnail_prefix", "/download/thumbnail_", "name prefix to use for storing thumbnails")
)

type SearchServiceRPC struct {
	pb.UnimplementedGetThumbnailServer
	types.SearchService
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

func main() {
	flag.Parse()

	logger, err := util.CreateLogger(*env)
	if err != nil {
		fmt.Println("Error setting up the logger:", err)
		return
	}
	filestore, err := storage.NewFilestore(storage.NewFilestoreOptions{ProjectID: *firestoreProject,
		Log:             logger,
		FilestoreLatest: *firestoreLatestPath,
		ThumbnailPerfix: *thumbnailPrefix,
	})

	if err != nil {
		fmt.Println("Error initilizing filestore", err)
		return
	}
	listener, err := net.Listen("tcp", *listenPort)
	if err != nil {
		logger.Fatal("failed to listen", zap.String("port", *listenPort), zap.Error(err))
		return
	}
	grpcServer := grpc.NewServer()
	searchService := types.SearchService{Log: logger, MetadataDB: filestore}
	pb.RegisterGetThumbnailServer(grpcServer, &SearchServiceRPC{SearchService: searchService})

	logger.Info("starting GRPC server", zap.String("port", *listenPort))
	if err := grpcServer.Serve(listener); err != nil {
		logger.Fatal("Failed to serve", zap.String("port", *listenPort))
	}
}
