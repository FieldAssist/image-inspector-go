package service

import (
	"context"
	"errors"
	"image"
	"testing"

	"github.com/anime-shed/image-inspector-go/internal/analyzer"
	"github.com/anime-shed/image-inspector-go/internal/repository"
	"github.com/anime-shed/image-inspector-go/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockImageRepository is a mock implementation of the ImageRepository interface.
// It allows for controlled testing of the service layer without actual network calls.


// MockAnalyzer is a mock implementation of the Analyzer interface.
// It is used to simulate the behavior of the analysis engine.


func TestImageAnalysisService_Analyze(t *testing.T) {
	dummyImage := image.NewRGBA(image.Rect(0, 0, 100, 50))
	dummyURL := "http://example.com/image.jpg"

	tests := []struct {
		name          string
		withQuality   bool
		withOCR       bool
		setupMocks    func(imageRepo *repository.MockImageRepository, analyzer *analyzer.MockAnalyzer)
		expectError   bool
		expectedCalls int
	}{
		{
			name:        "Successful analysis with quality and OCR",
			withQuality: true,
			withOCR:     true,
			setupMocks: func(imageRepo *repository.MockImageRepository, analyzer *analyzer.MockAnalyzer) {
				imageRepo.On("Fetch", mock.Anything, dummyURL).Return(dummyImage, nil)
				analyzer.On("Analyze", mock.Anything, dummyImage, mock.AnythingOfType("analyzer.ReportConfigurator"), mock.AnythingOfType("analyzer.ReportConfigurator")).Return(&models.ImageAnalysis{}, nil)
			},
			expectError:   false,
			expectedCalls: 1,
		},
		{
			name:        "Image fetch fails",
			withQuality: true,
			withOCR:     true,
			setupMocks: func(imageRepo *repository.MockImageRepository, analyzer *analyzer.MockAnalyzer) {
				imageRepo.On("Fetch", mock.Anything, dummyURL).Return(nil, errors.New("fetch error"))
			},
			expectError:   true,
			expectedCalls: 0,
		},
		{
			name:        "Analysis fails",
			withQuality: true,
			withOCR:     true,
			setupMocks: func(imageRepo *repository.MockImageRepository, analyzer *analyzer.MockAnalyzer) {
				imageRepo.On("Fetch", mock.Anything, dummyURL).Return(dummyImage, nil)
				analyzer.On("Analyze", mock.Anything, dummyImage, mock.Anything, mock.Anything).Return(nil, errors.New("analysis error"))
			},
			expectError:   true,
			expectedCalls: 1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(repository.MockImageRepository)
			mockAnalyzer := new(analyzer.MockAnalyzer)
			tc.setupMocks(mockRepo, mockAnalyzer)

			service := NewImageAnalysisService(mockRepo, mockAnalyzer)
			_, err := service.Analyze(context.Background(), dummyURL, tc.withQuality, tc.withOCR)

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
			mockAnalyzer.AssertExpectations(t)
		})
	}
}
