package repository

import (
	"context"
	"fmt"
	"image"
	"time"

	"github.com/anime-shed/image-inspector-go/internal/storage"
	"github.com/anime-shed/image-inspector-go/pkg/validation"
)

// HTTPImageRepository implements ImageRepository using HTTP storage
type HTTPImageRepository struct {
	fetcher   storage.ImageFetcher
	validator *validation.URLValidator
	timeout   time.Duration
}

// NewHTTPImageRepository creates a new HTTP-based image repository
func NewHTTPImageRepository(fetcher storage.ImageFetcher, timeout time.Duration) *HTTPImageRepository {
	return &HTTPImageRepository{
		fetcher:   fetcher,
		validator: validation.NewURLValidator(),
		timeout:   timeout,
	}
}

// Fetch retrieves an image from a URL
func (r *HTTPImageRepository) Fetch(ctx context.Context, imageURL string) (image.Image, error) {
	if err := r.validator.ValidateImageURL(imageURL); err != nil {
		return nil, fmt.Errorf("invalid image URL: %w", err)
	}
	return r.fetcher.FetchImage(ctx, imageURL)
}
