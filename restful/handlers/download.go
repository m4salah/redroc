package handlers

import (
	"context"
	"crypto/tls"
	"net/http"
	"time"

	"cloud.google.com/go/storage"
	"github.com/go-chi/chi/v5"
	pb "github.com/m4salah/redroc/grpc/protos"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

func Download(mux chi.Router, backendAddr string, log *zap.Logger, backendTimeout time.Duration) {
	mux.Get("/download/{imgName}", func(w http.ResponseWriter, r *http.Request) {
		imgName := chi.URLParam(r, "imgName")
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

		ctx, cancel := context.WithTimeout(context.Background(), backendTimeout)
		defer cancel()

		request := &pb.DownloadPhotoRequest{ImgName: imgName}

		client := pb.NewDownloadPhotoClient(conn)
		response, err := client.Download(ctx, request, grpc.WaitForReady(true))
		if err != nil {
			log.Error("downloading image failed", zap.Error(err))
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
