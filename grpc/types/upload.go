package types

import (
	"github.com/m4salah/redroc/grpc/storage"
	"go.uber.org/zap"
)

type UploadService struct {
	Log        *zap.Logger
	DB         storage.ObjectDB
	MetadataDB storage.MetadataDB
}
