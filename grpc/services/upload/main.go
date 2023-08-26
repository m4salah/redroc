package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"path"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	pb "github.com/m4salah/redroc/grpc/protos"
	"github.com/m4salah/redroc/grpc/storage"
	"github.com/m4salah/redroc/grpc/types"
	"github.com/m4salah/redroc/util"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

// release is set through the linker at build time, generally from a git sha.
// Used for logging and error reporting.
var release string

var (
	env                   = flag.String("env", "development", "Env")
	firestoreLatestPath   = flag.String("firestore_latest_path", "latest", "path for storing latest images")
	firestoreLatestPhotos = flag.Uint("firestore_latest_photos", 10, "number of latest images to store")
	firestoreProject      = flag.String("firestore_project", "carbon-relic-393513", "firestore project to use for storing tags")
	firestoreDryRun       = flag.Bool("firestore_dry_run", false, "disable firestore writes")
	listenPort            = flag.Int("listen_port", 8082, "start server on this port")
	storageBucket         = flag.String("storage_bucket", "sre-classroom-image-server_photos-3", "storage bucket to use for storing photos")
	storageDryRun         = flag.Bool("storage_dry_run", false, "disable storage bucket writes")
	thumbnailHeight       = flag.Uint("thumbnail_height", 180, "height of the generated photo thumbnail")
	thumbnailPrefix       = flag.String("thumbnail_prefix", "thumbnail_", "name prefix to use for storing thumbnails")
	socketUri             = flag.String("socket_uri", "ws://localhost:8080/ws", "address of the socket server")
	thumbnailWidth        = flag.Uint("thumbnail_width", 320, "width of the generated photo thumbnail")
	latestIdxFirestore    = rand.Uint32()
)

type UploadServiceRPC struct {
	pb.UnimplementedUploadPhotoServer
	types.UploadService
}

func (d *UploadServiceRPC) Upload(ctx context.Context, request *pb.UploadImageRequest) (*pb.UploadImageResponse, error) {

	// check if we are in dry run mode
	if *storageDryRun {
		return &pb.UploadImageResponse{}, nil
	}
	thumb, err := util.MakeThumbnail(request.Image, *thumbnailWidth, *thumbnailHeight)
	if err != nil {
		return nil, err
	}

	eg := new(errgroup.Group)

	// Store the original image
	eg.Go(func() error {
		return d.DB.Store(ctx, request.ObjName, request.Image)
	})

	// Store the thumbnail image
	eg.Go(func() error {
		return d.DB.Store(ctx, *thumbnailPrefix+request.ObjName, thumb)
	})

	// check if either the operation failed
	if err := eg.Wait(); err != nil {
		return nil, err
	}

	// TODO: Refactor this into own struct
	// boradcast the new image to all connected clients
	c, _, err := websocket.DefaultDialer.Dial(*socketUri, nil)
	if err != nil {
		log.Println("Error connecting:", zap.Error(err))
	}
	defer c.Close()

	// Send a message to the server
	err = c.WriteMessage(websocket.TextMessage, []byte("new image"))
	if err != nil {
		log.Println("Error writing message:", zap.Error(err))
	}

	log.Println("end of the upload")
	return &pb.UploadImageResponse{}, nil
}

func (d *UploadServiceRPC) CreateMetadata(ctx context.Context, request *pb.CreateMetadataRequest) (*pb.CreateMetadataResponse, error) {

	// check if we are in dry run mode
	if *firestoreDryRun {
		return &pb.CreateMetadataResponse{}, nil
	}

	eg := new(errgroup.Group)

	timestamp := time.Now().Unix()
	eg.Go(func() error {
		id := path.Join(request.User, request.ObjName)
		return d.MetadataDB.StorePath(ctx, id, timestamp)
	})

	eg.Go(func() error {
		var failure error
		for _, tag := range request.Hashtags {
			id := path.Join(tag, request.ObjName)
			err := d.MetadataDB.StorePathWithUser(ctx, request.User, id, timestamp)
			if err != nil {
				failure = err
				break
			}
		}
		return failure
	})

	eg.Go(func() error {
		index := atomic.AddUint32(&latestIdxFirestore, 1) % uint32(*firestoreLatestPhotos)
		return d.MetadataDB.StoreLatest(ctx, index, *firestoreLatestPath, request.ObjName)
	})

	if err := eg.Wait(); err != nil {
		return nil, err
	}
	return &pb.CreateMetadataResponse{}, nil
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
	filestore, err := storage.NewFilestore(storage.NewFilestoreOptions{ProjectID: *firestoreProject,
		Log:             logger,
		FilestoreLatest: *firestoreLatestPath,
		ThumbnailPerfix: *thumbnailPrefix})
	if err != nil {
		fmt.Println("Error initilizing filestore ", err)
		return
	}
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", *listenPort))
	if err != nil {
		logger.Fatal("failed to listen", zap.Int("port", *listenPort), zap.Error(err))
		return
	}
	grpcServer := grpc.NewServer()
	uploadService := types.UploadService{DB: bucket, Log: logger, MetadataDB: filestore}
	pb.RegisterUploadPhotoServer(grpcServer, &UploadServiceRPC{UploadService: uploadService})

	logger.Info("starting GRPC server", zap.Int("port", *listenPort))
	if err := grpcServer.Serve(listener); err != nil {
		logger.Fatal("Failed to serve", zap.Int("port", *listenPort))
	}
}
