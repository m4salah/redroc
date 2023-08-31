package storage

import (
	"context"
)

type MetadataDB interface {
	StorePath(ctx context.Context, path string, timestamp int64) error
	StorePathWithUser(ctx context.Context, user, path string, timestamp int64) error
	StoreLatest(ctx context.Context, index uint32, latest, objName string) error
	GetThumbnails(ctx context.Context, thumbnailCount int, keyword string) ([]string, error)
}

type ObjectDB interface {
	Store(ctx context.Context, objName string, data []byte) error
	Get(ctx context.Context, objName string) ([]byte, error)
}
