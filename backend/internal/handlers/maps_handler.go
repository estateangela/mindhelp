package handlers

import (
	"net/http"
	"strconv"

	"mindhelp-backend/internal/config"
	"mindhelp-backend/internal/database"
	"mindhelp-backend/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// MapsHandler 地圖處理器
type MapsHandler struct {
	cfg *config.Config
}

// NewMapsHandler 創建新的地圖處理器
func NewMapsHandler(cfg *config.Config) *MapsHandler {
	return &MapsHandler{
		cfg: cfg,
	}
}

// getDB 安全獲取資料庫連接
func (h *MapsHandler) getDB(c *gin.Context) (*gorm.DB, bool) {
	db, err := database.GetDBSafely()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "Database temporarily unavailable",
			"details": err.Error(),
		})
		return nil, false
	}
	return db, true
}

// GetAllAddresses 獲取所有地址資訊
// @Summary 獲取所有地址
// @Description 獲取諮商師、諮商所和推薦醫師的地址資訊
// @Tags maps
// @Accept json
// @Produce json
// @Param limit query int false "限制數量" default(100)
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /maps/addresses [get]
func (h *MapsHandler) GetAllAddresses(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "100")
	limit, _ := strconv.Atoi(limitStr)
	if limit <= 0 || limit > 1000 {
		limit = 100
	}

	db, ok := h.getDB(c)
	if !ok {
		return
	}

	// 獲取諮商師地址
	var counselors []models.Counselor
	if err := db.Select("id, name, work_location").
		Where("work_location IS NOT NULL AND work_location != ''").
		Limit(limit).
		Find(&counselors).Error; err != nil {
		// 記錄錯誤但繼續
	}

	// 獲取諮商所地址
	var centers []models.CounselingCenter
	if err := db.Select("id, name, address, phone").
		Where("address IS NOT NULL AND address != ''").
		Limit(limit).
		Find(&centers).Error; err != nil {
		// 記錄錯誤但繼續
	}

	// 獲取推薦醫師地址
	var doctors []models.RecommendedDoctor
	if err := db.Select("id, description").
		Where("description IS NOT NULL AND description != ''").
		Limit(limit).
		Find(&doctors).Error; err != nil {
		// 記錄錯誤但繼續
	}

	c.JSON(http.StatusOK, gin.H{
		"counselors":         counselors,
		"counseling_centers": centers,
		"recommended_doctors": doctors,
		"total": len(counselors) + len(centers) + len(doctors),
	})
}

// GetAddressesForGoogleMaps 為 Google Maps 獲取地址資訊
// @Summary 為 Google Maps 獲取地址
// @Description 獲取格式化的地址資訊供 Google Maps 使用
// @Tags maps
// @Accept json
// @Produce json
// @Param format query string false "輸出格式" Enums(json,geojson) default(json)
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /maps/google-addresses [get]
func (h *MapsHandler) GetAddressesForGoogleMaps(c *gin.Context) {
	format := c.DefaultQuery("format", "json")

	db, ok := h.getDB(c)
	if !ok {
		return
	}

	// 獲取諮商師工作地點
	var counselors []models.Counselor
	if err := db.Select("id, name, work_location").
		Where("work_location IS NOT NULL AND work_location != ''").
		Find(&counselors).Error; err != nil {
		// 記錄錯誤但繼續
	}

	// 獲取諮商所地址
	var centers []models.CounselingCenter
	if err := db.Select("id, name, address, phone").
		Where("address IS NOT NULL AND address != ''").
		Find(&centers).Error; err != nil {
		// 記錄錯誤但繼續
	}

	// 獲取推薦醫師地址
	var doctors []models.RecommendedDoctor
	if err := db.Select("id, description").
		Where("description IS NOT NULL AND description != ''").
		Find(&doctors).Error; err != nil {
		// 記錄錯誤但繼續
	}

	if format == "geojson" {
		// 返回 GeoJSON 格式
		features := make([]map[string]interface{}, 0)

		// 添加諮商師位置
		for _, counselor := range counselors {
			features = append(features, map[string]interface{}{
				"type": "Feature",
				"properties": map[string]interface{}{
					"id":   counselor.ID,
					"name": counselor.Name,
					"type": "counselor",
					"address": counselor.WorkLocation,
				},
				"geometry": map[string]interface{}{
					"type": "Point",
					"coordinates": []float64{0, 0}, // 需要實際座標
				},
			})
		}

		c.JSON(http.StatusOK, gin.H{
			"type": "FeatureCollection",
			"features": features,
		})
	} else {
		// 返回 JSON 格式
		var addresses []map[string]interface{}

		// 添加諮商師地址
		for _, counselor := range counselors {
			addresses = append(addresses, map[string]interface{}{
				"id":      counselor.ID,
				"name":    counselor.Name,
				"type":    "counselor",
				"address": counselor.WorkLocation,
			})
		}

		// 添加諮商所地址
		for _, center := range centers {
			addresses = append(addresses, map[string]interface{}{
				"id":      center.ID,
				"name":    center.Name,
				"type":    "counseling_center",
				"address": center.Address,
				"phone":   center.Phone,
			})
		}

		c.JSON(http.StatusOK, gin.H{
			"addresses": addresses,
			"total":     len(addresses),
		})
	}
}
