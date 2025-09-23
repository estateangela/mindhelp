package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"mindhelp-backend/internal/config"
	"mindhelp-backend/internal/dto"
	"mindhelp-backend/internal/vo"

	"github.com/gin-gonic/gin"
)

// GoogleMapsHandler Google Maps API 處理器
type GoogleMapsHandler struct {
	config *config.Config
	client *http.Client
}

// NewGoogleMapsHandler 創建新的 Google Maps 處理器
func NewGoogleMapsHandler(cfg *config.Config) *GoogleMapsHandler {
	return &GoogleMapsHandler{
		config: cfg,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Geocode 地理編碼 - 地址轉經緯度
// @Summary 地理編碼
// @Description 將地址轉換為經緯度座標
// @Tags google-maps
// @Accept json
// @Produce json
// @Param request body dto.GeocodeRequest true "地理編碼請求"
// @Success 200 {object} dto.GeocodeResponse
// @Failure 400 {object} vo.ErrorResponse
// @Failure 500 {object} vo.ErrorResponse
// @Router /google-maps/geocode [post]
func (h *GoogleMapsHandler) Geocode(c *gin.Context) {
	var req dto.GeocodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, vo.ErrorResponse{
			Code:    "INVALID_REQUEST",
			Message: "無效的請求格式",
			Error:   err.Error(),
		})
		return
	}

	// 驗證請求
	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, vo.ErrorResponse{
			Code:    "VALIDATION_ERROR",
			Message: "請求驗證失敗",
			Error:   err.Error(),
		})
		return
	}

	// 檢查 API Key
	if h.config.GoogleMaps.APIKey == "" {
		c.JSON(http.StatusInternalServerError, vo.ErrorResponse{
			Code:    "MISSING_API_KEY",
			Message: "Google Maps API Key 未設定",
		})
		return
	}

	// 構建請求 URL
	baseURL := h.config.GoogleMaps.GeocodingURL
	params := url.Values{}
	params.Add("address", req.Address)
	params.Add("key", h.config.GoogleMaps.APIKey)
	
	if req.Language != "" {
		params.Add("language", req.Language)
	}
	if req.Region != "" {
		params.Add("region", req.Region)
	}

	requestURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	// 發送請求
	resp, err := h.client.Get(requestURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, vo.ErrorResponse{
			Code:    "GOOGLE_MAPS_ERROR",
			Message: "Google Maps API 請求失敗",
			Error:   err.Error(),
		})
		return
	}
	defer resp.Body.Close()

	// 讀取回應
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, vo.ErrorResponse{
			Code:    "READ_RESPONSE_ERROR",
			Message: "讀取回應失敗",
			Error:   err.Error(),
		})
		return
	}

	// 解析回應
	var geocodeResp dto.GeocodeResponse
	if err := json.Unmarshal(body, &geocodeResp); err != nil {
		c.JSON(http.StatusInternalServerError, vo.ErrorResponse{
			Code:    "PARSE_RESPONSE_ERROR",
			Message: "解析回應失敗",
			Error:   err.Error(),
		})
		return
	}

	// 檢查 Google Maps API 狀態
	if geocodeResp.Status != "OK" && geocodeResp.Status != "ZERO_RESULTS" {
		var errorResp dto.GoogleMapsError
		json.Unmarshal(body, &errorResp)
		c.JSON(http.StatusBadRequest, vo.ErrorResponse{
			Code:    "GOOGLE_MAPS_API_ERROR",
			Message: errorResp.ErrorMessage,
			Error:   geocodeResp.Status,
		})
		return
	}

	c.JSON(http.StatusOK, geocodeResp)
}

