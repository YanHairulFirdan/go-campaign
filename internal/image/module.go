package image

import (
	"github.com/gofiber/fiber/v2"
	v1 "go-campaign.com/internal/image/transport/http/v1"
	"go-campaign.com/pkg/filesystem"
)

type HTTPDeps struct {
	FileSystem filesystem.Filesystem
}

func BootHttpV1(router fiber.Router, deps HTTPDeps) {
	handler := v1.NewImageHandler(deps.FileSystem)

	v1.RegisterRoute(router, handler)
}
