package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

func RetrievePublicToken(httpClient *http.Client, storeId int64) string {
	url := fmt.Sprintf("https://app.ecwid.com/storefront/api/v1/%d/initial-data", storeId)

	body := bytes.NewBufferString(`{}`)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, body)
	if err != nil {
		return ""
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return ""
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return ""
	}

	var r tokenResponse
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&r); err != nil {
		return ""
	}

	if r.StoreProfile.Value.AppsSettings.PublicTokens == nil {
		return ""
	}

	token, ok := r.StoreProfile.Value.AppsSettings.PublicTokens["ecwid-storefront"]
	if !ok {
		return ""
	}

	return token
}

type tokenResponse struct {
	StoreProfile struct {
		Value struct {
			AppsSettings struct {
				PublicTokens map[string]string `json:"publicTokens"`
			} `json:"appsSettings"`
		} `json:"value"`
	} `json:"storeProfile"`
}
