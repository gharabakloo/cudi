package repository

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"

	"cudi/internal/pkg/domain"
)

type repository struct {
	client *client.Client
}

var _ domain.ImageRepository = &repository{}

func New(client *client.Client) domain.ImageRepository {
	return &repository{
		client: client,
	}
}

func (r *repository) ReadImages(ctx context.Context, reference string) ([]types.ImageSummary, error) {
	images, err := r.client.ImageList(ctx, types.ImageListOptions{
		All:     false,
		Filters: filters.NewArgs(filters.Arg(domain.Reference, reference)),
	})
	return images, err
}

func (r *repository) DeleteImage(ctx context.Context, imageID string) ([]types.ImageDeleteResponseItem, error) {
	resp, err := r.client.ImageRemove(ctx, imageID, types.ImageRemoveOptions{
		Force:         false,
		PruneChildren: false,
	})
	return resp, err
}
