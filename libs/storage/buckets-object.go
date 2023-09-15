package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"

	"cloud.google.com/go/storage"
	"github.com/m4salah/redroc/libs/util"
)

type BucketsObject struct {
	bucketName string
}

// NewBucketsOptions for MetadateStorage.
type NewBucketsOptions struct {
	BucketName string
}

// NewBuckets with the given options.
// If no logger is provided, logs are discarded.
func NewBuckets(opts NewBucketsOptions) (*BucketsObject, error) {
	if opts.BucketName == "" {
		return nil, fmt.Errorf("BucketName must be provided")
	}
	return &BucketsObject{
		bucketName: opts.BucketName,
	}, nil
}

func (b *BucketsObject) Store(ctx context.Context, objName string, data []byte) error {
	slog.Info("Storing Image to cloud storage", slog.String("objName", objName))
	client, err := storage.NewClient(ctx)
	if err != nil {
		slog.Error("storage client failed", err)
		return err
	}
	defer client.Close()

	bucket := client.Bucket(b.bucketName)

	obj := bucket.Object(objName)
	w := obj.NewWriter(ctx)
	// DONE: make secret from env variable
	secret := util.GetStringOrDefault("ENCRYPTION_KEY", "")
	if secret == "" {
		return fmt.Errorf("ENCRYPTION_KEY must be provided")
	}
	encryptedData, err := util.EncryptAES(data, []byte(secret))
	if err != nil {
		return fmt.Errorf("storage encryption failed: %v", err)
	}
	r := bytes.NewReader(encryptedData)

	_, err = io.Copy(w, r)

	if err != nil {
		return fmt.Errorf("storage copy failed: %v", err)
	}

	err = w.Close()
	if err != nil {
		return fmt.Errorf("storage close failed: %v", err)
	}

	return nil
}

func (bo *BucketsObject) Get(ctx context.Context, objName string) ([]byte, error) {
	slog.Info("Getting Object from cloud storage", slog.String("objName", objName))
	client, err := storage.NewClient(ctx)
	if err != nil {
		slog.Error("storage client failed", err)
		return nil, err
	}
	defer client.Close()

	blob := client.Bucket(bo.bucketName).Object(objName)
	r, err := blob.NewReader(ctx)
	if err != nil {
		if err == storage.ErrBucketNotExist {
			slog.Error("Bucket doesn't exists", slog.String("bucketName", bo.bucketName), err)
			return nil, err
		}
		if err == storage.ErrObjectNotExist {
			slog.Error("Object doesn't exists", slog.String("objName", objName), err)
			return nil, err
		}
		slog.Error("storage reader failed", slog.String("objName", objName), slog.String("bucketName", bo.bucketName), err)
		return nil, err
	}
	defer r.Close()

	var b bytes.Buffer
	_, err = b.ReadFrom(r)
	if err != nil {
		slog.Error("storage reading failed", slog.String("objName", objName), err)
		return nil, err
	}
	secret := util.GetStringOrDefault("ENCRYPTION_KEY", "")
	if secret == "" {
		return nil, fmt.Errorf("ENCRYPTION_KEY must be provided")
	}
	return b.Bytes(), nil
}
