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
	categoriesPercent := float32(1)
	if status.categoriesCount > 0 {
		categoriesPercent = float32(status.categoriesProcessedCount) / float32(status.categoriesCount)
	}

	productsPercent := float32(1)
	if status.productsCount > 0 {
		productsPercent = float32(status.productsProcessedCount) / float32(status.productsCount)
	}

	imagesPercent := float32(0)
	imagesProcessed := status.imageDownloadErrors + status.imageDownloadSuccess
	if status.imageTotalCount > 0 {
		imagesPercent = float32(imagesProcessed) / float32(status.imageTotalCount)
	}

	totalPercent := categoriesPercent*0.15 + productsPercent*0.15 + imagesPercent*0.7

	fmt.Printf("[%3.f%%]: Images %d of %d (%.f%%). Categories %d of %d (%2.f%%). Products %d of %d (%2.f%%)\n",
		totalPercent*100,
		imagesProcessed,
		status.imageTotalCount,
		imagesPercent*100,
		status.categoriesProcessedCount,
		status.categoriesCount,
		categoriesPercent*100,
		status.productsProcessedCount,
		status.productsCount,
		productsPercent*100,
	)
}

func (status *Reporter) Done() {
	status.done <- nil
	close(status.done)

	fmt.Printf("[100%%]: Successfully downloaded: %d images, failed: %d images\n", status.imageDownloadSuccess, status.imageDownloadErrors)
}
