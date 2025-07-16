package ui

import (
	"bytes"
	"embed"
	"image"
	_ "image/png" // Register PNG decoder
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed assets
var assetsFS embed.FS

// loadImage loads an image from the embedded filesystem
func loadImage(path string) *ebiten.Image {
	imgData, err := assetsFS.ReadFile(path)
	if err != nil {
		log.Printf("Failed to load image %s: %v", path, err)
		return nil
	}

	img, _, err := image.Decode(bytes.NewReader(imgData))
	if err != nil {
		log.Printf("Failed to decode image %s: %v", path, err)
		return nil
	}

	return ebiten.NewImageFromImage(img)
}
