package handlers

import (
	"mime/multipart"
)

func Validate(config UploadConfig, image *multipart.FileHeader) string {
	contains := func(slice []string, item string) bool {
		for _, v := range slice {
			if v == item {
				return true
			}
		}
		return false
	}

	mimeType := image.Header.Get("Content-Type")

	if !contains(config.allowedTypes, mimeType) {
		return "Unsupported image type"
	}

	if image.Size < config.minImageSize {
		return "Image size is too small"
	}
	if image.Size > config.maxImageSize {
		return "Image size exceeds the maximum limit"
	}
	if image.Size == 0 {
		return "Image file is empty"
	}

	return ""
}
