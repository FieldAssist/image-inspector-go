# Go Image Analyzer

Go Image Analyzer is a web server application written in Go that fetches images from specified URLs and analyzes them for quality and content. It's built with a modular architecture, making it easy to extend and maintain.

## Features

-   **Secure Image Fetching**: Fetches images from URLs with built-in SSRF protection.
-   **Image Quality Analysis**:
    -   Blurriness detection (Laplacian variance)
    -   Overexposure and brightness assessment
    -   Average luminance and saturation levels
-   **Content Analysis**:
    -   Simulated Optical Character Recognition (OCR) to extract text.
-   **Simplified API**: A single, consolidated endpoint for all analysis types.
-   **Modular Architecture**: Easily extensible with new `FeatureAnalyzer`s.

## Prerequisites

-   [Go](https://golang.org/doc/install) 1.23 or higher
-   [Docker](https://docs.docker.com/get-docker/) (optional, for containerization)

## Installation

1.  Clone the repository:
    ```sh
    git clone https://github.com/anime-shed/image-inspector-go.git
    cd image-inspector-go
    ```
2.  Build the application:
    ```sh
    go build -o image-inspector-go ./cmd/api
    ```

## Usage

1.  Run the application:
    ```sh
    go run ./cmd/api/main.go
    ```

2.  The server will start and listen on port `8080`. You can interact with the API using tools like `curl`.

## API Endpoint

The application has been refactored to use a single, unified endpoint for all image analysis.

### `POST /analyze`

Analyzes an image from a given URL and runs the requested analyses.

#### Request Body (JSON)

| Field         | Type    | Description                                       | Required |
| ------------- | ------- | ------------------------------------------------- | -------- |
| `url`         | `string`  | The URL of the image to be analyzed.              | Yes      |
| `with_quality`| `boolean` | If `true`, a quality analysis will be performed.  | No       |
| `with_ocr`    | `boolean` | If `true`, an OCR analysis will be performed.     | No       |

---

### Usage Example

Analyze an image for both quality and OCR content.

```bash
curl -X POST http://localhost:8080/analyze \
  -H "Content-Type: application/json" \
  -d '{
    "url": "https://example.com/image.jpg",
    "with_quality": true,
    "with_ocr": true
  }'
```

### Sample Response

```json
{
    "image_id": "",
    "has_quality_report": true,
    "blurriness": 437.67,
    "overexposure": 106.22,
    "avg_luminance": 0.45,
    "avg_saturation": 0.24,
    "has_ocr_report": true,
    "extracted_text": "Simulated OCR text from image",
    "has_qr_report": false
}
```

## Architecture Overview

The analyzer is designed with a modular approach:

-   **`CoreAnalyzer`**: The central orchestrator that manages the analysis process. It takes an image and a set of requested analyses.
-   **`FeatureAnalyzer`**: An interface for specialized, single-purpose analyzers. Current implementations include:
    -   `QualityAnalyzer`: Assesses image quality (blur, exposure, etc.).
    -   `OCRAnalyzer`: Performs (simulated) text extraction.
-   **`Service Layer`**: Orchestrates the business logic, fetching the image and delegating to the `CoreAnalyzer`.
-   **`Delivery Layer`**: Handles HTTP requests and responses using the Gin framework.

This design allows new features (e.g., QR code detection, face detection) to be added by simply creating a new `FeatureAnalyzer` and integrating it into the `CoreAnalyzer`.

## Core Dependencies

-   **[Gin Web Framework](https://github.com/gin-gonic/gin)**: High-performance HTTP web framework.
-   **[GORM](https://gorm.io/)**: A developer-friendly ORM for Go.
-   **[SQLite](https://github.com/mattn/go-sqlite3)**: Used as the database driver for GORM.
-   **[Testify](https://github.com/stretchr/testify)**: Testing toolkit for Go.
-   **[Go standard library](https://pkg.go.dev/)**: For image processing and core functionalities.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
