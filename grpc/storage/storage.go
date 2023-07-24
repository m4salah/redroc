package storage

import "context"

type MetadataDB interface {
	Store(ctx context.Context, user, path string) error
}

type ObjectDB interface {
	Store(ctx context.Context, objName string, data []byte) error
	Get(ctx context.Context, objName string) ([]byte, error)
}
