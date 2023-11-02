package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"math/rand"
	"net"
	"path"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
	pb "github.com/m4salah/redroc/libs/proto"
	"github.com/m4salah/redroc/libs/storage"
	"github.com/m4salah/redroc/libs/util"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

// release is set through the linker at build time, generally from a git sha.
// Used for logging and error reporting.
var release string

var (
	env                   = flag.String("env", "development", "Env")
	firestoreLatestPath   = flag.String("firestore_latest_path", "latest", "path for storing latest images")
	firestoreLatestPhotos = flag.Uint("firestore_latest_photos", 30, "number of latest images to store")
	firestoreDryRun       = flag.Bool("firestore_dry_run", false, "disable firestore writes")
	listenPort            = flag.Int("listen_port", 8080, "start server on this port")
	storageDryRun         = flag.Bool("storage_dry_run", false, "disable storage bucket writes")
	thumbnailHeight       = flag.Uint("thumbnail_height", 180, "height of the generated photo thumbnail")
	thumbnailPrefix       = flag.String("thumbnail_prefix", "thumbnail_", "name prefix to use for storing thumbnails")
	thumbnailWidth        = flag.Uint("thumbnail_width", 320, "width of the generated photo thumbnail")
	latestIdxFirestore    = rand.Uint32()
)

type UploadServiceRPC struct {
	pb.UnimplementedUploadPhotoServer
	DB         storage.ObjectDB
	MetadataDB storage.MetadataDB
}

type Config struct {
	EncryptionKey    string `mapstructure:"ENCRYPTION_KEY"`
	SockerUri        string `mapstructure:"SOCKET_URI"`
	FilestoreProject string `mapstructure:"FILESTORE_PROJECT"`
	StorageBucket    string `mapstructure:"STORAGE_BUCKET"`
}

var config Config

func (d *UploadServiceRPC) Upload(ctx context.Context, request *pb.UploadImageRequest) (*pb.UploadImageResponse, error) {

	// check if we are in dry run mode
	if *storageDryRun {
		return &pb.UploadImageResponse{}, nil
	}
	thumb, err := util.MakeThumbnail(request.Image, *thumbnailWidth, *thumbnailHeight)
	if err != nil {
		slog.Error("error making the thumbnail", "error", err)
		return nil, err
	}

	eg := new(errgroup.Group)

	// encrypt the image and the thumbnail
	encryptedImage, err := util.EncryptAES(request.Image, []byte(config.EncryptionKey))
	if err != nil {
		slog.Error("error encrypting the image", "error", err)
		return nil, err
	}

	encryptedThumb, err := util.EncryptAES(thumb, []byte(config.EncryptionKey))
	if err != nil {
		slog.Error("error encrypting the thumbnail", "error", err)
		return nil, err
	}

	// Store the original image
	eg.Go(func() error {
		return d.DB.Store(ctx, request.ObjName, encryptedImage)
	})

	// Store the thumbnail image
	eg.Go(func() error {
		return d.DB.Store(ctx, *thumbnailPrefix+request.ObjName, encryptedThumb)
	})

	// check if either the operation failed
	if err := eg.Wait(); err != nil {
		slog.Error("error while uploading the image and the thumbnail", "error", err)
		return nil, err
	}

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
		slog.Error("error creating the metadata", "error", err)
		return nil, err
	}
	return &pb.CreateMetadataResponse{}, nil
}

func (d *UploadServiceRPC) ImageUploaded(ctx context.Context, request *pb.ImageUploadedRequest) (*pb.ImageUploadedResponse, error) {

	slog.Info("upload image triggered",
		slog.String("imageNmae", request.ObjName),
		slog.String("username", request.User),
		slog.Any("hashtags", request.Hashtags),
	)

	// TODO: Refactor this into own struct
	// boradcast the new image to all connected clients
	c, _, err := websocket.DefaultDialer.Dial(config.SockerUri, nil)
	if err != nil {
		slog.Error("error connecting to the websocket server", slog.String("error", err.Error()))
	} else {
		// Send a message to the server
		// TODO: refactor the message to it's own constant
		err = c.WriteMessage(websocket.TextMessage, []byte("new image"))
		if err != nil {
			slog.Error("error writing message:", slog.String("error", err.Error()))
		}
	}
	defer c.Close()

	return &pb.ImageUploadedResponse{}, nil
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
		fmt.Println("Error initializing Bucket", err)
		return
	}
	filestore, err := storage.NewFilestore(storage.NewFilestoreOptions{ProjectID: config.FilestoreProject,
		FilestoreLatest: *firestoreLatestPath,
		ThumbnailPerfix: *thumbnailPrefix})
	if err != nil {
		fmt.Println("Error initilizing filestore ", err)
		return
	}
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", *listenPort))
	if err != nil {
		slog.Error("failed to listen", slog.Int("port", *listenPort), slog.String("error", err.Error()))
		return
	}
	grpcServer := grpc.NewServer()
	pb.RegisterUploadPhotoServer(grpcServer, &UploadServiceRPC{DB: bucket, MetadataDB: filestore})

	slog.Info("starting GRPC server", slog.Int("port", *listenPort))
	if err := grpcServer.Serve(listener); err != nil {
		slog.Error("Failed to serve", slog.Int("port", *listenPort))
	}
}
