package loader

import (
	"fmt"

	"github.com/h2non/bimg"
)

// FastImage wraps raw image bytes and provides high-performance operations via libvips
type FastImage struct {
	buffer []byte
	img    *bimg.Image
}

// NewFastImage creates a new FastImage from a byte buffer
func NewFastImage(data []byte) *FastImage {
	return &FastImage{
		buffer: data,
		img:    bimg.NewImage(data),
	}
}

// Metadata returns image metadata (width, height, type) efficiently
func (f *FastImage) Metadata() (bimg.ImageMetadata, error) {
	return f.img.Metadata()
}

// Resize resizes the image to the specified dimensions using libvips (AVX2 optimized)
func (f *FastImage) Resize(width, height int) (*FastImage, error) {
	newOptions := bimg.Options{
		Width:  width,
		Height: height,
		Embed:  true,
	}

	newBuf, err := f.img.Process(newOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to resize image: %w", err)
	}

	return NewFastImage(newBuf), nil
}

// ConvertToGrayscale converts the image to grayscale
func (f *FastImage) ConvertToGrayscale() (*FastImage, error) {
	// bimg doesn't have a direct "grayscale" method on Image, but can use Process option
	// StripMetadata=true often helps reduce size
	// Interpretation: bimg.InterpretationBW
	
	// Since bimg.Options doesn't expose Interpretation directly in all versions, 
	// we use Colourspace change or similar if available, or force JPEG/PNG grayscale save options.
	// A strictly correct way in libvips is generic "colourspace" op.
	// For bimg, we might need to assume it's implicit or use a workaround if simple API is missing.
	// Checking bimg docs: Options has `Interpretation`.
	
	newOptions := bimg.Options{
		Interpretation: bimg.InterpretationBW,
	}

	newBuf, err := f.img.Process(newOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to convert to grayscale: %w", err)
	}

	return NewFastImage(newBuf), nil
}

// Buffer returns the underlying byte buffer
func (f *FastImage) Buffer() []byte {
	return f.buffer
}

// Size returns the size of the buffer in bytes
func (f *FastImage) Size() int {
	return len(f.buffer)
}
