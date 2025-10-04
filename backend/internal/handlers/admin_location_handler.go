package handlers

import (
	"net/http"
	"time"

	"mindhelp-backend/internal/database"
	"mindhelp-backend/internal/models"
	"mindhelp-backend/internal/vo"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AdminLocationHandler 管理員位置處理器
type AdminLocationHandler struct{}

// NewAdminLocationHandler 創建管理員位置處理器
func NewAdminLocationHandler() *AdminLocationHandler {
	return &AdminLocationHandler{}
}

// SeedLocations 添加測試位置數據
// @Summary 添加測試位置數據
// @Description 為測試目的添加一些示例位置數據
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} vo.Response
// @Failure 500 {object} vo.ErrorResponse
// @Router /admin/locations/seed [post]
func (h *AdminLocationHandler) SeedLocations(c *gin.Context) {
	db, err := database.GetDBSafely()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, vo.NewErrorResponse(
			"database_unavailable",
			"資料庫暫時無法使用，請稍後再試",
			"DATABASE_UNAVAILABLE",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	// 創建一個測試用戶 ID（如果不存在）
	testUserID := uuid.New()

	// 示例位置數據
	locations := []models.Location{
		{
			ID:          uuid.New(),
			UserID:      testUserID,
			Name:        "台北市立聯合醫院松德院區",
			Description: "專業心理健康服務，提供諮商、心理治療等服務",
			Address:     "台北市信義區松德路309號",
			Latitude:    25.0330,
			Longitude:   121.5654,
			Category:    "醫療機構",
			Phone:       "02-2726-3141",
			Website:     "https://tpech.gov.taipei/",
			Rating:      4.5,
			IsPublic:    true,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          uuid.New(),
			UserID:      testUserID,
			Name:        "台大醫院精神醫學部",
			Description: "台大醫院精神醫學部，提供專業心理健康諮詢",
			Address:     "台北市中正區中山南路7號",
			Latitude:    25.0408,
			Longitude:   121.5179,
			Category:    "醫療機構",
			Phone:       "02-2312-3456",
			Website:     "https://www.ntuh.gov.tw/",
			Rating:      4.7,
			IsPublic:    true,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          uuid.New(),
			UserID:      testUserID,
			Name:        "台北市立聯合醫院中興院區",
			Description: "提供心理健康服務和諮商輔導",
			Address:     "台北市大同區鄭州路145號",
			Latitude:    25.0661,
			Longitude:   121.5145,
			Category:    "醫療機構",
			Phone:       "02-2552-3234",
			Website:     "https://tpech.gov.taipei/",
			Rating:      4.3,
			IsPublic:    true,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          uuid.New(),
			UserID:      testUserID,
			Name:        "台北市立聯合醫院仁愛院區",
			Description: "綜合醫院，包含精神科門診服務",
			Address:     "台北市大安區仁愛路四段10號",
			Latitude:    25.0376,
			Longitude:   121.5440,
			Category:    "醫療機構",
			Phone:       "02-2709-3600",
			Website:     "https://tpech.gov.taipei/",
			Rating:      4.4,
			IsPublic:    true,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          uuid.New(),
			UserID:      testUserID,
			Name:        "台北市立聯合醫院忠孝院區",
			Description: "提供心理健康諮商和治療服務",
			Address:     "台北市南港區同德路87號",
			Latitude:    25.0562,
			Longitude:   121.6076,
			Category:    "醫療機構",
			Phone:       "02-2786-1288",
			Website:     "https://tpech.gov.taipei/",
			Rating:      4.2,
			IsPublic:    true,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	// 插入數據
	var createdCount int
	for _, location := range locations {
		if err := db.Create(&location).Error; err != nil {
			// 如果記錄已存在（根據名稱），跳過
			continue
		}
		createdCount++
	}

	c.JSON(http.StatusOK, vo.SuccessResponse(map[string]interface{}{
		"created_count":   createdCount,
		"total_locations": len(locations),
		"message":         "測試位置數據已添加",
	}, "Location data seeded successfully"))
}
