package storage

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	"go.uber.org/zap"
)

type FilestoreMetadata struct {
	projectID string
	log       *zap.Logger
}

// NewMetadateStorageOptions for MetadateStorage.
type NewFilestoreOptions struct {
	ProjectID string
	Log       *zap.Logger
}

// NewFilestore with the given options.
// If no logger is provided, logs are discarded.
func NewFilestore(opts NewFilestoreOptions) (*FilestoreMetadata, error) {
	if opts.ProjectID == "" {
		return nil, fmt.Errorf("ProjectID must be provided")
	}
	if opts.Log == nil {
		opts.Log = zap.NewNop()
	}
	return &FilestoreMetadata{
		log:       opts.Log,
		projectID: opts.ProjectID,
	}, nil
}

func (f *FilestoreMetadata) Store(ctx context.Context, user, path string) error {
	f.log.Info("Storing metadata")
	client, err := firestore.NewClient(ctx, f.projectID)
	if err != nil {
		f.log.Error("firestore client failed", zap.Error(err))
		return err
	}
	defer client.Close()

	// Store photo name under user.
	timestamp := time.Now().Unix()
	_, err = client.Doc(path).Set(ctx, map[string]interface{}{
		"uploaded_time": timestamp,
		"user":          user,
	})
	if err != nil {
		return fmt.Errorf("firestore create failed for %s: %v", path, err)
	}
	return nil
}
