package repository

import (
	"context"
	"image"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"errors"

	"github.com/anime-shed/image-inspector-go/pkg/models"
)

// MockImageRepository is a mock implementation of the ImageRepository interface.
// It is used in tests to isolate the service layer from the repository layer.

type MockImageRepository struct {
	mock.Mock
}

// Fetch is a mock implementation of the Fetch method.
func (m *MockImageRepository) Fetch(ctx context.Context, url string) (image.Image, error) {
	args := m.Called(ctx, url)
	// Handle the case where the first return value is nil
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(image.Image), args.Error(1)
}

// Save is a mock implementation of the Save method.
func (m *MockImage.Repository) Save(ctx context.Context, analysis *models.ImageAnalysis) error {
	args := m.Called(ctx, analysis)
	return args.Error(0)
}

func TestMockImageRepository_Fetch_Success(t *testing.T) {
	mockRepo := new(MockImageRepository)
	dummyImage := image.NewRGBA(image.Rect(0, 0, 10, 10))
	dummyURL := "http://example.com/image.png"

	mockRepo.On("Fetch", context.Background(), dummyURL).Return(dummyImage, nil)

	img, err := mockRepo.Fetch(context.Background(), dummyURL)

	assert.NoError(t, err)
	assert.NotNil(t, img)
	mockRepo.AssertExpectations(t)
}

func TestMockImageRepository_Fetch_Error(t *testing.T) {
	mockRepo := new(MockImageRepository)
	dummyURL := "http://example.com/image.png"
	expectedError := errors.New("network error")

	mockRepo.On("Fetch", context.Background(), dummyURL).Return(nil, expectedError)

	img, err := mockRepo.Fetch(context.Background(), dummyURL)

	assert.Error(t, err)
	assert.Nil(t, img)
	assert.Equal(t, expectedError, err)
	mockRepo.AssertExpectations(t)
}
