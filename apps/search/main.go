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
	firestoreLatestPath = flag.String("firestore_latest_path", "latest", "path for storing latest images")
	storageDryRun       = flag.Bool("storage_dry_run", false, "disable storage bucket writes")
	thumbnailCount      = flag.Int("thumbnail_count", 30, "number of thumbnails to return")
	thumbnailPrefix     = flag.String("thumbnail_prefix", "/download/thumbnail_", "name prefix to use for storing thumbnails")
)

type SearchServiceRPC struct {
	pb.UnimplementedGetThumbnailServer
	MetadataDB storage.MetadataDB
}

func (s *SearchServiceRPC) GetThumbnail(ctx context.Context, request *pb.GetThumbnailImagesRequest) (*pb.GetThumbnailImagesResponse, error) {
	slog.Info("receiving search request with params", slog.String("q", request.SearchKeyword))
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
	FilestoreProject string `env:"FILESTORE_PROJECT,notEmpty"`
	Env              string `env:"ENV,notEmpty"`
	Port             int    `env:"PORT,notEmpty"`
}

var config Config

func main() {
	flag.Parse()
	util.InitializeSlog(config.Env, release)

	// load env variables
	if err := util.LoadConfig(&config); err != nil {
		panic(err)
	}

	filestore, err := storage.NewFilestore(storage.NewFilestoreOptions{
		ProjectID:       config.FilestoreProject,
		FilestoreLatest: *firestoreLatestPath,
		ThumbnailPerfix: *thumbnailPrefix,
	})

	if err != nil {
		fmt.Println("Error initilizing filestore", err)
		return
	}
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", config.Port))
	if err != nil {
		slog.Error("failed to listen", slog.Int("port", config.Port), "error", err)
		return
	}
	grpcServer := grpc.NewServer()
	pb.RegisterGetThumbnailServer(grpcServer, &SearchServiceRPC{MetadataDB: filestore})

	slog.Info("starting GRPC server", slog.Int("port", config.Port))
	if err := grpcServer.Serve(listener); err != nil {
		slog.Error("Failed to serve", slog.Int("port", config.Port))
	}
}
