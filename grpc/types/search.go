package types

import (
	"github.com/m4salah/redroc/grpc/storage"
	"go.uber.org/zap"
)

type SearchService struct {
	Log        *zap.Logger
	MetadataDB storage.MetadataDB
}
