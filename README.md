# Ecwid Images Downloader

[![Go](https://img.shields.io/badge/lang-Go-blue.svg)](https://go.dev/)
[![Release](https://img.shields.io/github/v/release/turchenkoalex/ecwid-images-downloader)](https://github.com/turchenkoalex/ecwid-images-downloader/releases)
[![codeql](https://github.com/turchenkoalex/ecwid-images-downloader/actions/workflows/codeql.yml/badge.svg)](https://github.com/turchenkoalex/ecwid-images-downloader/actions/workflows/codeql.yml)

A lightweight **Go command-line utility** for bulk downloading product and category images from an [Ecwid](https://www.ecwid.com/) store.  

It is useful for:

- **Store owners** — backup all your product images, migrate them to another system, or simply keep a local copy.  
- **Engineers & integrators** — include it in scripts, CI/CD pipelines, or migration tools when working with Ecwid data.  

---

## ✨ Features

- Download images of **products**, **categories**, and **combinations/variations**.  
- **Skip already downloaded** images to avoid duplicates.  
- **Parallel downloads** to speed things up.  
- **Custom API fetch limit** (control how many items are fetched per request).  
- **Optional product names** in file names via `-include-names`.  
- **Verbose logging** for debugging.  
- Automatic retrieval of **public tokens** (no manual copy-pasting needed in most cases).  

---

## 🚀 Quick Start

Download the latest release from [GitHub Releases](https://github.com/turchenkoalex/ecwid-images-downloader/releases) or build from source with Go.

### Basic usage

```bash
./ecwid-images-downloader -store 123456
```

By default:

- Images are saved under `downloads/{storeId}`.  
- If no token is provided, the tool will attempt to automatically obtain a public token from the Instant Site API.  

### Using an explicit API token

If you prefer to provide your **Ecwid API v3 token** directly:

```bash
./ecwid-images-downloader -store 123456 -token secret_ecwid_api_v3_xxxxxxxxx
```

---

## ⚙️ Command-line flags

```
Usage of ./ecwid-images-downloader:
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
  -token string
    	Token to access API v3 (if not provided, will try to retrieve public token)
  -use-combinations
    	Download combination images
  -verbose
    	Detailed logs
```

---

## 📂 Examples

- **Download everything automatically (token fetched for you):**
  ```bash
  ./ecwid-images-downloader -store 123456
  ```

- **Use explicit API v3 token:**
  ```bash
  ./ecwid-images-downloader -store 123456 -token ecwid_api_v3_xxx
  ```

- **Add product names to file names:**
  ```bash
  ./ecwid-images-downloader -store 123456 -include-names
  ```

- **Skip categories, download only product images:**
  ```bash
  ./ecwid-images-downloader -store 123456 -skip-categories
  ```

- **Skip already downloaded images (resume safely):**
  ```bash
  ./ecwid-images-downloader -store 123456 -skip-downloaded
  ```

- **Download combination/variation images:**
  ```bash
  ./ecwid-images-downloader -store 123456 -use-combinations
  ```

- **Increase parallelism to 10 downloads at a time:**
  ```bash
  ./ecwid-images-downloader -store 123456 -parallelism 10
  ```

---

## 🛠 Project Structure

```
.
├── api/        # API client for Ecwid
├── cmd/        # Command entrypoint
├── status/     # Helpers for progress/status
├── go.mod
└── main.go
```

---

## 📌 Use Cases

- **Backup** all store images locally.  
- **Migration** from Ecwid to another e-commerce platform.  
- **Audit & analysis** of store media (duplicates, formats, optimization).  
- **Automation** in CI/CD or migration pipelines.  

---

## 🔑 Requirements

- [Go](https://go.dev/) (if building from source).  
- Internet access to fetch images.  
- Write permissions for the chosen download directory.  
- Store ID (you can find it in your Ecwid admin).  
- Optional: Ecwid API v3 token (only required if you don’t want auto token retrieval).  

---

## 🤝 Contributing

Contributions are welcome!  

1. Fork the repo.  
2. Create a feature branch.  
3. Add tests (where relevant).  
4. Open a pull request.  

Ideas for improvements:

- Progress bar or GUI wrapper.  
- Direct upload to cloud storage (S3, GCS, etc.).  
- Enhanced retry logic and error handling.  

---

## 📄 License

This project is licensed under the **Apache License 2.0** - see the [LICENSE](LICENSE) file for details.

```
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```

---
