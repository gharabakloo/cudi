package domain

import (
	"context"
)

type Cleanup struct {
	Images []Image `json:"images"`
}

type Image struct {
	Type       string   `json:"type"`
	Repository string   `json:"repository"`
	KeepNumber int      `json:"keepNumber"`
	KeepTags   []string `json:"keepTags"`
	RemoveTags []string `json:"removeTags"`
	OlderThan  string   `json:"olderThan"`
}

type CleanupUsecase interface {
	CleanupImages(ctx context.Context, cleanup Cleanup) error
	CleanupImage(ctx context.Context, image Image) error
}
