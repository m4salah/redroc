package types

import (
	"github.com/m4salah/redroc/grpc/storage"
	"go.uber.org/zap"
)

type DownloadService struct {
	Log *zap.Logger
	DB  storage.ObjectDB
}
