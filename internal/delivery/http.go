package delivery

import (
	"net/http"

	"github.com/anime-shed/image-inspector-go/internal/service"
	"github.com/gin-gonic/gin"
)

// AnalysisRequest defines the structure for the incoming JSON request.
// It allows clients to specify the image URL and the desired analyses.
type AnalysisRequest struct {
	URL         string `json:"url" binding:"required,url"`
	WithQuality bool   `json:"with_quality"`
	WithOCR     bool   `json:"with_ocr"`
}

// NewHandler sets up the routing for the application and returns the HTTP handler.
func NewHandler(analysisService *service.ImageAnalysisService) http.Handler {
	r := gin.Default()

	// A single, consolidated endpoint for all image analysis.
	r.POST("/analyze", analyzeImage(analysisService))

	return r
}

// analyzeImage is the handler function for the /analyze endpoint.
func analyzeImage(svc *service.ImageAnalysisService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req AnalysisRequest

		// 1. Decode and validate the incoming JSON request.
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
			return
		}

		// 2. Delegate the core business logic to the service layer.
		report, err := svc.Analyze(
			c.Request.Context(),
			req.URL,
			req.WithQuality,
			req.WithOCR,
		)

		// 3. Handle any errors returned from the service.
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Analysis failed", "details": err.Error()})
			return
		}

		// 4. Return the successful analysis report.
		c.JSON(http.StatusOK, report)
	}
}
