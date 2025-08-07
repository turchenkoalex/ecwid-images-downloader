package api

import (
	"fmt"
	"regexp"
	"strings"
)

// Image data
type Image struct {
	FileName string
	Dir      string
	URL      string
}

// Images - extract all available images from products structure
func (product Product) Images(includeNames bool) []Image {
	var images []Image
	for _, image := range product.Media.Images {
		var downloadableImage Image

		if includeNames {
			// Use product name in image file name
			downloadableImage.FileName = fmt.Sprintf("p%d-%s-%s.jpg", product.ID, image.ID, sanitizeFilename(product.Name))
		} else {
			// Use product ID and image ID in file name
			downloadableImage.FileName = fmt.Sprintf("p%d-%s.jpg", product.ID, image.ID)
		}

		downloadableImage.Dir = "products"

		if image.ImageOriginalURL != "" {
			downloadableImage.URL = image.ImageOriginalURL
		} else if image.Image1500pxURL != "" {
			downloadableImage.URL = image.Image1500pxURL
		} else if image.Image800pxURL != "" {
			downloadableImage.URL = image.Image800pxURL
		} else if image.Image400pxURL != "" {
			downloadableImage.URL = image.Image400pxURL
		} else {
			downloadableImage.URL = image.Image160pxURL
		}

		if downloadableImage.URL != "" {
			images = append(images, downloadableImage)
		} else {
			fmt.Printf("Not found image for product %d\n", product.ID)
		}
	}
	return images
}

// Image - get image
func (combination ProductCombination) Image(productId int, productName string, includeNames bool) *Image {
	var image Image

	if combination.OriginalImageUrl != "" {
		image.URL = combination.OriginalImageUrl
	} else if combination.ImageUrl != "" {
		image.URL = combination.ImageUrl
	} else if combination.HdThumbnailUrl != "" {
		image.URL = combination.HdThumbnailUrl
	} else if combination.ThumbnailUrl != "" {
		image.URL = combination.ThumbnailUrl
	} else {
		image.URL = combination.SmallThumbnailUrl
	}

	if image.URL == "" {
		return nil
	}

	if includeNames {
		// Use product name in image file name
		image.FileName = fmt.Sprintf("p%d-c%d-%s.jpg", productId, combination.CombinationNumber, sanitizeFilename(productName))
	} else {
		// Use product ID and combination number in file name
		image.FileName = fmt.Sprintf("p%d-c%d.jpg", productId, combination.CombinationNumber)
	}

	image.Dir = "products"

	return &image
}

// Image - get category image
func (category Category) Image(includeName bool) *Image {
	if category.OriginalImageUrl == "" {
		return nil
	}

	var downloadableImage Image

	if includeName {
		// Use category name in image file name
		downloadableImage.FileName = fmt.Sprintf("cat%d-%s.jpg", category.ID, sanitizeFilename(category.Name))
	} else {
		// Use category ID in file name
		downloadableImage.FileName = fmt.Sprintf("cat%d.jpg", category.ID)
	}

	downloadableImage.URL = category.OriginalImageUrl
	downloadableImage.Dir = "categories"
	return &downloadableImage
}

const maxFilenameLength = 255

var invalidChars = regexp.MustCompile(`[<>:"/\\|?*\x00-\x1F#,]`)

func sanitizeFilename(name string) string {
	name = strings.TrimSpace(name)
	name = strings.ReplaceAll(name, " ", "_")
	name = invalidChars.ReplaceAllString(name, "")

	if len(name) > maxFilenameLength {
		name = name[:maxFilenameLength]
	}

	if name == "" {
		name = "unnamed"
	}

	return name
}
