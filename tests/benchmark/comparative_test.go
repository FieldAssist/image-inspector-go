package benchmark

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/anime-shed/image-inspector-go/internal/config"
	"github.com/anime-shed/image-inspector-go/internal/container"
	"github.com/anime-shed/image-inspector-go/internal/service"
)

// Target image for benchmarking
const targetImageURL = "https://imagedetectionv2.blob.core.windows.net/images/234279/image-recognition/1811845c9c104102ab4ef08ef65cb68b"

func setupService() (service.ImageAnalysisService, error) {
	// Load configuration
	cfg, err := config.LoadFromEnv()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %v", err)
	}

	// Initialize dependency injection container
	c, err := container.NewContainer(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize container: %v", err)
	}

	return c.GetAnalysisService(), nil
}

func BenchmarkTargetImage(b *testing.B) {
	svc, err := setupService()
	if err != nil {
		b.Fatalf("Setup failed: %v", err)
	}

	// Pre-warm / Ensure connectivity
	ctx := context.Background()
	_, err = svc.AnalyzeImage(ctx, targetImageURL, false)
	if err != nil {
		b.Logf("Warning: Pre-warm request failed (network issue?): %v", err)
		// We don't fail here to allow benchmark to attempt running, but it might be connectivity dependent
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Create a separate context for each request to simulate real traffic
		// We use a generous timeout to ensure we measure processing time, not timeout errors
		reqCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
		res, err := svc.AnalyzeImage(reqCtx, targetImageURL, false)
		cancel()

		if err != nil {
			b.Errorf("Analysis failed: %v", err)
		} else {
             b.Logf("Processing Time: %.6f sec", res.ProcessingTimeSec)
        }
	}
}
