package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"log/slog"
	"math/rand"
	"net/http"
	"path"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
	"github.com/m4salah/redroc/libs/pubsub"
	"github.com/m4salah/redroc/libs/storage"
	"github.com/m4salah/redroc/libs/util"
	"golang.org/x/sync/errgroup"
)

// release is set through the linker at build time, generally from a git sha.
// Used for logging and error reporting.
var release string

var (
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
	DB         storage.ObjectDB
	MetadataDB storage.MetadataDB
}

type Config struct {
	EncryptionKey    string `env:"ENCRYPTION_KEY,notEmpty"`
	SockerUri        string `env:"SOCKET_URI,notEmpty"`
	FilestoreProject string `env:"FILESTORE_PROJECT,notEmpty"`
	StorageBucket    string `env:"STORAGE_BUCKET,notEmpty"`
	Env              string `env:"ENV,notEmpty"`
}

var config Config

func (d *UploadServiceRPC) Upload(ctx context.Context, request pubsub.UploadMessage) error {
	// check if we are in dry run mode
	if *storageDryRun {
		return nil
	}
	thumb, err := util.MakeThumbnail(request.Message.Data, *thumbnailWidth, *thumbnailHeight)
	if err != nil {
		slog.Error("error making the thumbnail", "error", err)
		return err
	}

	eg := new(errgroup.Group)

	// encrypt the image and the thumbnail
	encryptedImage, err := util.EncryptAES(request.Message.Data, []byte(config.EncryptionKey))
	if err != nil {
		slog.Error("error encrypting the image", "error", err)
		return err
	}

	encryptedThumb, err := util.EncryptAES(thumb, []byte(config.EncryptionKey))
	if err != nil {
		slog.Error("error encrypting the thumbnail", "error", err)
		return err
	}

	// Store the original image
	eg.Go(func() error {
		return d.DB.Store(ctx, request.Message.Attributes.ObjName, encryptedImage)
	})

	// Store the thumbnail image
	eg.Go(func() error {
		return d.DB.Store(ctx, *thumbnailPrefix+request.Message.Attributes.ObjName, encryptedThumb)
	})

	// check if either the operation failed
	if err := eg.Wait(); err != nil {
		slog.Error("error while uploading the image and the thumbnail", "error", err)
		return err
	}

	return nil
}

func (d *UploadServiceRPC) CreateMetadata(ctx context.Context, request pubsub.UploadMessage) error {

	// check if we are in dry run mode
	if *firestoreDryRun {
		return nil
	}

	// unmarshal hashtags from string to array of strings
	var hashtags []string
	if err := json.Unmarshal([]byte(request.Message.Attributes.Hashtags), &hashtags); err != nil {
		slog.Error("error while parsing hashtags", slog.String("error", err.Error()))
		return err
	}

	eg := new(errgroup.Group)

	timestamp := time.Now().Unix()
	eg.Go(func() error {
		id := path.Join(request.Message.Attributes.User, request.Message.Attributes.ObjName)
		return d.MetadataDB.StorePath(ctx, id, timestamp)
	})

	eg.Go(func() error {
		var failure error
		for _, tag := range hashtags {
			id := path.Join(tag, request.Message.Attributes.ObjName)
			err := d.MetadataDB.StorePathWithUser(ctx, request.Message.Attributes.User, id, timestamp)
			if err != nil {
				failure = err
				break
			}
		}
		return failure
	})

	eg.Go(func() error {
		index := atomic.AddUint32(&latestIdxFirestore, 1) % uint32(*firestoreLatestPhotos)
		return d.MetadataDB.StoreLatest(ctx, index, *firestoreLatestPath, request.Message.Attributes.ObjName)
	})

	if err := eg.Wait(); err != nil {
		slog.Error("error creating the metadata", "error", err)
		return err
	}
	return nil
}

func (d *UploadServiceRPC) ImageUploaded(ctx context.Context, request pubsub.UploadMessage) error {

	slog.Info("upload image triggered",
		slog.String("imageNmae", request.Message.Attributes.ObjName),
		slog.String("username", request.Message.Attributes.User),
		slog.Any("hashtags", request.Message.Attributes.Hashtags),
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

	return nil
}

func main() {
	flag.Parse()

	// load env variables
	if err := util.LoadConfig(&config); err != nil {
		panic(err)
	}

	util.InitializeSlog(config.Env, release)

	// pubsub code
	http.HandleFunc("/upload", uploadImage)

	slog.Info("starting server", slog.Int("port", *listenPort))
	if err := http.ListenAndServe(fmt.Sprintf(":%d", *listenPort), nil); err != nil {
		log.Fatal(err)
	}
}

// uploadImage receives and processes a Pub/Sub push message.
func uploadImage(w http.ResponseWriter, r *http.Request) {
	var m pubsub.UploadMessage
	body, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("ioutil.ReadAll", "error", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// byte slice unmarshalling handles base64 decoding.
	if err := json.Unmarshal(body, &m); err != nil {
		slog.Error("json.Unmarshal", "error", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	slog.Info("Message Info",
		slog.Any("attriutes", m.Message.Attributes),
		slog.String("id", m.Message.ID),
		slog.String("subscription", m.Subscription))

	bucket, err := storage.NewBuckets(storage.NewBucketsOptions{BucketName: config.StorageBucket})
	if err != nil {
		slog.Error("error initializing Bucket", "error", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	filestore, err := storage.NewFilestore(storage.NewFilestoreOptions{
		ProjectID:       config.FilestoreProject,
		FilestoreLatest: *firestoreLatestPath,
		ThumbnailPerfix: *thumbnailPrefix})

	if err != nil {
		slog.Error("error initilizing filestore ", "error", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	uploadService := &UploadServiceRPC{
		DB:         bucket,
		MetadataDB: filestore,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	eg := new(errgroup.Group)
	eg.Go(func() error {
		return uploadService.CreateMetadata(ctx, m)
	})

	eg.Go(func() error {
		return uploadService.Upload(ctx, m)
	})

	if err := eg.Wait(); err != nil {
		slog.Error("error creating the metadata or uploading the image", "error", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	// trigger image uploaded event in different goroutine
	go uploadService.ImageUploaded(ctx, m)

	slog.Info("image uploaded successfully",
		slog.Any("attriutes", m.Message.Attributes),
		slog.String("id", m.Message.ID),
		slog.String("subscription", m.Subscription))

	w.WriteHeader(http.StatusOK)
}
