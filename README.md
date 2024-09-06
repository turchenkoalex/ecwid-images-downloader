# ecwid-images-downloader [![codeql](https://github.com/turchenkoalex/ecwid-images-downloader/actions/workflows/codeql.yml/badge.svg)](https://github.com/turchenkoalex/ecwid-images-downloader/actions/workflows/codeql.yml)

Usage:

```shell
./ecwid-images-downloader -store SID -public-token public_TOKEN
```

You can get public_TOKEN if you open the storefront and in the browser console execute command `Ecwid.getAppPublicToken('ecwid-storefront')`,
the result of the form `public_XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX` use instead of `public_TOKEN`.

All args:

```shell
Usage of ./ecwid-images-downloader-arm64:
  -download-dir string
    	Dir for download images (default: downloads/storeId)
  -limit int
    	API v3 fetch limit (default 100)
  -parallelism int
    	Download parallelism (default 5)
  -public-token string
    	Public Access Token. You can get it from storefront Ecwid.getAppPublicToken('ecwid-storefront')
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
  -verbose
    	Detailed logs
```
