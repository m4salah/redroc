package handlers

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"time"

	"cloud.google.com/go/storage"
	"github.com/go-chi/chi/v5"
	pb "github.com/m4salah/redroc/libs/proto"
	"github.com/m4salah/redroc/libs/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// pingRequestWithAuth mints a new Identity Token for each request.
// This token has a 1 hour expiry and should be reused.
// audience must be the auto-assigned URL of a Cloud Run service or HTTP Cloud Function without port number.
func pingDownloadRequestWithAuth(backendTimeout time.Duration,
	backendAddr net.Addr,
	p *pb.DownloadPhotoRequest,
	audience string,
	skipAuth bool) (*pb.DownloadPhotoResponse, error) {
	creds, err := util.CreateTransportCredentials(skipAuth)
	if err != nil {
		slog.Error("failed to load system root CA cert pool")
	}
	conn, err := grpc.Dial(backendAddr.String(), grpc.WithTransportCredentials(creds))

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
	client := pb.NewDownloadPhotoClient(conn)
	return client.Download(ctx, p, grpc.WaitForReady(true))
}

func Download(mux chi.Router, backendAddr net.Addr, backendTimeout time.Duration, skipAuth bool) {
	mux.Get("/download/{imgName}", func(w http.ResponseWriter, r *http.Request) {
		imgName := chi.URLParam(r, "imgName")

		slog.Info("Downloading image", slog.String("imageName", imgName))
		request := &pb.DownloadPhotoRequest{ImgName: imgName}

		response, err := pingDownloadRequestWithAuth(backendTimeout, backendAddr, request, util.ExtractServiceURL(backendAddr), skipAuth)
		if err != nil {
			slog.Error("downloading image failed", slog.String("error", err.Error()))
			grpcStatus, ok := status.FromError(err)
			if ok && grpcStatus.Message() == storage.ErrObjectNotExist.Error() {
				http.Error(w, "Image not found", http.StatusNotFound)
				return
			}
			http.Error(w, "Downloading Image failed", http.StatusBadRequest)
			return
		}
		w.Write(response.ImgBlob)
	})
}
