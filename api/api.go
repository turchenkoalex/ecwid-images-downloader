package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

var apiBaseURL = "https://app.ecwid.com/api/v3"

// Products - https://api-docs.ecwid.com/reference/products#response
type Products struct {
	Total  int
	Count  int
	Offset int
	Limit  int
	Items  []Product
}

// Product - https://api-docs.ecwid.com/reference/products#productentry
type Product struct {
	ID    int
	Name  string
	Media ProductMedia
}

// Product combination - https://api-docs.ecwid.com/reference/variations#response
type ProductCombination struct {
	ID                int
	CombinationNumber int
	ThumbnailUrl      string
	ImageUrl          string
	SmallThumbnailUrl string
	HdThumbnailUrl    string
	OriginalImageUrl  string
}

// ProductMedia - https://api-docs.ecwid.com/reference/products#productmedia
type ProductMedia struct {
	Images []ProductImage
}

// ProductImage - https://api-docs.ecwid.com/reference/products#productimage
type ProductImage struct {
	ID               string
	ImageOriginalURL string
	Image1500pxURL   string
	Image800pxURL    string
	Image400pxURL    string
	Image160pxURL    string
}

// Categories - https://api-docs.ecwid.com/reference/categories#response
type Categories struct {
	Total  int
	Count  int
	Offset int
	Limit  int
	Items  []Category
}

// Category - https://api-docs.ecwid.com/reference/categories#category
type Category struct {
	ID               int
	Name             string
	OriginalImageUrl string
}

// LoadProducts - load products from api v3
func LoadProducts(httpClient *http.Client, storeID int, apiToken string, offset int, limit int) (Products, error) {
	products := &Products{}
	url := buildProductsURL(storeID, apiToken, offset, limit)
	err := readJSON(httpClient, url, products)
	if err != nil {
		return *products, err
	}

	return *products, nil
}

// LoadCategories - load categories from api v3
func LoadCategories(httpClient *http.Client, storeID int, apiToken string, offset int, limit int) (Categories, error) {
	categories := &Categories{}
	url := buildCategoriesURL(storeID, apiToken, offset, limit)
	err := readJSON(httpClient, url, categories)
	if err != nil {
		return *categories, err
	}

	return *categories, nil
}

// LoadProductCombinations - load product combinations from api v3
func LoadProductCombinations(httpClient *http.Client, storeID int, apiToken string, productId int) ([]ProductCombination, error) {
	var productCombinations []ProductCombination
	url := buildProductCombinationsURL(storeID, apiToken, productId)
	err := readJSON(httpClient, url, &productCombinations)
	if err != nil {
		return productCombinations, err
	}

	return productCombinations, nil
}

// LoadProductsTotalCount - load products total count
func LoadProductsTotalCount(httpClient *http.Client, storeID int, apiToken string) (int, error) {
	products, err := LoadProducts(httpClient, storeID, apiToken, 0, 0)
	if err != nil {
		return 0, err
	}

	return products.Total, nil
}

// LoadCategoriesTotalCount - load categories total count
func LoadCategoriesTotalCount(httpClient *http.Client, storeID int, apiToken string) (int, error) {
	categories, err := LoadCategories(httpClient, storeID, apiToken, 0, 0)
	if err != nil {
		return 0, err
	}

	return categories.Total, nil
}

func buildProductsURL(storeID int, apiToken string, offset int, limit int) string {
	return fmt.Sprintf("%s/%d/products?token=%s&limit=%d&offset=%d", apiBaseURL, storeID, apiToken, limit, offset)
}

func buildProductCombinationsURL(storeID int, apiToken string, productId int) string {
	return fmt.Sprintf("%s/%d/products/%d/combinations?token=%s", apiBaseURL, storeID, productId, apiToken)
}

func buildCategoriesURL(storeID int, apiToken string, offset int, limit int) string {
	return fmt.Sprintf("%s/%d/categories?token=%s&limit=%d&offset=%d", apiBaseURL, storeID, apiToken, limit, offset)
}

func readJSON(httpClient *http.Client, url string, target interface{}) error {
	response, err := httpClient.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode == 403 {
		return errors.New("INVALID_TOKEN")
	}

	return json.NewDecoder(response.Body).Decode(target)
}
