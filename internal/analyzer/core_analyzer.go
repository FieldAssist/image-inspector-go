package analyzer

import (
	"fmt"
	"image"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/anime-shed/image-inspector-go/internal/loader"
	"github.com/anime-shed/image-inspector-go/pkg/models"
	"github.com/anime-shed/image-inspector-go/pkg/validation"
)

// coreAnalyzer implements ImageAnalyzer interface with enhanced performance
// Implements optimizations from PERFORMANCE_OPTIMIZATION_ANALYSIS.md Phase 3
type coreAnalyzer struct {
	workerPool        *WorkerPool
	metricsCalculator MetricsCalculator
	qualityValidator  *validation.QualityValidator
	qrDetector        QRDetector

	// Enhanced memory pools with better sizing
	resultPool    sync.Pool
	
	// Performance monitoring
	analysisCount    int64
	totalProcessTime time.Duration
	mu               sync.RWMutex
}

// NewCoreAnalyzer creates a new core analyzer
func NewCoreAnalyzer() (ImageAnalyzer, error) {
	workerPool := NewWorkerPool(0) // Use default CPU count
	workerPool.Start()

	return &coreAnalyzer{
		workerPool:        workerPool,
		metricsCalculator: NewMetricsCalculator(),
		qualityValidator:  validation.NewQualityValidator(),
		qrDetector:        NewQRDetector(),

		// Enhanced memory pools
		resultPool: sync.Pool{
			New: func() interface{} {
				return &AnalysisResult{}
			},
		},
	}, nil
}

// Analyze performs basic image analysis with memory management
func (oca *coreAnalyzer) Analyze(img *loader.FastImage) (*AnalysisResult, error) {
	options := DefaultOptions()
	return oca.AnalyzeWithOptions(img, options)
}

// AnalyzeWithOCR performs OCR-specific image analysis (legacy method for backward compatibility)
func (oca *coreAnalyzer) AnalyzeWithOCR(img *loader.FastImage, expectedText string) (*AnalysisResult, error) {
	options := OCROptions().WithOCR(expectedText)
	return oca.AnalyzeWithOptions(img, options)
}

// AnalyzeWithOptions performs image analysis with enhanced parallel processing and memory optimization
func (oca *coreAnalyzer) AnalyzeWithOptions(img *loader.FastImage, options AnalysisOptions) (*AnalysisResult, error) {
	start := time.Now()
	defer func() {
		oca.updatePerformanceStats(time.Since(start))
	}()

	// Get result from pool and reset it efficiently
	result := oca.resultPool.Get().(*AnalysisResult)
	*result = AnalysisResult{} // Reset the result
	
	// Use defer with anonymous function to ensure cleanup even on panic
	defer func() {
		if r := recover(); r != nil {
			oca.resultPool.Put(result)
			panic(r) // Re-panic after cleanup
		}
	}()

	result.Timestamp = start

	// Set expected text in OCR result if provided
	if options.OCRExpectedText != "" {
		if result.OCRResult == nil {
			result.OCRResult = &models.OCRResult{}
		}
		result.OCRResult.ExpectedText = options.OCRExpectedText
		// OCR is not implemented yet, set error message
		result.OCRResult.OCRError = "OCR text extraction is not implemented in this version"
	}

	// Parallel processing of different analysis components
	if options.UseWorkerPool && !options.FastMode {
		oca.analyzeWithParallelProcessing(img, result, options)
	} else {
		oca.analyzeSequentially(img, result, options)
	}

	// Calculate processing time
	processingTime := time.Since(start).Seconds()
	result.ProcessingTimeSec = processingTime

	// Create a copy to return
	finalResult := *result
	// Ensure processing time is copied
	finalResult.ProcessingTimeSec = processingTime
	
	// Return result to pool before returning the copy
	oca.resultPool.Put(result)
	return &finalResult, nil
}

