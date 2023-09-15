package handlers

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"path"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/m4salah/redroc/libs/util"
	"golang.org/x/sync/errgroup"
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

func Upload(mux chi.Router, backendAddr string, backendTimeout time.Duration, skipAuth bool) {
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
			slog.Error("username must be provided", slog.String("error", err.Error()))
			http.Error(w, "username must be provided", http.StatusBadRequest)
			return
		}

		// Get the hashtags
		hashtags, err := util.GetTags(r)
		if err != nil {
			slog.Error("getting tags failes", slog.String("error", err.Error()))
			http.Error(w, "invalid tags", http.StatusBadRequest)
			return

		}

		// Access the array of strings using the field name
		blob, filename, err := util.GetPhoto(r)
		if err != nil {
			slog.Error("getting file failed", slog.String("error", err.Error()))
			http.Error(w, "getting file failed", http.StatusBadRequest)
			return
		}

		objName := uuid.New().String() + path.Ext(filename)

		eg := new(errgroup.Group)

		// uploading image
		eg.Go(func() error {
			slog.Info("uploading image", slog.String("imageName", objName))
			uploadRequest := &pb.UploadImageRequest{ObjName: objName, Image: blob}
			_, err = pingUploadRequestWithAuth(backendTimeout, backendAddr, uploadRequest, util.ExtractServiceURL(backendAddr), skipAuth)
			return err
		})

		// writing the metadata
		eg.Go(func() error {
			slog.Info("writing metadata for image",
				slog.String("imageName", objName),
				slog.String("username", username),
				slog.Any("hashtags", hashtags))

			metadataRequest := &pb.CreateMetadataRequest{
				ObjName: objName, User: username, Hashtags: hashtags}

			_, err = pingCreateMetadataRequestWithAuth(backendTimeout, backendAddr, metadataRequest, util.ExtractServiceURL(backendAddr), skipAuth)
			return err
		})

		if err := eg.Wait(); err != nil {
			slog.Error("error while uploading or creating metadata", "error", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		slog.Info("image Uploaded Successfully",
			slog.String("imageName", objName),
			slog.String("username", username),
			slog.Any("hashtags", hashtags))

		// Trigger the ImageUploaded event
		go pingTriggerImageUploaded(backendTimeout, backendAddr, &pb.ImageUploadedRequest{
			ObjName:  objName,
			User:     username,
			Hashtags: hashtags,
		}, util.ExtractServiceURL(backendAddr), skipAuth)

		fmt.Fprintln(w, "image Uploaded Successfully")
	})
}
