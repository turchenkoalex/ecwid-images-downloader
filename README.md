# ecwid-images-downloader [![codeql](https://github.com/turchenkoalex/ecwid-images-downloader/actions/workflows/codeql.yml/badge.svg)](https://github.com/turchenkoalex/ecwid-images-downloader/actions/workflows/codeql.yml)

Usage:
```shell
./ecwid-images-downloader -store SID -public-token public_TOKEN
```

All args:
```shell
Usage of ./ecwid-images-downloader:
  -download-dir string
    	Dir for download images (default: downloads/storeId)
  -limit int
    	API v3 fetch limit (default 100)
  -parallelism int
    	Download parallelism (default 5)
  -public-token string
    	Public Access Token. You can get it from storefront Ecwid.getAppPublicToken('ecwid-storefront')
  -skip-downloaded
    	Skip images already present on disk
  -store int
    	Store ID
  -use-combinations
    	Download combination images
  -verbose
    	Detailed logs
```
