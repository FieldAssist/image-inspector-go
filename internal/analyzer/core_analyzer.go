package analyzer

import (
	"context"
	"image"

	"github.com/anime-shed/image-inspector-go/pkg/models"
	"golang.org/x/sync/errgroup"
)

const (
	qualityAnalyzerName = "quality"
	ocrAnalyzerName     = "ocr"
)

// CoreAnalyzer is the central orchestrator for image analysis.
// It implements the Analyzer interface and coordinates various FeatureAnalyzers.
type CoreAnalyzer struct {
	featureAnalyzers map[string]FeatureAnalyzer
	ocrAnalyzer      *OCRAnalyzer // Keep a reference for closing
}

// NewCoreAnalyzer creates a new CoreAnalyzer and initializes all its feature analyzers.
func NewCoreAnalyzer() (*CoreAnalyzer, error) {
	ocrAnalyzer := NewOCRAnalyzer()

	return &CoreAnalyzer{
		featureAnalyzers: map[string]FeatureAnalyzer{
			qualityAnalyzerName: NewQualityAnalyzer(),
			ocrAnalyzerName:     ocrAnalyzer, // Reuse the same OCR analyzer
		},
		ocrAnalyzer: ocrAnalyzer,
	}, nil
}

// Analyze orchestrates the analysis by running the requested FeatureAnalyzers in parallel.
func (ca *CoreAnalyzer) Analyze(ctx context.Context, image image.Image, cfgs ...ReportConfigurator) (*models.ImageAnalysis, error) {
	config := &ReportConfig{}
	for _, cfg := range cfgs {
		cfg(config)
	}

	report := &models.ImageAnalysis{}

	// Use an error group to run analyzers concurrently for efficiency.
	eg, gCtx := errgroup.WithContext(ctx)

	if config.WithQuality {
		qualityAnalyzer := ca.featureAnalyzers[qualityAnalyzerName]
		eg.Go(func() error {
			return qualityAnalyzer.Analyze(gCtx, image, report)
		})
	}

	if config.WithOCR {
		ocrAnalyzer := ca.featureAnalyzers[ocrAnalyzerName]
		eg.Go(func() error {
			return ocrAnalyzer.Analyze(gCtx, image, report)
		})
	}

	// Wait for all selected analyzers to complete.
	if err := eg.Wait(); err != nil {
		return nil, err
	}

	return report, nil
}

// Close gracefully shuts down the OCR analyzer.
func (ca *CoreAnalyzer) Close() error {
	// Only the OCR analyzer currently needs a close method.
	if ca.ocrAnalyzer != nil {
		return ca.ocrAnalyzer.Close()
	}
	return nil
}
