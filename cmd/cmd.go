package cmd

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

type Options struct {
	PublicToken     string
	StoreID         int
	Parallelism     int
	FetchLimit      int
	UseCombinations bool
	Verbose         bool
	DownloadDir     string
	SkipDownloaded  bool
}

var options Options

func init() {
	flag.StringVar(&options.PublicToken, "public-token", "", "Public Access Token. You can get it from storefront Ecwid.getAppPublicToken('ecwid-storefront')")
	flag.IntVar(&options.StoreID, "store", 0, "Store ID")
	flag.IntVar(&options.Parallelism, "parallelism", 5, "Download parallelism")
	flag.IntVar(&options.FetchLimit, "limit", 100, "API v3 fetch limit")
	flag.BoolVar(&options.UseCombinations, "use-combinations", false, "Download combination images")
	flag.BoolVar(&options.Verbose, "verbose", false, "Detailed logs")
	flag.BoolVar(&options.SkipDownloaded, "skip-downloaded", false, "Skip images already present on disk")
	flag.StringVar(&options.DownloadDir, "download-dir", "", "Dir for download images (default: downloads/storeId)")
}

func ReadOptions() (Options, error) {
	flag.Parse()

	if options.PublicToken == "" {
		flag.Usage()
		return options, fmt.Errorf("please add public-token argument")
	}

	if options.StoreID == 0 {
		flag.Usage()
		return options, fmt.Errorf("please add store argument")
	}

	if options.Parallelism < 1 {
		options.Parallelism = 1
	}

	if options.Parallelism > 20 {
		options.Parallelism = 20
	}

	if options.FetchLimit < 1 {
		options.FetchLimit = 1
	}

	if options.FetchLimit > 100 {
		options.FetchLimit = 100
	}

	if options.DownloadDir == "" {
		options.DownloadDir = "downloads/" + strconv.Itoa(options.StoreID)
	}
	err := configureDirs(options.DownloadDir)
	if err != nil {
		return options, err
	}

	return options, nil
}

func configureDirs(downloadDir string) error {
	_ = os.MkdirAll(downloadDir, os.ModePerm)

	err := os.Chdir(downloadDir)
	if err != nil {
		return fmt.Errorf("Can't change active dirrectory to %s", downloadDir)
	}

	return nil
}
