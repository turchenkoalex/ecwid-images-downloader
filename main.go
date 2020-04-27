package main

import (
	"fmt"
	"github.com/turchenkoalex/ecwid-images-downloader/cmd"
	"github.com/turchenkoalex/ecwid-images-downloader/status"
	"net/http"
	"sync"
	"time"

	"github.com/turchenkoalex/ecwid-images-downloader/api"
)

func main() {
	fmt.Println("Ecwid Image Downloader 0.1")

	options, err := cmd.ReadOptions()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Start downloading product images (combinations mode: %v) for store %d with token %s to dir %s. (parallelism: %d)\n",
		options.UseCombinations,
		options.StoreID,
		options.PublicToken,
		options.DownloadDir,
		options.Parallelism,
	)

	httpClient := &http.Client{Timeout: 15 * time.Second}

	totalProductCount, err := api.LoadProductsTotalCount(httpClient, options.StoreID, options.PublicToken)
	if err != nil {
		fmt.Println("Error occurred while calculate products count", err)
		return
	}

	totalCategoriesCount, err := api.LoadCategoriesTotalCount(httpClient, options.StoreID, options.PublicToken)
	if err != nil {
		fmt.Println("Error occurred while calculate categories count", err)
		return
	}

	if totalProductCount == 0 && totalCategoriesCount == 0 {
		fmt.Println("No products and categories found. Nothing to download. Empty store catalog?")
		return
	}

	if options.Verbose {
		fmt.Printf("Found %d products and %d categories\n", totalProductCount, totalCategoriesCount)
	}

	// Репортилка о текущем статусе
	reporter := status.CreateReporter(totalProductCount, totalCategoriesCount)

	// репортаем состояние каждые 5 секунд
	reporter.Start(5 * time.Second)

	// Это очередь для скачивания, сюда будем накидывать все картинки которые нужно качать
	imagesChan := make(chan api.Image, 20000)

	// Запускаем parallelism параллельных задач на скачивание картинок
	downloadsWG := &sync.WaitGroup{}
	downloadsWG.Add(options.Parallelism)
	for jobID := 1; jobID <= options.Parallelism; jobID++ {
		go func() {
			defer downloadsWG.Done()
			cmd.DownloadImages(httpClient, options, imagesChan, reporter)
		}()
	}

	wg := &sync.WaitGroup{}

	// загрузим все товары и поставим загрузку картинок в очередь imagesChan
	wg.Add(1)
	go func() {
		defer wg.Done()
		cmd.DownloadProducts(httpClient, options, imagesChan, reporter)
	}()

	// загрузим все категории и поставим загрузку картинок в очередь imagesChan
	wg.Add(1)
	go func() {
		defer wg.Done()
		cmd.DownloadCategories(httpClient, options, imagesChan, reporter)
	}()

	// Ждем когда очедь картинок будет наполнена
	wg.Wait()

	// так как мы больше не будем писать в imagesChan, закрываем его
	close(imagesChan)

	// и ждем когда завершатся все задачи на скачивание
	downloadsWG.Wait()

	// Завершили все работы, останвливаем репортилку и выводим финальное сообщение
	reporter.Done()
}
