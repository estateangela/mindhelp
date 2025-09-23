package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"

	"mindhelp-backend/internal/config"
	"mindhelp-backend/internal/dto"

	"golang.org/x/time/rate"
)

// GoogleMapsService Google Maps API 服務
type GoogleMapsService struct {
	config     *config.Config
	client     *http.Client
	rateLimiter *rate.Limiter
	cache      *Cache
	mu         sync.RWMutex
}

// Cache 簡單的記憶體快取
type Cache struct {
	data map[string]CacheEntry
	mu   sync.RWMutex
}

// CacheEntry 快取條目
type CacheEntry struct {
	Value     interface{}
	ExpiresAt time.Time
}

// NewGoogleMapsService 創建新的 Google Maps 服務
func NewGoogleMapsService(cfg *config.Config) *GoogleMapsService {
	// 設定速率限制：每秒最多 10 個請求，突發最多 20 個
	limiter := rate.NewLimiter(rate.Limit(10), 20)
	
	return &GoogleMapsService{
		config:      cfg,
		client:      &http.Client{Timeout: 30 * time.Second},
		rateLimiter: limiter,
		cache: &Cache{
			data: make(map[string]CacheEntry),
		},
	}
}

// Get 從快取獲取資料
func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	entry, exists := c.data[key]
	if !exists {
		return nil, false
	}
	
	if time.Now().After(entry.ExpiresAt) {
		delete(c.data, key)
		return nil, false
	}
	
	return entry.Value, true
}

// Set 設定快取資料
func (c *Cache) Set(key string, value interface{}, duration time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	c.data[key] = CacheEntry{
		Value:     value,
		ExpiresAt: time.Now().Add(duration),
	}
}

