package cmd

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"

	"github.com/turchenkoalex/ecwid-images-downloader/api"
	"github.com/turchenkoalex/ecwid-images-downloader/status"
)

func DownloadProducts(ctx context.Context, httpClient *http.Client, options Options, publicToken string, imagesChan chan api.Image, status *status.Reporter) {
	limit := options.FetchLimit
	offset := 0
	total := status.GetTotalProductsCount()

	for offset < total {
		select {
		case <-ctx.Done():
			return
		default:
			// continue processing
		}

		products, err := api.LoadProducts(httpClient, options.StoreID, publicToken, offset, limit)
		if err != nil {
			fmt.Println("Download interrupted", err)
			return
		}

		wg := sync.WaitGroup{}

		for _, product := range products.Items {
			select {
			case <-ctx.Done():
				return
			default:
				// continue processing
			}

			images := product.Images(options.IncludeNames)

			for _, image := range images {
				imagesChan <- image
				status.MarkImageAdded()
			}

			if options.UseCombinations {
				productId := product.ID
				productName := product.Name

				wg.Add(1)
				go func() {
					defer wg.Done()
					// Загрузим параллельно комбинации товара и поставим их картинки в очередь
					downloadCombinations(ctx, httpClient, productId, productName, options, publicToken, imagesChan, status)
				}()
			}

			status.MarkProductProcessed()
		}

		wg.Wait()
		offset += limit
	}

	status.MarkAllProductsScheduled()
}

func downloadCombinations(ctx context.Context, httpClient *http.Client, productId int, productName string, options Options, publicToken string, imagesChan chan api.Image, status *status.Reporter) {
	combinations, err := api.LoadProductCombinations(httpClient, options.StoreID, publicToken, productId)
	if err == nil {
		for _, combination := range combinations {
			select {
			case <-ctx.Done():
				return
			default:
				// continue processing
			}

			image := combination.Image(productId, productName, options.IncludeNames)
			if image != nil {
				imagesChan <- *image
				status.MarkImageAdded()
			}
		}
	}
}

func DownloadCategories(ctx context.Context, httpClient *http.Client, options Options, publicToken string, imagesChan chan api.Image, status *status.Reporter) {
	limit := options.FetchLimit
	offset := 0
	total := status.GetTotalCategoriesCount()

	for offset < total {
		select {
		case <-ctx.Done():
			return
		default:
			// continue processing
		}

		categories, err := api.LoadCategories(httpClient, options.StoreID, publicToken, offset, limit)
		if err != nil {
			fmt.Println("Download interrupted", err)
			return
		}

		for _, category := range categories.Items {
			select {
			case <-ctx.Done():
				return
			default:
				// continue processing
			}

			image := category.Image(options.IncludeNames)
			if image != nil {
				imagesChan <- *image
				status.MarkImageAdded()
			}
			status.MarkCategoryProcessed()
		}

		offset += limit
	}

	status.MarkAllCategoriesScheduled()
}

func DownloadImages(ctx context.Context, httpClient *http.Client, options Options, imagesChan chan api.Image, status *status.Reporter) {
	for image := range imagesChan {
		select {
		case <-ctx.Done():
			return
		default:
			// continue processing
		}

		err := downloadFile(httpClient, options.SkipDownloaded, image)

		success := err == nil
		status.MarkImageDownloaded(success)

		if err != nil {
			fmt.Printf("Error occurred while download image from %s to file %s\n", image.URL, image.FileName)
		} else {
			if options.Verbose {
				fmt.Printf("Downloaded image url: %s to file: %s\n", image.URL, image.FileName)
			}
		}
	}
}

func downloadFile(client *http.Client, skipPresent bool, image api.Image) error {
	if _, err := os.Stat(image.FileName); err == nil {
		if skipPresent {
			return nil
		}
	}

	response, err := client.Get(image.URL)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(response.Body)

	// create the directory if it does not exist
	if err := os.MkdirAll(image.Dir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory for image: %w", err)
	}

	filePath := fmt.Sprintf("%s/%s", image.Dir, image.FileName)

	outputFile, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer func(outputFile *os.File) {
		_ = outputFile.Close()
	}(outputFile)

	_, err = io.Copy(outputFile, response.Body)
	return err
}
