package v1

import (
	"context"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go-campaign.com/internal/shared/http/response"
	"go-campaign.com/pkg/filesystem"
)

type ImageHandler struct {
	filesystem filesystem.Filesystem
}

func NewImageHandler(filesystem filesystem.Filesystem) *ImageHandler {
	return &ImageHandler{filesystem: filesystem}
}

func (h *ImageHandler) Upload(c *fiber.Ctx) error {
	module := c.FormValue("module", "default")

	config, exists := availableModules[module]
	if !exists {
		return c.Status(fiber.StatusBadRequest).JSON(
			response.NewErrorResponse(
				"error",
				"Invalid module",
				"Module not found",
			),
		)
	}

	files, err := c.MultipartForm()

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			response.NewErrorResponse(
				"error",
				"Invalid upload request",
				err.Error(),
			),
		)
	}

	images := files.File[config.inputField]

	if len(images) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(
			response.NewErrorResponse(
				"error",
				"No files uploaded",
				"Please upload at least one file",
			),
		)
	}

	uploadDir := fmt.Sprintf("%s/%s", config.uploadDir, module)

	physicalUploadDir := fmt.Sprintf("./public/%s", uploadDir)

	if err := h.filesystem.CreateDirectory(c.Context(), physicalUploadDir); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			response.NewErrorResponse(
				"error",
				"Failed to create upload directory",
				err.Error(),
			),
		)
	}

	validationErrors := validateUploadedImages(config, images)

	if len(validationErrors) > 0 {
		validationErrorsSlice := []map[string]string{validationErrors}
		return c.Status(fiber.StatusUnprocessableEntity).JSON(
			response.NewValidationErrorResponse(
				"error",
				"Validation failed",
				validationErrorsSlice,
			),
		)
	}

	baseFileURL := os.Getenv("BASE_FILE_URL")

	if baseFileURL == "" {
		return c.Status(fiber.StatusInternalServerError).JSON(
			response.NewErrorResponse(
				"error",
				"Base file URL not set",
				"Please set the BASE_FILE_URL environment variable",
			),
		)
	}

	uploadedImages, err := h.uploadImages(c.Context(), uploadDir, physicalUploadDir, baseFileURL, images)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			response.NewErrorResponse(
				"error",
				"Failed to save image",
				err.Error(),
			),
		)
	}

	return c.Status(fiber.StatusOK).JSON(
		response.NewResponse(
			"success",
			"Images uploaded successfully",
			map[string]interface{}{
				"images": uploadedImages,
			},
		),
	)
}

func (h *ImageHandler) Delete(c *fiber.Ctx) error {
	// get base file URL from environment variable
	baseFileURL := os.Getenv("BASE_FILE_URL")
	if baseFileURL == "" {
		return c.Status(fiber.StatusInternalServerError).JSON(
			response.NewErrorResponse(
				"error",
				"Base file URL not set",
				"Please set the BASE_FILE_URL environment variable",
			),
		)
	}

	// get url from form value
	imageURL := c.FormValue("image")
	if imageURL == "" {
		return c.Status(fiber.StatusBadRequest).JSON(
			response.NewErrorResponse(
				"error",
				"Image URL not provided",
				"Please provide the image URL to delete",
			),
		)
	}

	// check if the image URL starts with the base file URL
	if !strings.HasPrefix(imageURL, baseFileURL) {
		return c.Status(fiber.StatusBadRequest).JSON(
			response.NewErrorResponse(
				"error",
				"Invalid image URL",
				"Image URL does not match the base file URL",
			),
		)
	}

	// remove the base file URL from the image URL to get the physical path
	physicalPath := strings.TrimPrefix(imageURL, baseFileURL)
	physicalPath = fmt.Sprintf("./public/%s", physicalPath)

	if err := os.Remove(physicalPath); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			response.NewErrorResponse(
				"error",
				"Failed to delete image",
				err.Error(),
			),
		)
	}

	return c.Status(fiber.StatusOK).JSON(
		response.NewResponse(
			"success",
			"Image deleted successfully",
			map[string]any{
				"url": imageURL,
			},
		),
	)
}

func validateUploadedImages(config UploadConfig, images []*multipart.FileHeader) map[string]string {
	var errors = make(map[string]string)

	for index, image := range images {
		if err := Validate(config, image); err != "" {
			errors[fmt.Sprintf("image.%d", index)] = err
		}
	}

	return errors
}

func (h *ImageHandler) uploadImages(
	ctx context.Context,
	uploadDir,
	absoluteDir,
	baseFileURL string,
	images []*multipart.FileHeader,
) ([]string, error) {
	var uploaded = make([]string, 0, len(images))

	for _, image := range images {
		filename := generateFileName(image.Filename)

		err := h.saveUploadedImage(ctx, image, fmt.Sprintf("%s/%s", absoluteDir, filename))

		if err != nil {
			return nil, err
		}

		uploaded = append(uploaded, fmt.Sprintf("%s/%s/%s", baseFileURL, uploadDir, filename))
	}

	return uploaded, nil
}

func (h *ImageHandler) saveUploadedImage(ctx context.Context, image *multipart.FileHeader, path string) error {
	src, err := image.Open()

	if err != nil {
		return err
	}

	defer src.Close()

	return h.filesystem.SaveFile(ctx, src, path)
}

func generateFileName(originalName string) string {
	ext := strings.ToLower(filepath.Ext(originalName))

	return fmt.Sprintf("%s%s", uuid.NewString(), ext)
}
