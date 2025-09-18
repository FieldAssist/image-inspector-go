package analyzer

import (
	"context"
	"errors"
	"image"
	"testing"

	"github.com/anime-shed/image-inspector-go/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAnalyzer is a mock implementation of the Analyzer interface for testing.
type MockAnalyzer struct {
	mock.Mock
}

func (m *MockAnalyzer) Analyze(ctx context.Context, img image.Image, cfgs ...ReportConfigurator) (*models.ImageAnalysis, error) {
	args := m.Called(ctx, img, cfgs)
	// We can't directly compare funcs, so we handle configurators separately.

	// Create a config and apply the passed-in configurators to it.
	config := &ReportConfig{}
	for _, cfg := range cfgs {
		cfg(config)
	}

	// Here, you could add assertions based on the final config state if needed,
	// for example, checking if WithOCR() was passed.

	var report *models.ImageAnalysis
	if args.Get(0) != nil {
		report = args.Get(0).(*models.ImageAnalysis)
	}

	return report, args.Error(1)
}

func (m *MockAnalyzer) Close() error {
	args := m.Called()
	return args.Error(0)
}

// MockFeatureAnalyzer is a mock implementation of the FeatureAnalyzer interface.
type MockFeatureAnalyzer struct {
	mock.Mock
}

func (m *MockFeatureAnalyzer) Analyze(ctx context.Context, img image.Image, report *models.ImageAnalysis) error {
	args := m.Called(ctx, img, report)
	return args.Error(0)
}

func TestMockAnalyzer(t *testing.T) {
	t.Run("Analyze method mock", func(t *testing.T) {
		mockAnalyzer := new(MockAnalyzer)
		expectedReport := &models.ImageAnalysis{}
		img := image.NewRGBA(image.Rect(0, 0, 10, 10))

		// We expect the Analyze method to be called. We can't match the configurators directly,
		// so we use mock.Anything.
		mockAnalyzer.On("Analyze", mock.Anything, mock.Anything, mock.Anything).Return(expectedReport, nil)

		report, err := mockAnalyzer.Analyze(context.Background(), img)

		assert.NoError(t, err)
		assert.Equal(t, expectedReport, report)
		mockAnalyzer.AssertExpectations(t)
	})

	t.Run("Close method mock", func(t *testing.T) {
		mockAnalyzer := new(MockAnalyzer)
		mockAnalyzer.On("Close").Return(nil)

		err := mockAnalyzer.Close()

		assert.NoError(t, err)
		mockAnalyzer.AssertExpectations(t)
	})
}

func TestMockFeatureAnalyzer(t *testing.T) {
	t.Run("Analyze method mock with error", func(t *testing.T) {
		mockFeature := new(MockFeatureAnalyzer)
		img := image.NewRGBA(image.Rect(0, 0, 10, 10))
		report := &models.ImageAnalysis{}
		expectedError := errors.New("feature analysis failed")

		mockFeature.On("Analyze", mock.Anything, img, report).Return(expectedError)

		err := mockFeature.Analyze(context.Background(), img, report)

		assert.Equal(t, expectedError, err)
		mockFeature.AssertExpectations(t)
	})
}