// ReverseGeocode 反向地理編碼 - 經緯度轉地址
// @Summary 反向地理編碼
// @Description 將經緯度座標轉換為地址
// @Tags google-maps
// @Accept json
// @Produce json
// @Param request body dto.ReverseGeocodeRequest true "反向地理編碼請求"
// @Success 200 {object} dto.GeocodeResponse
// @Failure 400 {object} vo.ErrorResponse
// @Failure 500 {object} vo.ErrorResponse
// @Router /google-maps/reverse-geocode [post]
func (h *GoogleMapsHandler) ReverseGeocode(c *gin.Context) {
	var req dto.ReverseGeocodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, vo.ErrorResponse{
			Code:    "INVALID_REQUEST",
			Message: "無效的請求格式",
			Error:   err.Error(),
		})
		return
	}

	// 驗證請求
	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, vo.ErrorResponse{
			Code:    "VALIDATION_ERROR",
			Message: "請求驗證失敗",
			Error:   err.Error(),
		})
		return
	}

	// 檢查 API Key
	if h.config.GoogleMaps.APIKey == "" {
		c.JSON(http.StatusInternalServerError, vo.ErrorResponse{
			Code:    "MISSING_API_KEY",
			Message: "Google Maps API Key 未設定",
		})
		return
	}

	// 構建請求 URL
	baseURL := h.config.GoogleMaps.GeocodingURL
	params := url.Values{}
	params.Add("latlng", fmt.Sprintf("%.8f,%.8f", req.Latitude, req.Longitude))
	params.Add("key", h.config.GoogleMaps.APIKey)
	
	if req.Language != "" {
		params.Add("language", req.Language)
	}
	if req.ResultType != "" {
		params.Add("result_type", req.ResultType)
	}
	if req.LocationType != "" {
		params.Add("location_type", req.LocationType)
	}

	requestURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	// 發送請求
	resp, err := h.client.Get(requestURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, vo.ErrorResponse{
			Code:    "GOOGLE_MAPS_ERROR",
			Message: "Google Maps API 請求失敗",
			Error:   err.Error(),
		})
		return
	}
	defer resp.Body.Close()

	// 讀取回應
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, vo.ErrorResponse{
			Code:    "READ_RESPONSE_ERROR",
			Message: "讀取回應失敗",
			Error:   err.Error(),
		})
		return
	}

	// 解析回應
	var geocodeResp dto.GeocodeResponse
	if err := json.Unmarshal(body, &geocodeResp); err != nil {
		c.JSON(http.StatusInternalServerError, vo.ErrorResponse{
			Code:    "PARSE_RESPONSE_ERROR",
			Message: "解析回應失敗",
			Error:   err.Error(),
		})
		return
	}

	// 檢查 Google Maps API 狀態
	if geocodeResp.Status != "OK" && geocodeResp.Status != "ZERO_RESULTS" {
		var errorResp dto.GoogleMapsError
		json.Unmarshal(body, &errorResp)
		c.JSON(http.StatusBadRequest, vo.ErrorResponse{
			Code:    "GOOGLE_MAPS_API_ERROR",
			Message: errorResp.ErrorMessage,
			Error:   geocodeResp.Status,
		})
		return
	}

	c.JSON(http.StatusOK, geocodeResp)
}

// SearchPlaces 搜尋地點
// @Summary 搜尋地點
// @Description 搜尋附近的心靈健康相關地點
// @Tags google-maps
// @Accept json
// @Produce json
// @Param request body dto.PlacesSearchRequest true "地點搜尋請求"
// @Success 200 {object} dto.PlacesSearchResponse
// @Failure 400 {object} vo.ErrorResponse
// @Failure 500 {object} vo.ErrorResponse
// @Router /google-maps/search-places [post]
func (h *GoogleMapsHandler) SearchPlaces(c *gin.Context) {
	var req dto.PlacesSearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, vo.ErrorResponse{
			Code:    "INVALID_REQUEST",
			Message: "無效的請求格式",
			Error:   err.Error(),
		})
		return
	}

	// 驗證請求
	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, vo.ErrorResponse{
			Code:    "VALIDATION_ERROR",
			Message: "請求驗證失敗",
			Error:   err.Error(),
		})
		return
	}

	// 檢查 API Key
	if h.config.GoogleMaps.APIKey == "" {
		c.JSON(http.StatusInternalServerError, vo.ErrorResponse{
			Code:    "MISSING_API_KEY",
			Message: "Google Maps API Key 未設定",
		})
		return
	}

	// 構建請求 URL
	baseURL := h.config.GoogleMaps.PlacesURL + "/textsearch/json"
	params := url.Values{}
	params.Add("key", h.config.GoogleMaps.APIKey)
	
	if req.Query != "" {
		params.Add("query", req.Query)
	}
	if req.Location != "" {
		params.Add("location", req.Location)
	}
	if req.Radius > 0 {
		params.Add("radius", strconv.Itoa(req.Radius))
	}
	if req.Type != "" {
		params.Add("type", req.Type)
	}
	if req.Language != "" {
		params.Add("language", req.Language)
	}
	if req.Region != "" {
		params.Add("region", req.Region)
	}

	requestURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	// 發送請求
	resp, err := h.client.Get(requestURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, vo.ErrorResponse{
			Code:    "GOOGLE_MAPS_ERROR",
			Message: "Google Maps API 請求失敗",
			Error:   err.Error(),
		})
		return
	}
	defer resp.Body.Close()

	// 讀取回應
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, vo.ErrorResponse{
			Code:    "READ_RESPONSE_ERROR",
			Message: "讀取回應失敗",
			Error:   err.Error(),
		})
		return
	}

	// 解析回應
	var placesResp dto.PlacesSearchResponse
	if err := json.Unmarshal(body, &placesResp); err != nil {
		c.JSON(http.StatusInternalServerError, vo.ErrorResponse{
			Code:    "PARSE_RESPONSE_ERROR",
			Message: "解析回應失敗",
			Error:   err.Error(),
		})
		return
	}

	// 檢查 Google Maps API 狀態
	if placesResp.Status != "OK" && placesResp.Status != "ZERO_RESULTS" {
		var errorResp dto.GoogleMapsError
		json.Unmarshal(body, &errorResp)
		c.JSON(http.StatusBadRequest, vo.ErrorResponse{
			Code:    "GOOGLE_MAPS_API_ERROR",
			Message: errorResp.ErrorMessage,
			Error:   placesResp.Status,
		})
		return
	}

	c.JSON(http.StatusOK, placesResp)
}

