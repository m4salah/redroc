package handlers

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"path"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/m4salah/redroc/apps/server/types"
	"github.com/m4salah/redroc/libs/util"
	"google.golang.org/grpc"

	pb "github.com/m4salah/redroc/libs/proto"
)

const (
	_4MB = 4 * 1024 * 1024
)

// pingRequestWithAuth mints a new Identity Token for each request.
// This token has a 1 hour expiry and should be reused.
// audience must be the auto-assigned URL of a Cloud Run service or HTTP Cloud Function without port number.
func pingUploadRequestWithAuth(backendTimeout time.Duration,
	backendAddr string,
	p *pb.UploadImageRequest,
	audience string,
	skipAuth bool) (*pb.UploadImageResponse, error) {

	creds, err := util.CreateTransportCredentials(skipAuth)
	if err != nil {
		slog.Error("failed to load system root CA cert pool")
	}

	conn, err := grpc.Dial(backendAddr, grpc.WithTransportCredentials(creds))

	if err != nil {
		slog.Error("cannot dial to grpc service", slog.String("error", err.Error()))
		return nil, fmt.Errorf("grpc.Dial: %w", err)
	}

	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), backendTimeout)
	defer cancel()

	ctx, err = util.GetAuthContext(ctx, audience, skipAuth)
	if err != nil {
		return nil, fmt.Errorf("error get auth context: %w", err)
	}

	// Send the request.
	client := pb.NewUploadPhotoClient(conn)
	return client.Upload(ctx, p, grpc.WaitForReady(true))
}

func pingCreateMetadataRequestWithAuth(backendTimeout time.Duration,
	backendAddr string,
	p *pb.CreateMetadataRequest,
	audience string,
	skipAuth bool) (*pb.CreateMetadataResponse, error) {

	creds, err := util.CreateTransportCredentials(skipAuth)
	if err != nil {
		slog.Error("failed to load system root CA cert pool")
	}

	conn, err := grpc.Dial(backendAddr, grpc.WithTransportCredentials(creds))

	if err != nil {
		slog.Error("Cannot dial to grpc service", slog.String("error", err.Error()))
		return nil, fmt.Errorf("grpc.Dial: %w", err)
	}

	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), backendTimeout)
	defer cancel()

	ctx, err = util.GetAuthContext(ctx, audience, skipAuth)
	if err != nil {
		return nil, fmt.Errorf("error get auth context: %w", err)
	}

	// Send the request.
	client := pb.NewUploadPhotoClient(conn)
	return client.CreateMetadata(ctx, p, grpc.WaitForReady(true))
}

func pingTriggerImageUploaded(backendTimeout time.Duration,
	backendAddr string,
	p *pb.ImageUploadedRequest,
	audience string,
	skipAuth bool) (*pb.ImageUploadedResponse, error) {

	creds, err := util.CreateTransportCredentials(skipAuth)
	if err != nil {
		slog.Error("failed to load system root CA cert pool")
		return nil, fmt.Errorf("error creating CA cert pool: %w", err)
	}

	conn, err := grpc.Dial(backendAddr, grpc.WithTransportCredentials(creds))

	if err != nil {
		slog.Error("cannot dial to grpc service", slog.String("error", err.Error()))
		return nil, fmt.Errorf("grpc.Dial: %w", err)
	}

	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), backendTimeout)
	defer cancel()

	ctx, err = util.GetAuthContext(ctx, audience, skipAuth)
	if err != nil {
		slog.Error("cannot get auth context", slog.String("error", err.Error()))
		return nil, fmt.Errorf("error get auth context: %w", err)
	}

	// Send the request.
	client := pb.NewUploadPhotoClient(conn)
	imageUploaded, err := client.ImageUploaded(ctx, p, grpc.WaitForReady(true))

	if err != nil {
		slog.Error("cannot trigger the image uploaded event", slog.String("error", err.Error()))
		return nil, fmt.Errorf("error triggering the image uploaded: %w", err)
	}

	return imageUploaded, nil
}

func Upload(mux chi.Router, config types.Config) {
	mux.Post("/upload", func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, _4MB)
		_, _, err := r.FormFile("file")
		if err != nil {
			slog.Error("file is too large the limit is 4MB", slog.String("error", err.Error()))
			http.Error(w, "file is too large the limit is 4MB", http.StatusBadRequest)
			return
		}
		err = r.ParseMultipartForm(_4MB)
		if err != nil {
			http.Error(w, "error parsing form data", http.StatusBadRequest)
			return
		}

		// Get the username
		username := r.FormValue("username")
		if username == "" {
			slog.Error("username must be provided", slog.Any("error", err))
			http.Error(w, "username must be provided", http.StatusBadRequest)
			return
		}

		// Get the hashtags
		hashtags := r.FormValue("hashtags")

		// Access the array of strings using the field name
		bolb, filename, err := util.GetPhoto(r)
		if err != nil {
			slog.Error("getting file failed", slog.Any("error", err))
			http.Error(w, "getting file failed", http.StatusBadRequest)
			return
		}

		objName := uuid.New().String() + path.Ext(filename)

		// publish image upload message
		ctx := context.Background()
		client, err := pubsub.NewClient(ctx, config.ProjectID)
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		defer client.Close()

		t := client.Topic(config.TopicID)
		result := t.Publish(ctx, &pubsub.Message{
			Data: bolb,
			Attributes: map[string]string{
				"objName":  objName,
				"hashtags": hashtags,
				"user":     username,
			},
		})

		// Block until the result is returned and a server-generated
		// ID is returned for the published message.
		if _, err := result.Get(ctx); err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "image with name %s scheduled for upload successfully \n", objName)
	})
}
