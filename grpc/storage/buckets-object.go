package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/storage"
	"github.com/m4salah/redroc/util"
	"go.uber.org/zap"
)

type BucketsObject struct {
	bucketName string
	log        *zap.Logger
}

// NewBucketsOptions for MetadateStorage.
type NewBucketsOptions struct {
	BucketName string
	Log        *zap.Logger
}

// NewBuckets with the given options.
// If no logger is provided, logs are discarded.
func NewBuckets(opts NewBucketsOptions) (*BucketsObject, error) {
	if opts.BucketName == "" {
		return nil, fmt.Errorf("BucketName must be provided")
	}
	if opts.Log == nil {
		opts.Log = zap.NewNop()
	}
	return &BucketsObject{
		log:        opts.Log,
		bucketName: opts.BucketName,
	}, nil
}

func (b *BucketsObject) Store(ctx context.Context, objName string, data []byte) error {
	b.log.Info("Storing Image to cloud storage", zap.String("objName", objName))
	client, err := storage.NewClient(ctx)
	if err != nil {
		b.log.Error("storage client failed", zap.Error(err))
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
	bo.log.Info("Getting Object from cloud storage", zap.String("objName", objName))
	client, err := storage.NewClient(ctx)
	if err != nil {
		bo.log.Error("storage client failed", zap.Error(err))
		return nil, err
	}
	defer client.Close()

	blob := client.Bucket(bo.bucketName).Object(objName)
	r, err := blob.NewReader(ctx)
	if err != nil {
		if err == storage.ErrBucketNotExist {
			bo.log.Error("Bucket doesn't exists", zap.String("bucketName", bo.bucketName), zap.Error(err))
			return nil, err
		}
		if err == storage.ErrObjectNotExist {
			bo.log.Error("Object doesn't exists", zap.String("objName", objName), zap.Error(err))
			return nil, err
		}
		bo.log.Error("storage reader failed", zap.String("objName", objName), zap.String("bucketName", bo.bucketName), zap.Error(err))
		return nil, err
	}
	defer r.Close()

	var b bytes.Buffer
	_, err = b.ReadFrom(r)
	if err != nil {
		bo.log.Error("storage reading failed", zap.String("objName", objName), zap.Error(err))
		return nil, err
	}
	// DONE: make secret from env variable
	secret := util.GetStringOrDefault("ENCRYPTION_KEY", "")
	if secret == "" {
		return nil, fmt.Errorf("ENCRYPTION_KEY must be provided")
	}
	return util.DecryptAES(b.Bytes(), []byte(secret))
}
