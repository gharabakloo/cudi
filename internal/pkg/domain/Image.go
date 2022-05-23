package domain

import (
	"context"

	"github.com/docker/docker/api/types"
)

type ImageRepository interface {
	ReadImages(ctx context.Context, reference string) ([]types.ImageSummary, error)
	DeleteImage(ctx context.Context, imageID string) ([]types.ImageDeleteResponseItem, error)
}