// GeocodeWithCache 帶快取的地理編碼
func (s *GoogleMapsService) GeocodeWithCache(ctx context.Context, req dto.GeocodeRequest) (*dto.GeocodeResponse, error) {
	// 生成快取鍵
	cacheKey := fmt.Sprintf("geocode:%s:%s:%s", req.Address, req.Language, req.Region)
	
	// 嘗試從快取獲取
	if cached, found := s.cache.Get(cacheKey); found {
		if response, ok := cached.(*dto.GeocodeResponse); ok {
			return response, nil
		}
	}

	// 速率限制
	if err := s.rateLimiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limit error: %w", err)
	}

	// 檢查 API Key
	if s.config.GoogleMaps.APIKey == "" {
		return nil, fmt.Errorf("Google Maps API Key not configured")
	}

	// 構建請求
	params := url.Values{}
	params.Add("address", req.Address)
	params.Add("key", s.config.GoogleMaps.APIKey)
	
	if req.Language != "" {
		params.Add("language", req.Language)
	}
	if req.Region != "" {
		params.Add("region", req.Region)
	}

	requestURL := fmt.Sprintf("%s?%s", s.config.GoogleMaps.GeocodingURL, params.Encode())

	// 發送請求
	httpReq, err := http.NewRequestWithContext(ctx, "GET", requestURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := s.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// 讀取回應
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// 解析回應
	var geocodeResp dto.GeocodeResponse
	if err := json.Unmarshal(body, &geocodeResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// 檢查 API 狀態
	if geocodeResp.Status != "OK" && geocodeResp.Status != "ZERO_RESULTS" {
		return nil, fmt.Errorf("Google Maps API error: %s", geocodeResp.Status)
	}

	// 快取結果（快取 1 小時）
	s.cache.Set(cacheKey, &geocodeResp, time.Hour)

	return &geocodeResp, nil
}

// ReverseGeocodeWithCache 帶快取的反向地理編碼
func (s *GoogleMapsService) ReverseGeocodeWithCache(ctx context.Context, req dto.ReverseGeocodeRequest) (*dto.GeocodeResponse, error) {
	// 生成快取鍵
	cacheKey := fmt.Sprintf("reverse_geocode:%.6f,%.6f:%s:%s:%s", 
		req.Latitude, req.Longitude, req.Language, req.ResultType, req.LocationType)
	
	// 嘗試從快取獲取
	if cached, found := s.cache.Get(cacheKey); found {
		if response, ok := cached.(*dto.GeocodeResponse); ok {
			return response, nil
		}
	}

	// 速率限制
	if err := s.rateLimiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limit error: %w", err)
	}

	// 檢查 API Key
	if s.config.GoogleMaps.APIKey == "" {
		return nil, fmt.Errorf("Google Maps API Key not configured")
	}

	// 構建請求
	params := url.Values{}
	params.Add("latlng", fmt.Sprintf("%.8f,%.8f", req.Latitude, req.Longitude))
	params.Add("key", s.config.GoogleMaps.APIKey)
	
	if req.Language != "" {
		params.Add("language", req.Language)
	}
	if req.ResultType != "" {
		params.Add("result_type", req.ResultType)
	}
	if req.LocationType != "" {
		params.Add("location_type", req.LocationType)
	}

	requestURL := fmt.Sprintf("%s?%s", s.config.GoogleMaps.GeocodingURL, params.Encode())

	// 發送請求
	httpReq, err := http.NewRequestWithContext(ctx, "GET", requestURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := s.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// 讀取回應
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// 解析回應
	var geocodeResp dto.GeocodeResponse
	if err := json.Unmarshal(body, &geocodeResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// 檢查 API 狀態
	if geocodeResp.Status != "OK" && geocodeResp.Status != "ZERO_RESULTS" {
		return nil, fmt.Errorf("Google Maps API error: %s", geocodeResp.Status)
	}

	// 快取結果（快取 24 小時，座標不太會變）
	s.cache.Set(cacheKey, &geocodeResp, 24*time.Hour)

	return &geocodeResp, nil
}

// BatchGeocode 批次地理編碼
func (s *GoogleMapsService) BatchGeocode(ctx context.Context, addresses []string) ([]dto.GeocodeResponse, error) {
	if len(addresses) > 100 {
		return nil, fmt.Errorf("too many addresses: maximum 100 allowed")
	}

	// 使用 goroutine 並行處理，但限制並發數
	semaphore := make(chan struct{}, 5) // 最多同時 5 個請求
	var wg sync.WaitGroup
	results := make([]dto.GeocodeResponse, len(addresses))
	errors := make([]error, len(addresses))

	for i, address := range addresses {
		wg.Add(1)
		go func(index int, addr string) {
			defer wg.Done()
			
			// 獲取信號量
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			req := dto.GeocodeRequest{
				Address:  addr,
				Language: "zh-TW",
				Region:   "tw",
			}

			resp, err := s.GeocodeWithCache(ctx, req)
			if err != nil {
				errors[index] = err
				return
			}

			results[index] = *resp
		}(i, address)
	}

	wg.Wait()

	// 檢查是否有錯誤
	var hasError bool
	for _, err := range errors {
		if err != nil {
			hasError = true
			break
		}
	}

	if hasError {
		return nil, fmt.Errorf("some geocoding requests failed")
	}

	return results, nil
}

// SearchNearbyMentalHealthServices 搜尋附近心理健康服務（增強版）
func (s *GoogleMapsService) SearchNearbyMentalHealthServices(ctx context.Context, latitude, longitude float64, radius int, keyword string) (*dto.PlacesSearchResponse, error) {
	// 生成快取鍵
	cacheKey := fmt.Sprintf("mental_health:%.6f,%.6f:%d:%s", latitude, longitude, radius, keyword)
	
	// 嘗試從快取獲取
	if cached, found := s.cache.Get(cacheKey); found {
		if response, ok := cached.(*dto.PlacesSearchResponse); ok {
			return response, nil
		}
	}

	// 構建多個搜尋查詢
	queries := []string{
		fmt.Sprintf("%s 台灣", keyword),
		"精神科診所 台灣",
		"心理健康中心 台灣",
		"諮商中心 台灣",
		"心理治療 台灣",
	}

	var allResults []dto.PlaceResult
	location := fmt.Sprintf("%.8f,%.8f", latitude, longitude)

	// 對每個查詢進行搜尋
	for _, query := range queries {
		// 速率限制
		if err := s.rateLimiter.Wait(ctx); err != nil {
			continue
		}

		// 構建請求 URL
		baseURL := s.config.GoogleMaps.PlacesURL + "/textsearch/json"
		params := url.Values{}
		params.Add("query", query)
		params.Add("location", location)
		params.Add("radius", strconv.Itoa(radius))
		params.Add("type", "health")
		params.Add("language", "zh-TW")
		params.Add("region", "tw")
		params.Add("key", s.config.GoogleMaps.APIKey)

		requestURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

		// 發送請求
		httpReq, err := http.NewRequestWithContext(ctx, "GET", requestURL, nil)
		if err != nil {
			continue
		}

		resp, err := s.client.Do(httpReq)
		if err != nil {
			continue
		}

		// 讀取回應
		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			continue
		}

		// 解析回應
		var placesResp dto.PlacesSearchResponse
		if err := json.Unmarshal(body, &placesResp); err != nil {
			continue
		}

		// 如果查詢成功，添加結果
		if placesResp.Status == "OK" {
			allResults = append(allResults, placesResp.Results...)
		}
	}

	// 去重複結果（根據 PlaceID）
	uniqueResults := make(map[string]dto.PlaceResult)
	for _, result := range allResults {
		if _, exists := uniqueResults[result.PlaceID]; !exists {
			uniqueResults[result.PlaceID] = result
		}
	}

	// 轉換為切片並排序（按評分排序）
	var finalResults []dto.PlaceResult
	for _, result := range uniqueResults {
		finalResults = append(finalResults, result)
	}

	// 簡單按評分排序（降序）
	for i := 0; i < len(finalResults)-1; i++ {
		for j := i + 1; j < len(finalResults); j++ {
			if finalResults[i].Rating < finalResults[j].Rating {
				finalResults[i], finalResults[j] = finalResults[j], finalResults[i]
			}
		}
	}

	// 限制結果數量
	if len(finalResults) > 20 {
		finalResults = finalResults[:20]
	}

	response := &dto.PlacesSearchResponse{
		Results: finalResults,
		Status:  "OK",
	}

	// 快取結果（快取 30 分鐘）
	s.cache.Set(cacheKey, response, 30*time.Minute)

	return response, nil
}

// GetDirectionsWithCache 帶快取的路線規劃
func (s *GoogleMapsService) GetDirectionsWithCache(ctx context.Context, req dto.DirectionsRequest) (*dto.DirectionsResponse, error) {
	// 生成快取鍵
	cacheKey := fmt.Sprintf("directions:%s:%s:%s:%s:%t", 
		req.Origin, req.Destination, req.Mode, req.Language, req.Alternatives)
	
	// 嘗試從快取獲取
	if cached, found := s.cache.Get(cacheKey); found {
		if response, ok := cached.(*dto.DirectionsResponse); ok {
			return response, nil
		}
	}

	// 速率限制
	if err := s.rateLimiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limit error: %w", err)
	}

	// 檢查 API Key
	if s.config.GoogleMaps.APIKey == "" {
		return nil, fmt.Errorf("Google Maps API Key not configured")
	}

	// 構建請求 URL
	params := url.Values{}
	params.Add("origin", req.Origin)
	params.Add("destination", req.Destination)
	params.Add("key", s.config.GoogleMaps.APIKey)
	
	if req.Mode != "" {
		params.Add("mode", req.Mode)
	}
	if req.Language != "" {
		params.Add("language", req.Language)
	}
	if req.Region != "" {
		params.Add("region", req.Region)
	}
	if req.Alternatives {
		params.Add("alternatives", "true")
	}
	if req.Avoid != "" {
		params.Add("avoid", req.Avoid)
	}
	if req.Units != "" {
		params.Add("units", req.Units)
	}

	requestURL := fmt.Sprintf("%s?%s", s.config.GoogleMaps.DirectionsURL, params.Encode())

	// 發送請求
	httpReq, err := http.NewRequestWithContext(ctx, "GET", requestURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := s.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// 讀取回應
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// 解析回應
	var directionsResp dto.DirectionsResponse
	if err := json.Unmarshal(body, &directionsResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// 檢查 API 狀態
	if directionsResp.Status != "OK" && directionsResp.Status != "ZERO_RESULTS" {
		return nil, fmt.Errorf("Google Maps API error: %s", directionsResp.Status)
	}

	// 快取結果（快取 15 分鐘，交通狀況會變化）
	s.cache.Set(cacheKey, &directionsResp, 15*time.Minute)

	return &directionsResp, nil
}

// ClearCache 清除快取
func (s *GoogleMapsService) ClearCache() {
	s.cache.mu.Lock()
	defer s.cache.mu.Unlock()
	s.cache.data = make(map[string]CacheEntry)
}

// GetCacheStats 獲取快取統計
func (s *GoogleMapsService) GetCacheStats() map[string]int {
	s.cache.mu.RLock()
	defer s.cache.mu.RUnlock()
	
	stats := map[string]int{
		"total_entries": len(s.cache.data),
		"expired_entries": 0,
	}
	
	now := time.Now()
	for _, entry := range s.cache.data {
		if now.After(entry.ExpiresAt) {
			stats["expired_entries"]++
		}
	}
	
	return stats
}