// GetDirections 路線規劃
// @Summary 路線規劃
// @Description 計算從起點到終點的最佳路線
// @Tags google-maps
// @Accept json
// @Produce json
// @Param request body dto.DirectionsRequest true "路線規劃請求"
// @Success 200 {object} dto.DirectionsResponse
// @Failure 400 {object} vo.ErrorResponse
// @Failure 500 {object} vo.ErrorResponse
// @Router /google-maps/directions [post]
func (h *GoogleMapsHandler) GetDirections(c *gin.Context) {
	var req dto.DirectionsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, vo.ErrorResponse{
			Code:    "INVALID_REQUEST",
			Message: "無效的請求格式",
			Error:   err.Error(),
		})
		return
	}

	// 驗證請求
	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, vo.ErrorResponse{
			Code:    "VALIDATION_ERROR",
			Message: "請求驗證失敗",
			Error:   err.Error(),
		})
		return
	}

	// 檢查 API Key
	if h.config.GoogleMaps.APIKey == "" {
		c.JSON(http.StatusInternalServerError, vo.ErrorResponse{
			Code:    "MISSING_API_KEY",
			Message: "Google Maps API Key 未設定",
		})
		return
	}

	// 構建請求 URL
	baseURL := h.config.GoogleMaps.DirectionsURL
	params := url.Values{}
	params.Add("origin", req.Origin)
	params.Add("destination", req.Destination)
	params.Add("key", h.config.GoogleMaps.APIKey)
	
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

	requestURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	// 發送請求
	resp, err := h.client.Get(requestURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, vo.ErrorResponse{
			Code:    "GOOGLE_MAPS_ERROR",
			Message: "Google Maps API 請求失敗",
			Error:   err.Error(),
		})
		return
	}
	defer resp.Body.Close()

	// 讀取回應
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, vo.ErrorResponse{
			Code:    "READ_RESPONSE_ERROR",
			Message: "讀取回應失敗",
			Error:   err.Error(),
		})
		return
	}

	// 解析回應
	var directionsResp dto.DirectionsResponse
	if err := json.Unmarshal(body, &directionsResp); err != nil {
		c.JSON(http.StatusInternalServerError, vo.ErrorResponse{
			Code:    "PARSE_RESPONSE_ERROR",
			Message: "解析回應失敗",
			Error:   err.Error(),
		})
		return
	}

	// 檢查 Google Maps API 狀態
	if directionsResp.Status != "OK" && directionsResp.Status != "ZERO_RESULTS" {
		var errorResp dto.GoogleMapsError
		json.Unmarshal(body, &errorResp)
		c.JSON(http.StatusBadRequest, vo.ErrorResponse{
			Code:    "GOOGLE_MAPS_API_ERROR",
			Message: errorResp.ErrorMessage,
			Error:   directionsResp.Status,
		})
		return
	}

	c.JSON(http.StatusOK, directionsResp)
}

