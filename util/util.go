package util

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"strings"

	"github.com/nfnt/resize"
	"go.uber.org/zap"
	"google.golang.org/api/idtoken"
	"google.golang.org/grpc/credentials"
	grpcMetadata "google.golang.org/grpc/metadata"
)

func CreateLogger(env string) (*zap.Logger, error) {
	switch env {
	case "production":
		return zap.NewProduction()
	case "development":
		return zap.NewDevelopment()
	default:
		return zap.NewNop(), nil
	}
}

func GetPhoto(r *http.Request) ([]byte, string, error) {
	image, header, err := r.FormFile("file")
	if err != nil {
		return nil, "", fmt.Errorf("failed to get file: %v", err)
	}
	defer image.Close()

	var b bytes.Buffer
	_, err = io.Copy(&b, image)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read file: %v", err)
	}

	return b.Bytes(), header.Filename, nil
}

// GetTags function tries to read and decode photo tags from the request.
func GetTags(r *http.Request) ([]string, error) {
	var tags []string

	data := r.FormValue("hashtags")
	if len(data) == 0 {
		return tags, nil
	}

	err := json.Unmarshal([]byte(data), &tags)
	if err != nil {
		return nil, err
	}

	return tags, nil
}

func MakeThumbnail(photo []byte, width, height uint) ([]byte, error) {
	r := bytes.NewReader(photo)
	img, format, err := image.Decode(r)
	if err != nil {
		return nil, fmt.Errorf("photo decode failed: %v", err)
	}

	thumb := resize.Thumbnail(width, height, img, resize.NearestNeighbor)

	var b bytes.Buffer
	if format == "jpeg" {
		err = jpeg.Encode(&b, thumb, nil)
	} else if format == "png" {
		err = png.Encode(&b, thumb)
	} else if format == "gif" {
		err = gif.Encode(&b, thumb, nil)
	} else {
		err = fmt.Errorf("unsuported image format: %s", format)
	}
	if err != nil {
		return nil, fmt.Errorf("thumbnail encode failed: %v", err)
	}

	return b.Bytes(), nil
}

func GetAuthContext(ctx context.Context, audience string) (context.Context, error) {
	// Create an identity token.
	// With a global TokenSource tokens would be reused and auto-refreshed at need.
	// A given TokenSource is specific to the audience.
	tokenSource, err := idtoken.NewTokenSource(ctx, audience)
	if err != nil {
		return nil, fmt.Errorf("idtoken.NewTokenSource: %w", err)
	}
	token, err := tokenSource.Token()
	if err != nil {
		return nil, fmt.Errorf("TokenSource.Token: %w", err)
	}

	// Add token to gRPC Request.
	ctx = grpcMetadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+token.AccessToken)
	return ctx, nil
}

func ExtractServiceURL(addr string) string {
	return "https://" + strings.Split(addr, ":")[0]
}

// CreateTransportCredentials creates a new TLS credentials instance with the system root CA pool.
//
// This is used to create a secure connection to the server.
func CreateTransportCredentials() (credentials.TransportCredentials, error) {
	systemRoots, err := x509.SystemCertPool()
	if err != nil {
		return nil, fmt.Errorf("failed to load system root CA cert pool")
	}
	creds := credentials.NewTLS(&tls.Config{
		RootCAs: systemRoots,
	})
	return creds, nil
}