// analyzeWithParallelProcessing performs analysis using parallel worker pool
func (oca *coreAnalyzer) analyzeWithParallelProcessing(img *loader.FastImage, result *AnalysisResult, options AnalysisOptions) {
	var wg sync.WaitGroup
	var mu sync.Mutex

	// Set resolution information (needed for quality validation)
	meta, err := img.Metadata()
	width, height := 0, 0
	if err == nil {
		width, height = meta.Size.Width, meta.Size.Height
	}
	result.Metrics.Resolution = fmt.Sprintf("%dx%d", width, height)

	// Basic metrics calculation
	wg.Add(1)
	oca.workerPool.Submit(func() {
		defer wg.Done()
		metrics, err := oca.metricsCalculator.CalculateBasicMetrics(img)
		if err == nil {
			mu.Lock()
			result.Metrics.AvgLuminance = metrics.avgLuminance
			result.Metrics.AvgSaturation = metrics.avgSaturation
			result.Metrics.ChannelBalance = [3]float64{metrics.avgR, metrics.avgG, metrics.avgB}
			mu.Unlock()
		}
	})

	// Laplacian variance calculation
	wg.Add(1)
	oca.workerPool.Submit(func() {
		defer wg.Done()
		laplacianVar, err := oca.metricsCalculator.CalculateLaplacianVariance(img)
		if err == nil {
			mu.Lock()
			result.Metrics.LaplacianVar = laplacianVar
			result.Quality.Blurry = laplacianVar <= options.BlurThreshold
			mu.Unlock()
		}
	})

	// QR detection (if enabled)
	if !options.SkipQRDetection {
		wg.Add(1)
		oca.workerPool.Submit(func() {
			defer wg.Done()
			qrDetected := oca.qrDetector.DetectQRCode(img)

			mu.Lock()
			result.Quality.QRDetected = qrDetected
			mu.Unlock()
		})
	}

	// Enhanced quality checks for OCR mode
	if options.OCRMode {
		wg.Add(1)
		oca.workerPool.Submit(func() {
			defer wg.Done()
			oca.performEnhancedQualityChecks(img, result, options)
		})
	}

	wg.Wait()

	// Perform quality validation to populate error messages
	oca.performQualityValidation(result, options)

	// Post-process results
	oca.finalizeAnalysisResults(result, options)
}

// analyzeSequentially performs analysis without parallel processing (for fast mode)
func (oca *coreAnalyzer) analyzeSequentially(img *loader.FastImage, result *AnalysisResult, options AnalysisOptions) {
	// Set resolution information
	meta, err := img.Metadata()
	width, height := 0, 0
	if err == nil {
		width, height = meta.Size.Width, meta.Size.Height
	}
	result.Metrics.Resolution = fmt.Sprintf("%dx%d", width, height)

	// Calculate basic metrics
	metrics, err := oca.metricsCalculator.CalculateBasicMetrics(img)
	if err == nil {
		result.Metrics.AvgLuminance = metrics.avgLuminance
		result.Metrics.AvgSaturation = metrics.avgSaturation
		result.Metrics.ChannelBalance = [3]float64{metrics.avgR, metrics.avgG, metrics.avgB}
	}

	// Calculate Laplacian variance for blur detection
	lapVar, err := oca.metricsCalculator.CalculateLaplacianVariance(img)
	if err == nil {
		result.Metrics.LaplacianVar = lapVar
		result.Quality.Blurry = result.Metrics.LaplacianVar <= options.BlurThreshold
	}

	// Check for overexposure and oversaturation
	result.Quality.Overexposed = metrics.avgLuminance > options.OverexposureThreshold
	result.Quality.Oversaturated = metrics.avgSaturation > options.OversaturationThreshold

	// Check white balance (skip if disabled)
	if !options.SkipWhiteBalance {
		result.Quality.IncorrectWB = oca.hasWhiteBalanceIssue(metrics.avgR, metrics.avgG, metrics.avgB)
	}

	// Detect QR codes (skip if disabled)
	if !options.SkipQRDetection {
		result.Quality.QRDetected = oca.qrDetector.DetectQRCode(img)
	}

	// Enhanced quality checks for OCR mode
	if options.OCRMode {
		oca.performEnhancedQualityChecks(img, result, options)
	}

	// Perform quality validation to populate error messages
	oca.performQualityValidation(result, options)

	oca.finalizeAnalysisResults(result, options)
}

