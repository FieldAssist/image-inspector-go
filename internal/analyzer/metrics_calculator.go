package analyzer

import (
	"math"

	"github.com/anime-shed/image-inspector-go/internal/loader"
	"gonum.org/v1/gonum/stat"
)

// metricsCalculator implements MetricsCalculator interface with hybrid Libvips/Go optimizations
type metricsCalculator struct {
}

// NewMetricsCalculator creates a new metrics calculator
func NewMetricsCalculator() MetricsCalculator {
	return &metricsCalculator{}
}

// CalculateBasicMetrics computes basic image metrics using resizing for speed
func (omc *metricsCalculator) CalculateBasicMetrics(img *loader.FastImage) (metrics, error) {
	// Optimization: Resize to a small fixed size (e.g., 100x100) to approximate average colors rapidly
	// This uses libvips AVX2 resize which is extremely fast
	small, err := img.Resize(100, 100)
	if err != nil {
		return metrics{}, err
	}
	
	meta, err := small.Metadata()
	if err != nil {
		return metrics{}, err
	}

	buf := small.Buffer()
	totalPixels := float64(meta.Size.Width * meta.Size.Height)
	if totalPixels == 0 {
		return metrics{}, nil
	}

	var totalLum, totalSat, totalR, totalG, totalB float64

	// Assuming RGB input (libvips usually gives sRGB)
	// Iterate through buffer (3 bytes per pixel for RGB)
	// Safety check
	channels := 3
	if len(buf) < int(totalPixels)*channels {
		channels = 1 // Grayscale fallback?
	}

	for i := 0; i < len(buf); i += channels {
		if i+2 >= len(buf) && channels >= 3 {
			break
		}

		var r, g, b float64
		if channels >= 3 {
			r = float64(buf[i]) / 255.0
			g = float64(buf[i+1]) / 255.0
			b = float64(buf[i+2]) / 255.0
		} else {
			val := float64(buf[i]) / 255.0
			r, g, b = val, val, val
		}

		// Calculate HSV (Luminance/Saturation approximation)
		_, s, v := omc.rgbToHSV(r, g, b)
		
		totalR += r
		totalG += g
		totalB += b
		totalSat += s
		totalLum += v
	}

	return metrics{
		avgLuminance:  totalLum / totalPixels,
		avgSaturation: totalSat / totalPixels,
		avgR:          totalR / totalPixels,
		avgG:          totalG / totalPixels,
		avgB:          totalB / totalPixels,
	}, nil
}

// CalculateLaplacianVariance computes Laplacian variance for blur detection
func (omc *metricsCalculator) CalculateLaplacianVariance(img *loader.FastImage) (float64, error) {
	// 1. Convert to grayscale using libvips
	gray, err := img.ConvertToGrayscale()
	if err != nil {
		return 0, err
	}

	// 2. Resize to fixed working resolution (e.g. 512px width) to normalize blur metric
	// and ensure consistent performance regardless of input size
	resized, err := gray.Resize(512, 512) // Aspect ratio might change, but fine for blur "score" check
	if err != nil {
		return 0, err
	}

	meta, err := resized.Metadata()
	if err != nil {
		return 0, err
	}
	width, height := meta.Size.Width, meta.Size.Height
	buf := resized.Buffer()

	// 3. Compute Laplacian on raw byte buffer
	// Kernel: [0, 1, 0; 1, -4, 1; 0, 1, 0]
	
	// We use a slice to store laplacian values for variance calc
	// Since 512x512 is small (262k), we can allocate or pool. 
	// For simplicity and speed, we just iterate.
	
	var values []float64
	// Pre-allocate
	values = make([]float64, 0, (width-2)*(height-2))

	stride := width
	
	for y := 1; y < height-1; y++ {
		rowOffset := y * stride
		prevRowOffset := (y - 1) * stride
		nextRowOffset := (y + 1) * stride
		
		for x := 1; x < width-1; x++ {
			// Access pixels directly from 1D buffer
			center := float64(buf[rowOffset+x])
			top    := float64(buf[prevRowOffset+x])
			bottom := float64(buf[nextRowOffset+x])
			left   := float64(buf[rowOffset+x-1])
			right  := float64(buf[rowOffset+x+1])

			laplacian := -4*center + top + bottom + left + right
			values = append(values, laplacian)
		}
	}

	if len(values) == 0 {
		return 0, nil
	}

	return stat.Variance(values, nil), nil
}

// CalculateBrightness computes average brightness
func (omc *metricsCalculator) CalculateBrightness(img *loader.FastImage) (float64, error) {
	// Optimization: Resize to 1x1!
	// The single pixel value represents the average brightness of the image (computed by libvips)
	
	// First convert to grayscale to ensure brightness
	gray, err := img.ConvertToGrayscale()
	if err != nil {
		return 0, err
	}
	
	tiny, err := gray.Resize(1, 1)
	if err != nil {
		return 0, err
	}

	buf := tiny.Buffer()
	if len(buf) == 0 {
		return 0, nil
	}

	return float64(buf[0]), nil
}

// DetectSkew placeholder - reimplement if needed using edge detection on resized image
func (omc *metricsCalculator) DetectSkew(img *loader.FastImage) (*float64, error) {
	// Simplified: return nil for now to prioritize core performance
	// Skew detection is expensive
	return nil, nil
}

// DetectContours placeholder
func (omc *metricsCalculator) DetectContours(img *loader.FastImage) (int, error) {
	return 0, nil // Disabled for high-perf mode
}

// rgbToHSV helper
func (omc *metricsCalculator) rgbToHSV(r, g, b float64) (h, s, v float64) {
	max := math.Max(r, math.Max(g, b))
	min := math.Min(r, math.Min(g, b))
	delta := max - min
	v = max
	if max == 0 {
		s = 0
	} else {
		s = delta / max
	}
	// Hue calculation skipped for perf if not strictly needed, or keep it
	if delta == 0 {
		h = 0
	} else if max == r {
		h = 60 * (((g - b) / delta) + 0)
	} else if max == g {
		h = 60 * (((b - r) / delta) + 2)
	} else {
		h = 60 * (((r - g) / delta) + 4)
	}
	if h < 0 {
		h += 360
	}
	return h, s, v
}
