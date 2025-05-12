// serpapi.go
package serpapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"
)

// SearchRequest 定义通用搜索参数
type SearchRequest struct {
	Query  string
	Engine string
	APIKey string
}

// 缓存结构体
type SearchResultCache struct {
	Result string
	Expiry time.Time
}

var (
	cache      = make(map[string]SearchResultCache)
	cacheMutex sync.RWMutex
)

func PerformSearch(req SearchRequest) (string, error) {
	cacheMutex.RLock()
	cached, found := cache[req.Query]
	cacheMutex.RUnlock()

	if found && cached.Expiry.After(time.Now()) {
		return cached.Result, nil
	}

	var resultStr string
	switch req.Engine {
	case EngineGoogle:
		resultStr, _ = GoogleSearch(req.APIKey, req.Query)
	case EngineBing:
		resultStr, _ = BingSearch(req.APIKey, req.Query)
	default:
		return "", fmt.Errorf("unsupported search engine: %s", req.Engine)
	}

	cacheMutex.Lock()
	cache[req.Query] = SearchResultCache{
		Result: resultStr,
		Expiry: time.Now().Add(5 * time.Minute),
	}
	cacheMutex.Unlock()

	//fmt.Printf("Raw Result for [%s][%s]:\n%s\n", req.Engine, req.Query, resultStr)

	return resultStr, nil
}

func GoogleSearch(apiKey string, query string) (string, error) {
	baseURL := "https://serpapi.com/search"

	params := url.Values{}
	params.Add("q", query)
	params.Add("engine", "google")
	params.Add("api_key", apiKey)

	fullURL := baseURL + "?" + params.Encode()
	resp, err := http.Get(fullURL)
	if err != nil {
		return "", fmt.Errorf("HTTP request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	if len(body) > 0 && body[0] == '<' {
		return "", fmt.Errorf("received HTML instead of JSON: possible invalid API key or rate limit exceeded")
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("failed to unmarshal JSON: %v, raw response: %s", err, string(body))
	}

	organicResults := result["organic_results"].([]interface{})
	var output bytes.Buffer

	for i, item := range organicResults {
		itemMap := item.(map[string]interface{})
		output.WriteString(fmt.Sprintf("[%d] %s\n%s\n%s\n\n",
			i+1,
			itemMap["title"],
			itemMap["link"],
			itemMap["snippet"],
		))
	}

	return output.String(), nil
}

func BingSearch(apiKey string, query string) (string, error) {
	baseURL := "https://serpapi.com/search"

	params := url.Values{}
	params.Add("q", query)
	params.Add("engine", "bing")
	params.Add("api_key", apiKey)

	fullURL := baseURL + "?" + params.Encode()
	resp, err := http.Get(fullURL)
	if err != nil {
		return "", fmt.Errorf("HTTP request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	if len(body) > 0 && body[0] == '<' {
		return "", fmt.Errorf("received HTML instead of JSON: possible invalid API key or rate limit exceeded")
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("failed to unmarshal JSON: %v, raw response: %s", err, string(body))
	}

	organicResults := result["organic_results"].([]interface{})
	var output bytes.Buffer

	for i, item := range organicResults {
		itemMap := item.(map[string]interface{})
		output.WriteString(fmt.Sprintf("[%d] %s\n%s\n%s\n\n",
			i+1,
			itemMap["title"],
			itemMap["link"],
			itemMap["snippet"],
		))
	}

	return output.String(), nil
}
