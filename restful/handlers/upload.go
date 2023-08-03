package handlers

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"path"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/m4salah/redroc/util"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	pb "github.com/m4salah/redroc/grpc/protos"
)

const (
	MB = 1 << 20
)

func Upload(mux chi.Router, backendAddr string, log *zap.Logger, backendTimeout time.Duration) {
	mux.Post("/upload", func(w http.ResponseWriter, r *http.Request) {
		creds := credentials.NewTLS(&tls.Config{
			InsecureSkipVerify: true,
		})
		conn, err := grpc.Dial(backendAddr, grpc.WithTransportCredentials(creds))
		if err != nil {
			log.Error("Cannot dial to grpc service", zap.Error(err))
			http.Error(w, "Cannot dial download service", http.StatusBadRequest)
			return
		}
		defer conn.Close()

		err = r.ParseMultipartForm(5 * MB)
		if err != nil {
			http.Error(w, "Error parsing form data", http.StatusBadRequest)
			return
		}

		// Get the username
		username := r.FormValue("username")
		if username == "" {
			log.Error("user name must be provided", zap.Error(err))
			http.Error(w, "username must be provided", http.StatusBadRequest)
			return
		}

		// Get the hashtags
		hashtags, err := util.GetTags(r)
		if err != nil {
			log.Error("getting tags failes", zap.Error(err))
			http.Error(w, "Invalid tags", http.StatusBadRequest)
			return

		}

		// Access the array of strings using the field name
		blob, filename, err := util.GetPhoto(r)
		if err != nil {
			log.Error("Getting file failed", zap.Error(err))
			http.Error(w, "Getting file failed", http.StatusBadRequest)
			return
		}

		client := pb.NewUploadPhotoClient(conn)

		ctx, cancel := context.WithTimeout(context.Background(), backendTimeout)
		defer cancel()

		objName := uuid.New().String() + path.Ext(filename)

		eg := new(errgroup.Group)

		eg.Go(func() error {
			uploadRequest := &pb.UploadImageRequest{ObjName: objName, Image: blob}
			_, err = client.Upload(ctx, uploadRequest, grpc.WaitForReady(true))
			if err != nil {
				return fmt.Errorf("photo upload failed: %v", err)
			}
			return nil
		})

		eg.Go(func() error {
			metadataRequest := &pb.CreateMetadataRequest{
				ObjName: objName, User: username, Hashtags: hashtags}
			_, err = client.CreateMetadata(ctx, metadataRequest, grpc.WaitForReady(true))
			if err != nil {
				return fmt.Errorf("metadata create failed: %v", err)
			}
			return nil
		})
		if err := eg.Wait(); err != nil {
			log.Error("Error while uploading or creating metadata", zap.Error(err))
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		fmt.Fprintln(w, "Image Uploaded Successfully")
	})
}
