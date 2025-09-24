# ecwid-images-downloader [![codeql](https://github.com/turchenkoalex/ecwid-images-downloader/actions/workflows/codeql.yml/badge.svg)](https://github.com/turchenkoalex/ecwid-images-downloader/actions/workflows/codeql.yml)

Usage:

```shell
./ecwid-images-downloader -store SID
```

How to use:

```shell
Usage of ./ecwid-images-downloader-darwin-arm64:
  -download-dir string
    	Dir for download images (default: downloads/storeId)
  -include-names
    	Use product names in image file names
  -limit int
    	API v3 fetch limit (default 100)
  -parallelism int
    	Download parallelism (default 5)
  -skip-categories
    	Skip categories images
  -skip-downloaded
    	Skip images already present on disk
  -skip-products
    	Skip product images
  -store int
    	Store ID
  -use-combinations
    	Download combination images
  -token string
    	APIv3 token (if not provided, retrieve public token from Instant Site API)
  -verbose
    	Detailed logs
```
