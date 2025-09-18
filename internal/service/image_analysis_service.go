package service

import (
	"context"

	"github.com/anime-shed/image-inspector-go/internal/analyzer"
	"github.com/anime-shed/image-inspector-go/internal/repository"
	"github.com/anime-shed/image-inspector-go/pkg/models"
)

// ImageAnalysisService provides a high-level interface for image analysis operations.
// It orchestrates the fetching of images and the execution of various analyses.
type ImageAnalysisService struct {
	imageRepository repository.ImageRepository
	imageAnalyzer   analyzer.Analyzer
}

// NewImageAnalysisService creates a new service instance with the given dependencies.
// It depends on interfaces (DIP) for better testability and flexibility.
func NewImageAnalysisService(repo repository.ImageRepository, anlz analyzer.Analyzer) *ImageAnalysisService {
	return &ImageAnalysisService{
		imageRepository: repo,
		imageAnalyzer:   anlz,
	}
}

// Analyze is the unified method for performing image analysis.
// It fetches the image from the provided URL and runs the analysis based on the boolean flags.
func (s *ImageAnalysisService) Analyze(ctx context.Context, imageURL string, withQuality, withOCR bool) (*models.ImageAnalysis, error) {
	// 1. Fetch the image using the repository
	image, err := s.imageRepository.Fetch(ctx, imageURL)
	if err != nil {
		return nil, err // Errors are handled by the caller (transport layer)
	}

	// 2. Build the analysis configuration dynamically using functional options
	var configs []analyzer.ReportConfigurator
	if withQuality {
		configs = append(configs, analyzer.WithQuality())
	}
	if withOCR {
		configs = append(configs, analyzer.WithOCR())
	}

	// 3. Delegate the analysis to the analyzer component
	return s.imageAnalyzer.Analyze(ctx, image, configs...)
}