// GetDistanceMatrix 距離矩陣
// @Summary 距離矩陣
// @Description 計算多點間的距離和時間
// @Tags google-maps
// @Accept json
// @Produce json
// @Param request body dto.DistanceMatrixRequest true "距離矩陣請求"
// @Success 200 {object} dto.DistanceMatrixResponse
// @Failure 400 {object} vo.ErrorResponse
// @Failure 500 {object} vo.ErrorResponse
// @Router /google-maps/distance-matrix [post]
func (h *GoogleMapsHandler) GetDistanceMatrix(c *gin.Context) {
	var req dto.DistanceMatrixRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, vo.ErrorResponse{
			Code:    "INVALID_REQUEST",
			Message: "無效的請求格式",
			Error:   err.Error(),
		})
		return
	}

	// 驗證請求
	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, vo.ErrorResponse{
			Code:    "VALIDATION_ERROR",
			Message: "請求驗證失敗",
			Error:   err.Error(),
		})
		return
	}

	// 檢查 API Key
	if h.config.GoogleMaps.APIKey == "" {
		c.JSON(http.StatusInternalServerError, vo.ErrorResponse{
			Code:    "MISSING_API_KEY",
			Message: "Google Maps API Key 未設定",
		})
		return
	}

	// 構建請求 URL
	baseURL := h.config.GoogleMaps.DistanceMatrixURL
	params := url.Values{}
	params.Add("origins", strings.Join(req.Origins, "|"))
	params.Add("destinations", strings.Join(req.Destinations, "|"))
	params.Add("key", h.config.GoogleMaps.APIKey)
	
	if req.Mode != "" {
		params.Add("mode", req.Mode)
	}
	if req.Language != "" {
		params.Add("language", req.Language)
	}
	if req.Region != "" {
		params.Add("region", req.Region)
	}
	if req.Units != "" {
		params.Add("units", req.Units)
	}
	if req.TrafficModel != "" {
		params.Add("traffic_model", req.TrafficModel)
	}
	if req.DepartureTime != "" {
		params.Add("departure_time", req.DepartureTime)
	}
	if req.ArrivalTime != "" {
		params.Add("arrival_time", req.ArrivalTime)
	}

	requestURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	// 發送請求
	resp, err := h.client.Get(requestURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, vo.ErrorResponse{
			Code:    "GOOGLE_MAPS_ERROR",
			Message: "Google Maps API 請求失敗",
			Error:   err.Error(),
		})
		return
	}
	defer resp.Body.Close()

	// 讀取回應
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, vo.ErrorResponse{
			Code:    "READ_RESPONSE_ERROR",
			Message: "讀取回應失敗",
			Error:   err.Error(),
		})
		return
	}

	// 解析回應
	var distanceResp dto.DistanceMatrixResponse
	if err := json.Unmarshal(body, &distanceResp); err != nil {
		c.JSON(http.StatusInternalServerError, vo.ErrorResponse{
			Code:    "PARSE_RESPONSE_ERROR",
			Message: "解析回應失敗",
			Error:   err.Error(),
		})
		return
	}

	// 檢查 Google Maps API 狀態
	if distanceResp.Status != "OK" && distanceResp.Status != "ZERO_RESULTS" {
		var errorResp dto.GoogleMapsError
		json.Unmarshal(body, &errorResp)
		c.JSON(http.StatusBadRequest, vo.ErrorResponse{
			Code:    "GOOGLE_MAPS_API_ERROR",
			Message: errorResp.ErrorMessage,
			Error:   distanceResp.Status,
		})
		return
	}

	c.JSON(http.StatusOK, distanceResp)
}

// GetNearbyMentalHealthServices 搜尋附近心靈健康服務
// @Summary 搜尋附近心靈健康服務
// @Description 搜尋附近的心靈健康相關服務，包括醫院、診所、諮商中心等
// @Tags google-maps
// @Accept json
// @Produce json
// @Param latitude query number true "緯度"
// @Param longitude query number true "經度"
// @Param radius query int false "搜尋半徑（公尺）" default(5000)
// @Param type query string false "服務類型" Enums(hospital,health,establishment)
// @Success 200 {object} dto.PlacesSearchResponse
// @Failure 400 {object} vo.ErrorResponse
// @Failure 500 {object} vo.ErrorResponse
// @Router /google-maps/nearby-mental-health [get]
func (h *GoogleMapsHandler) GetNearbyMentalHealthServices(c *gin.Context) {
	// 解析查詢參數
	latitudeStr := c.Query("latitude")
	longitudeStr := c.Query("longitude")
	radiusStr := c.DefaultQuery("radius", "5000")
	serviceType := c.DefaultQuery("type", "health")

	latitude, err := strconv.ParseFloat(latitudeStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, vo.ErrorResponse{
			Code:    "INVALID_LATITUDE",
			Message: "無效的緯度值",
			Error:   err.Error(),
		})
		return
	}

	longitude, err := strconv.ParseFloat(longitudeStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, vo.ErrorResponse{
			Code:    "INVALID_LONGITUDE",
			Message: "無效的經度值",
			Error:   err.Error(),
		})
		return
	}

	radius, err := strconv.Atoi(radiusStr)
	if err != nil || radius <= 0 {
		radius = 5000
	}

	// 構建搜尋請求
	searchReq := dto.PlacesSearchRequest{
		Query:    "心理諮商 精神科 心理健康 諮商中心",
		Location: fmt.Sprintf("%.8f,%.8f", latitude, longitude),
		Radius:   radius,
		Type:     serviceType,
		Language: "zh-TW",
		Region:   "tw",
	}

	// 使用現有的搜尋地點方法
	searchReqJSON, _ := json.Marshal(searchReq)
	
	// 創建新的請求上下文來模擬內部調用
	mockCtx := &gin.Context{}
	mockCtx.Set("request_body", searchReqJSON)
	
	// 直接調用搜尋方法（這裡需要重構以支援內部調用）
	c.JSON(http.StatusOK, gin.H{
		"message": "此功能需要整合現有的搜尋地點方法",
		"request": searchReq,
	})
}
