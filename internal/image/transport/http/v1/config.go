package v1

type UploadConfig struct {
	minImageSize int64
	maxImageSize int64
	allowedTypes []string
	uploadDir    string
	inputField   string
}

var availableModules = map[string]UploadConfig{
	"campaign": {
		minImageSize: 1024,            // 1 KB
		maxImageSize: 5 * 1024 * 1024, // 5 MB
		allowedTypes: []string{"image/jpeg", "image/png", "image/gif"},
		uploadDir:    "uploads",
		inputField:   "images",
	},
	"default": {
		minImageSize: 512,              // 512 bytes
		maxImageSize: 10 * 1024 * 1024, // 10 MB
		allowedTypes: []string{"image/jpeg", "image/png", "image/gif", "image/webp"},
		uploadDir:    "uploads",
		inputField:   "image",
	},
}
