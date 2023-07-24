package download

import (
	"context"

	"github.com/m4salah/redroc/grpc/storage"
	"go.uber.org/zap"
)

type DownloadService struct {
	Log *zap.Logger
	DB  storage.ObjectDB
}

func (d *DownloadService) Download(ctx context.Context, imageName string) ([]byte, error) {
	image, err := d.DB.Get(ctx, imageName)
	if err != nil {
		return nil, err
	}
	return image, nil
}
