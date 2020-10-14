package api

import (
	"fmt"
)

// Image data
type Image struct {
	FileName string
	URL      string
}

// Image - extract all available images from products structure
func (product Product) Images() []Image {
	var images []Image
	for _, image := range product.Media.Images {
		var downloadableImage Image
		downloadableImage.FileName = fmt.Sprintf("p%d-%s.jpg", product.ID, image.ID)

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
func (combination ProductCombination) Image(productId int) *Image {
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

	image.FileName = fmt.Sprintf("p%d-c%d.jpg", productId, combination.CombinationNumber)

	return &image
}

// Image - get category image
func (category Category) Image() *Image {
	if category.OriginalImageUrl == "" {
		return nil
	}

	var downloadableImage Image
	downloadableImage.FileName = fmt.Sprintf("cat%d.jpg", category.ID)
	downloadableImage.URL = category.OriginalImageUrl
	return &downloadableImage
}
