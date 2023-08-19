package file_save

import (
	"errors"
	"mime/multipart"

	"github.com/cloudwego/hertz/pkg/app"
)

var (
	errMissingKey = errors.New("missing key")
)

// Extractor defines the function to get file from app.Request Context
type Extractor func(c *app.RequestContext) (*multipart.FileHeader, error)

func NewExtractor(key string) (Extractor, error) {
	if key == "" {
		return nil, errMissingKey
	}
	return func(c *app.RequestContext) (*multipart.FileHeader, error) {
		file, err := c.FormFile(key)
		if err != nil {
			return nil, err
		}
		return file, nil
	}, nil
}
