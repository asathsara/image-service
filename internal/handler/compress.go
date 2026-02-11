package handler

import (
	"image-service/internal/service"
	"io"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func CompressHandler(c *fiber.Ctx) error {
	// Parse multipart file
	fileHead, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "No file uploaded"})
	}

	// Validate mimetype
	if !service.IsSupported(fileHead.Header.Get("Content-Type")) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Unsupported media type. Allowed: jpeg, png, webp"})
	}

	// Get file buffer
	file, err := fileHead.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to open file"})
	}
	defer file.Close()

	buffer, err := io.ReadAll(file)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to read file"})
	}

	// Parse query params
	qualityRaw := c.Query("quality", "75")
	quality, err := strconv.Atoi(qualityRaw)
	if err != nil || quality < 1 || quality > 100 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Quality must be between 1 and 100"})
	}

	format := c.Query("format", "webp")

	// Process image
	output, contentType, err := service.ProcessImage(buffer, quality, format)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Send response
	c.Set("Content-Type", contentType)
	return c.Send(output)
}
