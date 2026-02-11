package service

import (
	"fmt"
	"strings"

	"github.com/davidbyttow/govips/v2/vips"
)

func ProcessImage(buffer []byte, quality int, format string) ([]byte, string, error) {
	// Import image from buffer
	img, err := vips.NewImageFromBuffer(buffer)
	if err != nil {
		return nil, "", fmt.Errorf("failed to load image: %v", err)
	}
	defer img.Close()

	// Strip metadata (optimization)
	img.RemoveMetadata()

	format = strings.ToLower(format)
	var out []byte
	var contentType string

	switch format {
	case "jpeg", "jpg":
		ep := vips.NewJpegExportParams()
		ep.Quality = quality
		out, _, err = img.ExportJpeg(ep)
		contentType = "image/jpeg"
	case "png":
		ep := vips.NewPngExportParams()
		ep.Compression = 9
		out, _, err = img.ExportPng(ep)
		contentType = "image/png"
	default: // webp as default
		ep := vips.NewWebpExportParams()
		ep.Quality = quality
		out, _, err = img.ExportWebp(ep)
		contentType = "image/webp"
	}

	if err != nil {
		return nil, "", fmt.Errorf("failed to export image: %v", err)
	}

	return out, contentType, nil
}

func IsSupported(mimetype string) bool {
	supported := map[string]bool{
		"image/jpeg": true,
		"image/jpg":  true,
		"image/png":  true,
		"image/webp": true,
	}
	return supported[strings.ToLower(mimetype)]
}
