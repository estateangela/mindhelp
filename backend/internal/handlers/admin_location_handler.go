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

// SeedLocations 添加台北地區心理健康服務機構數據
// @Summary 添加台北地區心理健康服務機構數據
// @Description 為測試目的添加台北地區所有已知的心理健康服務機構數據
// @Tags location
// @Accept json
// @Produce json
// @Success 200 {object} vo.Response
// @Failure 500 {object} vo.ErrorResponse
// @Router /locations/seed [post]
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

	// 台北地區心理健康服務機構數據
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
		{
			ID:          uuid.New(),
			UserID:      testUserID,
			Name:        "台北榮總精神醫學部",
			Description: "台北榮民總醫院精神醫學部，提供專業精神醫療服務",
			Address:     "台北市北投區石牌路二段201號",
			Latitude:    25.1189,
			Longitude:   121.5184,
			Category:    "醫療機構",
			Phone:       "02-2871-2121",
			Website:     "https://www.vghtpe.gov.tw/",
			Rating:      4.6,
			IsPublic:    true,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          uuid.New(),
			UserID:      testUserID,
			Name:        "三軍總醫院精神醫學部",
			Description: "三軍總醫院精神醫學部，提供軍民精神醫療服務",
			Address:     "台北市內湖區成功路二段325號",
			Latitude:    25.0682,
			Longitude:   121.5944,
			Category:    "醫療機構",
			Phone:       "02-8792-3311",
			Website:     "https://www.tsgh.ndmctsgh.edu.tw/",
			Rating:      4.4,
			IsPublic:    true,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          uuid.New(),
			UserID:      testUserID,
			Name:        "馬偕紀念醫院精神醫學部",
			Description: "馬偕紀念醫院精神醫學部，提供基督教精神醫療服務",
			Address:     "台北市中山區中山北路二段92號",
			Latitude:    25.0588,
			Longitude:   121.5219,
			Category:    "醫療機構",
			Phone:       "02-2543-3535",
			Website:     "https://www.mmh.org.tw/",
			Rating:      4.3,
			IsPublic:    true,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          uuid.New(),
			UserID:      testUserID,
			Name:        "新光醫院精神科",
			Description: "新光醫院精神科，提供專業精神醫療服務",
			Address:     "台北市士林區文昌路95號",
			Latitude:    25.0969,
			Longitude:   121.5225,
			Category:    "醫療機構",
			Phone:       "02-2833-2211",
			Website:     "https://www.skh.org.tw/",
			Rating:      4.1,
			IsPublic:    true,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          uuid.New(),
			UserID:      testUserID,
			Name:        "台北市立萬芳醫院精神科",
			Description: "台北市立萬芳醫院精神科，提供社區精神醫療服務",
			Address:     "台北市文山區興隆路三段111號",
			Latitude:    24.9894,
			Longitude:   121.5675,
			Category:    "醫療機構",
			Phone:       "02-2930-7930",
			Website:     "https://www.wanfang.gov.taipei/",
			Rating:      4.2,
			IsPublic:    true,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          uuid.New(),
			UserID:      testUserID,
			Name:        "台北市立關渡醫院精神科",
			Description: "台北市立關渡醫院精神科，提供精神復健服務",
			Address:     "台北市北投區知行路225巷12號",
			Latitude:    25.1214,
			Longitude:   121.4694,
			Category:    "醫療機構",
			Phone:       "02-2858-7000",
			Website:     "https://www.gandau.gov.taipei/",
			Rating:      4.0,
			IsPublic:    true,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          uuid.New(),
			UserID:      testUserID,
			Name:        "台北市立療養院",
			Description: "台北市立療養院，專業精神醫療機構",
			Address:     "台北市松山區三民路5號",
			Latitude:    25.0631,
			Longitude:   121.5597,
			Category:    "醫療機構",
			Phone:       "02-2726-3141",
			Website:     "https://www.tpech.gov.taipei/",
			Rating:      4.3,
			IsPublic:    true,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          uuid.New(),
			UserID:      testUserID,
			Name:        "台北市立聯合醫院陽明院區精神科",
			Description: "台北市立聯合醫院陽明院區精神科，提供精神醫療服務",
			Address:     "台北市士林區雨聲街105號",
			Latitude:    25.1047,
			Longitude:   121.5229,
			Category:    "醫療機構",
			Phone:       "02-2835-3456",
			Website:     "https://tpech.gov.taipei/",
			Rating:      4.1,
			IsPublic:    true,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          uuid.New(),
			UserID:      testUserID,
			Name:        "台北市立聯合醫院和平婦幼院區精神科",
			Description: "台北市立聯合醫院和平婦幼院區精神科，提供婦幼精神醫療服務",
			Address:     "台北市中正區中華路二段33號",
			Latitude:    25.0426,
			Longitude:   121.5083,
			Category:    "醫療機構",
			Phone:       "02-2388-9595",
			Website:     "https://tpech.gov.taipei/",
			Rating:      4.2,
			IsPublic:    true,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          uuid.New(),
			UserID:      testUserID,
			Name:        "台北市立聯合醫院林森中醫昆明院區精神科",
			Description: "台北市立聯合醫院林森中醫昆明院區精神科，提供中西醫整合精神醫療",
			Address:     "台北市萬華區昆明街100號",
			Latitude:    25.0408,
			Longitude:   121.5073,
			Category:    "醫療機構",
			Phone:       "02-2388-7088",
			Website:     "https://tpech.gov.taipei/",
			Rating:      4.0,
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