// performEnhancedQualityChecks performs additional quality checks for OCR with optimizations
func (oca *coreAnalyzer) performEnhancedQualityChecks(img *loader.FastImage, result *AnalysisResult, options AnalysisOptions) {
	meta, _ := img.Metadata()
	width, height := meta.Size.Width, meta.Size.Height
	
	// Set resolution information
	result.Metrics.Resolution = fmt.Sprintf("%dx%d", width, height)
	result.Quality.IsLowResolution = width*height < 800000 || width < 800 || height < 1000

	// Calculate brightness
	brightness, err := oca.metricsCalculator.CalculateBrightness(img)
	if err == nil {
		result.Metrics.Brightness = brightness
		result.Quality.IsTooDark = result.Metrics.Brightness < 75
		result.Quality.IsTooBright = result.Metrics.Brightness > 230
	}

	// Detect skew
	skewAngle, err := oca.metricsCalculator.DetectSkew(img)
	if err == nil && skewAngle != nil {
		result.Quality.SkewAngle = skewAngle
		result.Quality.IsSkewed = *skewAngle > 5 || *skewAngle < -5
	}

	// Count contours (skip if disabled)
	if !options.SkipContourDetection {
		contours, _ := oca.metricsCalculator.DetectContours(img)
		result.Metrics.NumContours = contours
	}

	// Simple document edge detection (skip if disabled)
	if !options.SkipEdgeDetection {
		result.Quality.HasDocumentEdges = oca.detectDocumentEdges(img)
	}

	// Perform quality validation using QualityValidator
	oca.performQualityValidation(result, options)
}

// performQualityValidation uses QualityValidator to generate quality error messages
func (oca *coreAnalyzer) performQualityValidation(result *AnalysisResult, options AnalysisOptions) {
	// Prepare metrics for validation
	width := oca.getWidthFromResolution(result.Metrics.Resolution)
	height := oca.getHeightFromResolution(result.Metrics.Resolution)

	metrics := validation.ImageQualityMetrics{
		Width:            width,
		Height:           height,
		LaplacianVar:     result.Metrics.LaplacianVar,
		Brightness:       result.Metrics.Brightness,
		AvgLuminance:     result.Metrics.AvgLuminance,
		AvgSaturation:    result.Metrics.AvgSaturation,
		ChannelBalance:   result.Metrics.ChannelBalance,
		Overexposed:      result.Quality.Overexposed,
		Oversaturated:    result.Quality.Oversaturated,
		IncorrectWB:      result.Quality.IncorrectWB,
		IsTooDark:        result.Quality.IsTooDark,
		IsTooBright:      result.Quality.IsTooBright,
		IsSkewed:         result.Quality.IsSkewed,
		HasDocumentEdges: result.Quality.HasDocumentEdges,
		SkewAngle:        result.Quality.SkewAngle,
	}

	// Perform appropriate validation based on mode
	var issues []validation.QualityIssue
	if options.OCRMode {
		issues = oca.qualityValidator.ValidateOCRQuality(metrics)
	} else {
		issues = oca.qualityValidator.ValidateBasicQuality(metrics)
	}

	// Convert issues to error messages
	if len(issues) > 0 {
		result.Errors = oca.qualityValidator.ConvertIssuesToMessages(issues)
	}
}

// getWidthFromResolution extracts width from resolution string
func (oca *coreAnalyzer) getWidthFromResolution(resolution string) int {
	width, _ := oca.parseResolution(resolution)
	return width
}

// getHeightFromResolution extracts height from resolution string
func (oca *coreAnalyzer) getHeightFromResolution(resolution string) int {
	_, height := oca.parseResolution(resolution)
	return height
}

