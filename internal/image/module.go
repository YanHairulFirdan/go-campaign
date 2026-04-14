package image

import (
	"github.com/gofiber/fiber/v2"
	"go-campaign.com/internal/app"
	v1 "go-campaign.com/internal/image/transport/http/v1"
)

// type ImageModule struct {
// }

// var _ app.Module = (*ImageModule)(nil)

// func NewImageModule() *ImageModule {
// 	return &ImageModule{}
// }

func Boot(router fiber.Router, deps *app.Dependencies) {
	handler := v1.NewImageHandler(deps.FileSystem)

	v1.RegisterRoute(router, handler)
}
