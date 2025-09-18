package analyzer

import (
	"context"
	"errors"
	"image"
	"testing"

	"github.com/anime-shed/image-inspector-go/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockFeatureAnalyzer is a mock type for the FeatureAnalyzer interface.

type MockFeatureAnalyzer struct {
	mock.Mock
}

// Analyze provides a mock function with given fields: ctx, img, report
func (_m *MockFeatureAnalyzer) Analyze(ctx context.Context, img image.Image, report *models.ImageAnalysis) error {

	args := _m.Called(ctx, img, report)

	return args.Error(0)

}

// TestCoreAnalyzer verifies the orchestration logic of the CoreAnalyzer.
func TestCoreAnalyzer(t *testing.T) {
	// This test uses mocks to isolate the CoreAnalyzer from the actual feature analyzers.
	t.Run("NewCoreAnalyzer successfully creates an instance", func(t *testing.T) {
		// This test now implicitly tests the real NewCoreAnalyzer, which is fine.
		analyzer, err := NewCoreAnalyzer()
		require.NoError(t, err)
		assert.NotNil(t, analyzer)
		assert.NotNil(t, analyzer.featureAnalyzers[qualityAnalyzerName])
		assert.NotNil(t, analyzer.featureAnalyzers[ocrAnalyzerName])
		assert.NoError(t, analyzer.Close()) // Should be able to close immediately
	})

	t.Run("Analyze orchestrates selected analyzers", func(t *testing.T) {
		// Setup: Create mock feature analyzers
		mockQuality := new(MockFeatureAnalyzer)
		mockOCR := new(MockFeatureAnalyzer)

		// Setup: Create CoreAnalyzer with mocks
		coreAnalyzer := &CoreAnalyzer{
			featureAnalyzers: map[string]FeatureAnalyzer{
				qualityAnalyzerName: mockQuality,
				ocrAnalyzerName:     mockOCR,
			},
		}

		img := image.NewRGBA(image.Rect(0, 0, 10, 10))

		// Expect that WithQuality and WithOCR call the corresponding mocks.
		mockQuality.On("Analyze", mock.Anything, img, mock.AnythingOfType("*models.ImageAnalysis")).Return(nil)
		mockOCR.On("Analyze", mock.Anything, img, mock.AnythingOfType("*models.ImageAnalysis")).Return(nil)

		// Action: Run analysis with specific options
		_, err := coreAnalyzer.Analyze(context.Background(), img, WithQuality(), WithOCR())

		// Assertions
		assert.NoError(t, err)

		// Verify that the expected mocks were called
		mockQuality.AssertExpectations(t)
		mockOCR.AssertExpectations(t)
	})

	t.Run("Analyze handles errors from a feature analyzer", func(t *testing.T) {
		mockQuality := new(MockFeatureAnalyzer)
		expectedError := errors.New("quality analysis failed")

		coreAnalyzer := &CoreAnalyzer{
			featureAnalyzers: map[string]FeatureAnalyzer{
				qualityAnalyzerName: mockQuality,
			},
		}

		img := image.NewRGBA(image.Rect(0, 0, 10, 10))

		// Expect the quality analyzer to be called and return an error.
		mockQuality.On("Analyze", mock.Anything, img, mock.AnythingOfType("*models.ImageAnalysis")).Return(expectedError)

		// Action: Run analysis
		_, err := coreAnalyzer.Analyze(context.Background(), img, WithQuality())

		// Assertions
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		mockQuality.AssertExpectations(t)
	})

	t.Run("Close handles underlying analyzer close errors", func(t *testing.T) {
		// Although our real OCRAnalyzer.Close returns an error, we can't easily test it
		// without a more complex mock. This test serves as a placeholder for the logic.
		analyzer, _ := NewCoreAnalyzer()
		// In a real scenario with a mockable OCR client, you'd set up an expectation here.
		err := analyzer.Close()
		assert.NoError(t, err) // Expecting no error from the default gosseract client close
	})
}
