package main

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"os"
)

func GetAssets(attributes []string) ([]image.Image, error) {
	// Get list of assets for the image.
	paths := []string{
		fmt.Sprintf("./var/assets/backgrounds/Bg-%s.png", attributes[1]),
		fmt.Sprintf("./var/assets/bases/Base-%s.png", attributes[6]),
		fmt.Sprintf("./var/assets/outfits/Outfit%s.png", attributes[5]),
		fmt.Sprintf("./var/assets/hairs/Hair-%s.png", attributes[3]),
		fmt.Sprintf("./var/assets/eyes/Eyes-%s.png", attributes[2]),
		fmt.Sprintf("./var/assets/lips/Lips-%s.png", attributes[4]),
		fmt.Sprintf("./var/assets/accessories/Acc-%s.png", attributes[0]),
	}
	var assets []image.Image
	for _, path := range paths {
		fmt.Println(path)

		f, err := os.Open(path)
		if err != nil {
			panic("Error opening an asset file.")
		}
		defer f.Close()

		asset, err := png.Decode(f)
		if err != nil {
			panic("Error decoding an asset file.")
		}

		assets = append(assets, asset)
	}

	return assets, nil
}

func CreateImage(assets []image.Image) (image.Image, error) {
	// Build the image using the assets.
	img := image.NewRGBA(image.Rect(0, 0, 1336, 1336))
	for _, asset := range assets {
		draw.Draw(img, asset.Bounds(), asset, image.ZP, draw.Over)
	}

	return img, nil
}

func CreatePNGFile(img image.Image, name string) (*os.File, error) {
	// Create file output.
	file, err := os.Create(fmt.Sprintf("./var/images/%s.png", name))
	if err != nil {
		panic("Error creating the file.")
	}

	// Encode the image to png in the file.
	err = png.Encode(file, img)
	if err != nil {
		panic("Error encoding the file.")
	}

	return file, nil
}
