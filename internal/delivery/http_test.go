package delivery

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"image"
	"image/png"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/anime-shed/image-inspector-go/internal/analyzer"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockService is a mock implementation of the Service for handler tests.
type MockService struct {
	mock.Mock
}

// ProcessImage mocks the service call.
func (m *MockService) ProcessImage(ctx context.Context, r io.Reader, opts ...analyzer.ReportConfigurator) (string, error) {
	// Pass the actual arguments received to the mock framework.
	args := m.Called(ctx, r, opts)
	return args.String(0), args.Error(1)
}

func TestHTTPHandler(t *testing.T) {
	// Create a dummy image to be used in request bodies for sub-tests.
	var imgBuffer bytes.Buffer
	require.NoError(t, png.Encode(&imgBuffer, image.NewRGBA(image.Rect(0, 0, 1, 1))))
	imageBytes := imgBuffer.Bytes()

	t.Run("handleAnalyze with all options successfully", func(t *testing.T) {
		// Setup: Create a new mock and handler for this specific test case.
		mockService := new(MockService)
		httpHandler := NewHTTPHandler(mockService)
		router := mux.NewRouter()
		httpHandler.RegisterRoutes(router)

		// Setup: Prepare the HTTP request and recorder.
		targetURL := "/analyze?quality=true&ocr=true"
		req := httptest.NewRequest(http.MethodPost, targetURL, bytes.NewReader(imageBytes))
		rr := httptest.NewRecorder()

		// Expectation: The service's ProcessImage method will be called.
		expectedImageID := "test-id-123"
		mockService.On("ProcessImage",
			mock.Anything, // context
			mock.Anything, // io.Reader
			mock.MatchedBy(func(opts []analyzer.ReportConfigurator) bool { return len(opts) == 2 }),
		).Return(expectedImageID, nil)

		// Action: Serve the HTTP request.
		router.ServeHTTP(rr, req)

		// Assertions.
		assert.Equal(t, http.StatusCreated, rr.Code)
		var respBody map[string]string
		require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &respBody))
		assert.Equal(t, expectedImageID, respBody["image_id"])
		mockService.AssertExpectations(t)
	})

	t.Run("handleAnalyze returns error on service failure", func(t *testing.T) {
		mockService := new(MockService)
		httpHandler := NewHTTPHandler(mockService)
		router := mux.NewRouter()
		httpHandler.RegisterRoutes(router)

		req := httptest.NewRequest(http.MethodPost, "/analyze", bytes.NewReader(imageBytes))
		rr := httptest.NewRecorder()

		expectedError := errors.New("service failed")
		mockService.On("ProcessImage", mock.Anything, mock.Anything, mock.Anything).Return("", expectedError)

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("handleAnalyze returns bad request for invalid image format", func(t *testing.T) {
		mockService := new(MockService)
		httpHandler := NewHTTPHandler(mockService)
		router := mux.NewRouter()
		httpHandler.RegisterRoutes(router)

		invalidFormatErr := errors.New("failed to decode image: image: unknown format")
		mockService.On("ProcessImage", mock.Anything, mock.Anything, mock.Anything).Return("", invalidFormatErr)

		req := httptest.NewRequest(http.MethodPost, "/analyze", bytes.NewReader([]byte("not an image")))
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Contains(t, rr.Body.String(), "Invalid image format")
		mockService.AssertExpectations(t)
	})

	t.Run("RegisterRoutes registers the correct route and method", func(t *testing.T) {
		router := mux.NewRouter()
		httpHandler := NewHTTPHandler(nil) // Service can be nil for this test
		httpHandler.RegisterRoutes(router)

		// Check if a POST request to /analyze is matched.
		req := httptest.NewRequest(http.MethodPost, "/analyze", nil)
		var match mux.RouteMatch
		assert.True(t, router.Match(req, &match))
		assert.NotNil(t, match.Handler)

		// Check that a GET request is not matched.
		req = httptest.NewRequest(http.MethodGet, "/analyze", nil)
		assert.False(t, router.Match(req, &match))
	})
}
