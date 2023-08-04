package handlers

import (
	"context"
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

	pb "github.com/m4salah/redroc/grpc/protos"
)

const (
	MB = 1 << 20
)

// pingRequestWithAuth mints a new Identity Token for each request.
// This token has a 1 hour expiry and should be reused.
// audience must be the auto-assigned URL of a Cloud Run service or HTTP Cloud Function without port number.
func pingUploadRequestWithAuth(backendTimeout time.Duration, backendAddr string, log *zap.Logger, p *pb.UploadImageRequest, audience string) (*pb.UploadImageResponse, error) {

	creds, err := util.CreateTransportCredentials()
	if err != nil {
		log.Fatal("failed to load system root CA cert pool")
	}

	conn, err := grpc.Dial(backendAddr, grpc.WithTransportCredentials(creds))

	if err != nil {
		log.Error("Cannot dial to grpc service", zap.Error(err))
		return nil, fmt.Errorf("grpc.Dial: %w", err)
	}

	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), backendTimeout)
	defer cancel()

	ctx, err = util.GetAuthContext(ctx, audience)
	if err != nil {
		return nil, fmt.Errorf("error get auth context: %w", err)
	}

	// Send the request.
	client := pb.NewUploadPhotoClient(conn)
	return client.Upload(ctx, p, grpc.WaitForReady(true))
}

func pingCreateMetadataRequestWithAuth(backendTimeout time.Duration, backendAddr string, log *zap.Logger, p *pb.CreateMetadataRequest, audience string) (*pb.CreateMetadataResponse, error) {

	creds, err := util.CreateTransportCredentials()
	if err != nil {
		log.Fatal("failed to load system root CA cert pool")
	}

	conn, err := grpc.Dial(backendAddr, grpc.WithTransportCredentials(creds))

	if err != nil {
		log.Error("Cannot dial to grpc service", zap.Error(err))
		return nil, fmt.Errorf("grpc.Dial: %w", err)
	}

	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), backendTimeout)
	defer cancel()

	ctx, err = util.GetAuthContext(ctx, audience)
	if err != nil {
		return nil, fmt.Errorf("error get auth context: %w", err)
	}

	// Send the request.
	client := pb.NewUploadPhotoClient(conn)
	return client.CreateMetadata(ctx, p, grpc.WaitForReady(true))
}

func Upload(mux chi.Router, backendAddr string, log *zap.Logger, backendTimeout time.Duration) {
	mux.Post("/upload", func(w http.ResponseWriter, r *http.Request) {

		err := r.ParseMultipartForm(5 * MB)
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

		objName := uuid.New().String() + path.Ext(filename)

		eg := new(errgroup.Group)

		eg.Go(func() error {
			uploadRequest := &pb.UploadImageRequest{ObjName: objName, Image: blob}
			_, err = pingUploadRequestWithAuth(backendTimeout, backendAddr, log, uploadRequest, util.ExtractServiceURL(backendAddr))
			if err != nil {
				return fmt.Errorf("photo upload failed: %v", err)
			}
			return nil
		})

		eg.Go(func() error {
			metadataRequest := &pb.CreateMetadataRequest{
				ObjName: objName, User: username, Hashtags: hashtags}
			_, err = pingCreateMetadataRequestWithAuth(backendTimeout, backendAddr, log, metadataRequest, util.ExtractServiceURL(backendAddr))
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
