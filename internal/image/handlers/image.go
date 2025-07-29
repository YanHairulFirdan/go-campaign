package handlers

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"go-campaign.com/internal/shared/http/response"
)

func Upload(c *fiber.Ctx) error {
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

	if len(files.File) == 0 {
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

	if _, err := os.Stat(physicalUploadDir); os.IsNotExist(err) {
		if err := os.MkdirAll(physicalUploadDir, os.ModePerm); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(
				response.NewErrorResponse(
					"error",
					"Failed to create upload directory",
					err.Error(),
				),
			)
		}
	}

	images := files.File["images"]

	uploadedImages := make([]string, 0, len(images))
	validationErrors := map[string]string{}
	baseFileURL := os.Getenv("BASE_FILE_URL")

	for index, fileHeader := range images {
		errors := Validate(config, fileHeader)
		if errors != "" {
			return c.Status(fiber.StatusInternalServerError).JSON(
				response.NewErrorResponse(
					"error",
					"Failed to validate image",
					errors,
				),
			)
		}
		if len(errors) > 0 {
			validationErrors[fmt.Sprintf("image.%d", index)] = errors
			continue
		}

		extension := strings.ToLower(fileHeader.Filename[strings.LastIndex(fileHeader.Filename, "."):])

		fileHeader.Filename = fmt.Sprintf("%d%s", time.Now().UnixNano(), extension)

		if err := c.SaveFile(fileHeader, fmt.Sprintf("%s/%s", physicalUploadDir, fileHeader.Filename)); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(
				response.NewErrorResponse(
					"error",
					"Failed to save image",
					err.Error(),
				),
			)
		}

		uploadedImages = append(uploadedImages, fmt.Sprintf("%s/%s/%s", baseFileURL, uploadDir, fileHeader.Filename))
	}

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

func Delete(c *fiber.Ctx) error {
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
