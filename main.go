package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/turchenkoalex/ecwid-images-downloader/cmd"
	"github.com/turchenkoalex/ecwid-images-downloader/status"

	"github.com/turchenkoalex/ecwid-images-downloader/api"
)

func main() {
	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	cmd.PrintVersion()

	options, err := cmd.ReadOptions()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
		return
	}

	if options.SkipProducts && options.SkipCategories {
		fmt.Println("Skip categories and products in same time not allowed")
		os.Exit(1)
		return
	}

	subject := ""
	if options.SkipCategories {
		subject = "products"
	} else if options.SkipProducts {
		subject = "categories"
	} else {
		subject = "products and categories"
	}

	httpClient := &http.Client{Timeout: 15 * time.Second}

	var apiToken string
	if len(options.Token) > 0 {
		apiToken = options.Token
	} else {
		apiToken = api.RetrievePublicToken(httpClient, options.StoreID)

		if len(apiToken) == 0 {
			fmt.Printf("Can't retrieve public token for store %d. Please check that store ID is correct and store has instant site or provide token manually.\n", options.StoreID)
			os.Exit(1)
			return
		}
	}

	if len(apiToken) == 0 {
		fmt.Printf("No token provided for store %d. Please provide token manually with -token argument.\n", options.StoreID)
		os.Exit(1)
		return
	}

	fmt.Printf("Start downloading %s images (combinations mode: %v) for store %d with token %s to dir %s. (parallelism: %d)\n",
		subject,
		options.UseCombinations,
		options.StoreID,
		apiToken,
		options.DownloadDir,
		options.Parallelism,
	)

	totalProductCount := 0
	totalCategoriesCount := 0

	if !options.SkipProducts {
		totalProductCount, err = api.LoadProductsTotalCount(httpClient, options.StoreID, apiToken)
		if err != nil {
			fmt.Println("Error occurred while calculate products count", err)
			return
		}
	}

	if !options.SkipCategories {
		totalCategoriesCount, err = api.LoadCategoriesTotalCount(httpClient, options.StoreID, apiToken)
		if err != nil {
			fmt.Println("Error occurred while calculate categories count", err)
			return
		}
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
			cmd.DownloadImages(ctx, httpClient, options, imagesChan, reporter)
		}()
	}

	wg := &sync.WaitGroup{}

	// загрузим все товары и поставим загрузку картинок в очередь imagesChan
	if !options.SkipProducts {
		wg.Add(1)
		go func() {
			defer wg.Done()
			cmd.DownloadProducts(ctx, httpClient, options, apiToken, imagesChan, reporter)
		}()
	} else {
		reporter.MarkAllProductsScheduled()
	}

	// загрузим все категории и поставим загрузку картинок в очередь imagesChan
	if !options.SkipCategories {
		wg.Add(1)
		go func() {
			defer wg.Done()
			cmd.DownloadCategories(ctx, httpClient, options, apiToken, imagesChan, reporter)
		}()
	} else {
		reporter.MarkAllProductsScheduled()
	}

	// Ждем когда очедь картинок будет наполнена
	wg.Wait()

	// так как мы больше не будем писать в imagesChan, закрываем его
	close(imagesChan)

	// и ждем когда завершатся все задачи на скачивание
	downloadsWG.Wait()

	// Завершили все работы, останвливаем репортилку и выводим финальное сообщение
	reporter.Done()
}
