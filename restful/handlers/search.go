package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	pb "github.com/m4salah/redroc/grpc/protos"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Search(mux chi.Router, backendAddr string, log *zap.Logger, backendTimeout time.Duration) {
	mux.Get("/search", func(w http.ResponseWriter, r *http.Request) {
		conn, err := grpc.Dial(backendAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Error("Cannot dial to grpc service", zap.Error(err))
			http.Error(w, "Cannot dial download service", http.StatusBadRequest)
			return
		}
		defer conn.Close()

		ctx, cancel := context.WithTimeout(context.Background(), backendTimeout)
		defer cancel()

		// get the query string
		q := r.URL.Query().Get("q")
		log.Info("searching for", zap.String("q", q))
		request := &pb.GetThumbnailImagesRequest{SearchKeyword: q}
		client := pb.NewGetThumbnailClient(conn)
		response, err := client.GetThumbnail(ctx, request, grpc.WaitForReady(true))
		if err != nil {
			log.Error("search request failed", zap.Error(err))
			http.Error(w, "search request failed", http.StatusBadRequest)
			return
		}

		urls, err := json.Marshal(response.StorageUrl)
		if err != nil {
			log.Error("json marshal failed", zap.Error(err))
			http.Error(w, "Something Wrong", http.StatusInternalServerError)
			return
		}
		w.Write(urls)
	})
}
