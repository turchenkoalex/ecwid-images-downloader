package status

import (
	"fmt"
	"sync/atomic"
	"time"
)

type Reporter struct {
	imageDownloadErrors      int32
	imageDownloadSuccess     int32
	imageTotalCount          int32
	productsCount            int
	productsProcessedCount   int32
	categoriesCount          int
	categoriesProcessedCount int32
	allCategoriesScheduled   int32
	allProductsScheduled     int32
	done                     chan interface{}
}

func CreateReporter(productsCount int, categoriesCount int) *Reporter {
	return &Reporter{
		imageDownloadErrors:      0,
		imageDownloadSuccess:     0,
		imageTotalCount:          0,
		productsCount:            productsCount,
		productsProcessedCount:   0,
		categoriesCount:          categoriesCount,
		categoriesProcessedCount: 0,
		allCategoriesScheduled:   0,
		allProductsScheduled:     0,
		done:                     make(chan interface{}),
	}
}

func (status *Reporter) GetTotalProductsCount() int {
	return status.productsCount
}

func (status *Reporter) GetTotalCategoriesCount() int {
	return status.categoriesCount
}

func (status *Reporter) MarkImageAdded() {
	atomic.AddInt32(&status.imageTotalCount, 1)
}

func (status *Reporter) MarkProductProcessed() {
	atomic.AddInt32(&status.productsProcessedCount, 1)
}

func (status *Reporter) MarkCategoryProcessed() {
	atomic.AddInt32(&status.categoriesProcessedCount, 1)
}

func (status *Reporter) MarkAllCategoriesScheduled() {
	atomic.StoreInt32(&status.allCategoriesScheduled, 1)
}

func (status *Reporter) MarkAllProductsScheduled() {
	atomic.StoreInt32(&status.allProductsScheduled, 1)
}

func (status *Reporter) allImagesScheduled() bool {
	return 1 == atomic.LoadInt32(&status.allCategoriesScheduled) &&
		1 == atomic.LoadInt32(&status.allProductsScheduled)
}

func (status *Reporter) MarkImageDownloaded(success bool) {
	if success {
		atomic.AddInt32(&status.imageDownloadSuccess, 1)
	} else {
		atomic.AddInt32(&status.imageDownloadErrors, 1)
	}
}

func (status *Reporter) Start(duration time.Duration) {
	status.printStatus()

	go func() {
		ticker := time.NewTicker(duration)
		defer ticker.Stop()
		for {
			select {
			case <-status.done:
				return
			case <-ticker.C:
				status.printStatus()
			}
		}
	}()
}

func (status *Reporter) printStatus() {
	categoriesProcessedCount := atomic.LoadInt32(&status.categoriesProcessedCount)
	productsProcessedCount := atomic.LoadInt32(&status.productsProcessedCount)
	imagesProcessed := atomic.LoadInt32(&status.imageDownloadErrors) + atomic.LoadInt32(&status.imageDownloadSuccess)
	imagesTotalCount := atomic.LoadInt32(&status.imageTotalCount)

	categoriesPercent := float32(1)
	if status.categoriesCount > 0 {
		categoriesPercent = float32(categoriesProcessedCount) / float32(status.categoriesCount)
	}

	productsPercent := float32(1)
	if status.productsCount > 0 {
		productsPercent = float32(productsProcessedCount) / float32(status.productsCount)
	}

	allImagesScheduled := status.allImagesScheduled()

	imagesPercent := float32(0)
	if imagesTotalCount > 0 && allImagesScheduled {
		imagesPercent = float32(imagesProcessed) / float32(imagesTotalCount)
	}

	var imagesPercentString string
	if allImagesScheduled {
		imagesPercentString = fmt.Sprintf("%.f%%", imagesPercent*100)
	} else {
		imagesPercentString = "??%"
	}

	totalPercent := categoriesPercent*0.15 + productsPercent*0.15 + imagesPercent*0.7

	fmt.Printf("[%3.f%%]: Images %d of %d (%s). Categories %d of %d (%2.f%%). Products %d of %d (%2.f%%)\n",
		totalPercent*100,
		imagesProcessed,
		imagesTotalCount,
		imagesPercentString,
		categoriesProcessedCount,
		status.categoriesCount,
		categoriesPercent*100,
		productsProcessedCount,
		status.productsCount,
		productsPercent*100,
	)
}

func (status *Reporter) Done() {
	status.done <- nil
	close(status.done)

	fmt.Printf("[100%%]: Successfully downloaded: %d images, failed: %d images\n", status.imageDownloadSuccess, status.imageDownloadErrors)
}
