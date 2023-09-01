package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	pb "github.com/m4salah/redroc/libs/proto"
	"github.com/m4salah/redroc/libs/util"
	"google.golang.org/grpc"
)

// pingRequestWithAuth mints a new Identity Token for each request.
// This token has a 1 hour expiry and should be reused.
// audience must be the auto-assigned URL of a Cloud Run service or HTTP Cloud Function without port number.
func pingSearchRequestWithAuth(backendTimeout time.Duration,
	backendAddr string,
	p *pb.GetThumbnailImagesRequest,
	audience string,
	skipAuth bool) (*pb.GetThumbnailImagesResponse, error) {

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
	client := pb.NewGetThumbnailClient(conn)
	return client.GetThumbnail(ctx, p, grpc.WaitForReady(true))
}

func Search(mux chi.Router, backendAddr string, backendTimeout time.Duration, skipAuth bool) {
	mux.Get("/search", func(w http.ResponseWriter, r *http.Request) {
		// get the query string
		q := r.URL.Query().Get("q")

		slog.Info("Searching for", slog.String("searchQuery", q))

		request := &pb.GetThumbnailImagesRequest{SearchKeyword: q}
		response, err := pingSearchRequestWithAuth(backendTimeout, backendAddr, request, util.ExtractServiceURL(backendAddr), skipAuth)
		if err != nil {
			slog.Error("search request failed", slog.String("error", err.Error()))
			http.Error(w, "search request failed", http.StatusBadRequest)
			return
		}

		urls, err := json.Marshal(response.StorageUrl)
		if err != nil {
			slog.Error("json marshal failed", slog.String("error", err.Error()))
			http.Error(w, "Something Wrong", http.StatusInternalServerError)
			return
		}
		w.Write(urls)
	})
}