// finalizeAnalysisResults performs final processing and validation
func (oca *coreAnalyzer) finalizeAnalysisResults(result *AnalysisResult, options AnalysisOptions) {
	// Set overall validity based on quality checks and validation errors
	hasQualityIssues := result.Quality.Blurry ||
		result.Quality.Overexposed ||
		result.Quality.Oversaturated ||
		(options.OCRMode && (result.Quality.IsTooDark || result.Quality.IsTooBright))

	// Also consider validation errors from QualityValidator
	hasValidationErrors := len(result.Errors) > 0

	// Image is valid only if it has no quality issues AND no validation errors
	result.Quality.IsValid = !hasQualityIssues && !hasValidationErrors
}

// hasWhiteBalanceIssue checks for white balance issues
func (oca *coreAnalyzer) hasWhiteBalanceIssue(avgR, avgG, avgB float64) bool {
	threshold := 0.15
	maxChannel := maxFloat64(avgR, maxFloat64(avgG, avgB))
	minChannel := minFloat64(avgR, minFloat64(avgG, avgB))
	return (maxChannel - minChannel) > threshold
}

// detectDocumentEdges performs basic document edge detection
func (oca *coreAnalyzer) detectDocumentEdges(img *loader.FastImage) bool {
	// Zero-copy conversion to gray for manual pixel check
	grayImg, err := img.ConvertToGrayscale()
	if err != nil {
		return false
	}
	
	meta, err := grayImg.Metadata()
	if err != nil {
		return false
	}
	width, height := meta.Size.Width, meta.Size.Height
	
	buf := grayImg.Buffer()
	
	// Access pixel at (x, y) assuming stride=width
	at := func(x, y int) uint8 {
		if x < 0 || x >= width || y < 0 || y >= height {
			return 0
		}
		return buf[y*width+x]
	}

	// Simple heuristic: check if corners are significantly different from center
	corners := []image.Point{
		{10, 10}, // Top-left
		{width - 10, 10}, // Top-right
		{10, height - 10}, // Bottom-left
		{width - 10, height - 10}, // Bottom-right
	}

	center := int(at(width/2, height/2))
	differentCorners := 0

	for _, corner := range corners {
		if corner.X >= 0 && corner.X < width && corner.Y >= 0 && corner.Y < height {
			cornerValue := int(at(corner.X, corner.Y))
			if abs(cornerValue-center) > 30 {
				differentCorners++
			}
		}
	}

	return differentCorners >= 2
}

// updatePerformanceStats updates internal performance statistics
func (oca *coreAnalyzer) updatePerformanceStats(duration time.Duration) {
	oca.mu.Lock()
	defer oca.mu.Unlock()

	oca.analysisCount++
	oca.totalProcessTime += duration
}

// GetPerformanceStats returns current performance statistics
func (oca *coreAnalyzer) GetPerformanceStats() (int64, time.Duration) {
	oca.mu.RLock()
	defer oca.mu.RUnlock()

	return oca.analysisCount, oca.totalProcessTime
}

// Helper functions
func maxFloat64(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

func minFloat64(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

// abs returns the absolute value of an integer
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// parseResolution parses the resolution string (e.g., "1920x1080") and returns width and height
func (oca *coreAnalyzer) parseResolution(resolution string) (int, int) {
	if resolution == "" {
		return 0, 0
	}

	parts := strings.Split(resolution, "x")
	if len(parts) != 2 {
		return 0, 0
	}

	width, err1 := strconv.Atoi(parts[0])
	height, err2 := strconv.Atoi(parts[1])

	if err1 != nil || err2 != nil {
		return 0, 0
	}

	return width, height
}

// Cleanup shuts down the analyzer and releases resources
func (oca *coreAnalyzer) Cleanup() error {
	if oca.workerPool != nil {
		oca.workerPool.Close()
	}
	return nil
}

// Close shuts down the analyzer and releases resources
func (oca *coreAnalyzer) Close() error {
	return oca.Cleanup()
}
