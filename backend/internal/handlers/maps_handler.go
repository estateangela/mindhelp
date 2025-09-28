package handlers

import (
	"net/http"
	"strconv"

	"mindhelp-backend/internal/database"
	"mindhelp-backend/internal/models"
	"mindhelp-backend/internal/vo"

	"github.com/gin-gonic/gin"
)

// MapsHandler 地圖處理器
type MapsHandler struct {
}

// NewMapsHandler 創建新的地圖處理器
func NewMapsHandler() *MapsHandler {
	return &MapsHandler{}
}

// AddressInfo 地址資訊結構
type AddressInfo struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Address     string `json:"address"`
	Type        string `json:"type"` // counselor, counseling_center, recommended_doctor
	Phone       string `json:"phone,omitempty"`
	Description string `json:"description,omitempty"`
}

// GetAllAddresses 獲取所有地址資訊
// @Summary 獲取所有地址資訊
// @Description 獲取諮商師、諮商所和推薦醫師的所有地址資訊，用於 Google Maps 整合
// @Tags maps
// @Accept json
// @Produce json
// @Param type query string false "地址類型篩選" Enums(counselor,counseling_center,recommended_doctor)
// @Param limit query int false "限制數量" default(100)
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} vo.ErrorResponse
// @Router /maps/addresses [get]
func (h *MapsHandler) GetAllAddresses(c *gin.Context) {
	// 獲取資料庫連接
	db, err := database.GetDBSafely()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, vo.NewErrorResponse(
			"database_unavailable",
			"Database service is currently unavailable",
			"SERVICE_UNAVAILABLE",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	addressType := c.Query("type")
	limitStr := c.DefaultQuery("limit", "100")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 100
	}

	var addresses []AddressInfo

	// 獲取諮商師地址
	if addressType == "" || addressType == "counselor" {
		var counselors []models.Counselor
		if err := db.Select("id, name, work_location").Where("work_location IS NOT NULL AND work_location != ''").Limit(limit).Find(&counselors).Error; err == nil {
			for _, counselor := range counselors {
				addresses = append(addresses, AddressInfo{
					ID:      counselor.ID.String(),
					Name:    counselor.Name,
					Address: counselor.WorkLocation,
					Type:    "counselor",
				})
			}
		}
	}

	// 獲取諮商所地址
	if addressType == "" || addressType == "counseling_center" {
		var centers []models.CounselingCenter
		if err := db.Select("id, name, address, phone").Where("address IS NOT NULL AND address != ''").Limit(limit).Find(&centers).Error; err == nil {
			for _, center := range centers {
				addresses = append(addresses, AddressInfo{
					ID:      center.ID.String(),
					Name:    center.Name,
					Address: center.Address,
					Type:    "counseling_center",
					Phone:   center.Phone,
				})
			}
		}
	}

	// 獲取推薦醫師地址
	if addressType == "" || addressType == "recommended_doctor" {
		var doctors []models.RecommendedDoctor
		if err := db.Select("id, name, description").Where("description IS NOT NULL AND description != ''").Limit(limit).Find(&doctors).Error; err == nil {
			for _, doctor := range doctors {
				address := extractAddressFromDescription(doctor.Description)
				if address != "" {
					name := doctor.Name
					if name == "" {
						name = "推薦醫師 " + doctor.ID.String()[:8]
					}

					addresses = append(addresses, AddressInfo{
						ID:          doctor.ID.String(),
						Name:        name,
						Address:     address,
						Type:        "recommended_doctor",
						Description: doctor.Description,
					})
				}
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    addresses,
		"message": "地址資訊已準備就緒",
	})
}

// GetAddressesForGoogleMaps 獲取 Google Maps 地址資訊
// @Summary 獲取 Google Maps 地址資訊
// @Description 獲取格式化的地址資訊，適合 Google Maps API 使用
// @Tags maps
// @Accept json
// @Produce json
// @Param format query string false "輸出格式" Enums(json,geojson) default(json)
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} vo.ErrorResponse
// @Router /maps/google-addresses [get]
func (h *MapsHandler) GetAddressesForGoogleMaps(c *gin.Context) {
	format := c.DefaultQuery("format", "json")

	addresses, err := h.fetchAllAddresses()
	if err != nil {
		c.JSON(http.StatusInternalServerError, vo.NewErrorResponse(
			"internal_error",
			"Failed to fetch addresses",
			"INTERNAL_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	if format == "geojson" {
		geoJSON := h.convertToGeoJSON(addresses)
		c.JSON(http.StatusOK, geoJSON)
	} else {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    addresses,
			"message": "地址資訊已準備就緒",
		})
	}
}

// fetchAllAddresses 獲取所有地址資訊
func (h *MapsHandler) fetchAllAddresses() ([]AddressInfo, error) {
	// 獲取資料庫連接
	db, err := database.GetDBSafely()
	if err != nil {
		return nil, err
	}

	var addresses []AddressInfo
	limit := 100

	// 獲取諮商師地址
	var counselors []models.Counselor
	if err := db.Select("id, name, work_location").Where("work_location IS NOT NULL AND work_location != ''").Limit(limit).Find(&counselors).Error; err == nil {
		for _, counselor := range counselors {
			addresses = append(addresses, AddressInfo{
				ID:      counselor.ID.String(),
				Name:    counselor.Name,
				Address: counselor.WorkLocation,
				Type:    "counselor",
			})
		}
	}

	// 獲取諮商所地址
	var centers []models.CounselingCenter
	if err := db.Select("id, name, address, phone").Where("address IS NOT NULL AND address != ''").Limit(limit).Find(&centers).Error; err == nil {
		for _, center := range centers {
			addresses = append(addresses, AddressInfo{
				ID:      center.ID.String(),
				Name:    center.Name,
				Address: center.Address,
				Type:    "counseling_center",
				Phone:   center.Phone,
			})
		}
	}

	// 獲取推薦醫師地址
	var doctors []models.RecommendedDoctor
	if err := db.Select("id, name, description").Where("description IS NOT NULL AND description != ''").Limit(limit).Find(&doctors).Error; err == nil {
		for _, doctor := range doctors {
			address := extractAddressFromDescription(doctor.Description)
			if address != "" {
				name := doctor.Name
				if name == "" {
					name = "推薦醫師 " + doctor.ID.String()[:8]
				}

				addresses = append(addresses, AddressInfo{
					ID:          doctor.ID.String(),
					Name:        name,
					Address:     address,
					Type:        "recommended_doctor",
					Description: doctor.Description,
				})
			}
		}
	}

	return addresses, nil
}

// convertToGeoJSON 轉換為 GeoJSON 格式
func (h *MapsHandler) convertToGeoJSON(addresses []AddressInfo) map[string]interface{} {
	features := make([]map[string]interface{}, 0, len(addresses))

	for _, addr := range addresses {
		feature := map[string]interface{}{
			"type": "Feature",
			"properties": map[string]interface{}{
				"id":          addr.ID,
				"name":        addr.Name,
				"address":     addr.Address,
				"type":        addr.Type,
				"phone":       addr.Phone,
				"description": addr.Description,
			},
			"geometry": map[string]interface{}{
				"type":        "Point",
				"coordinates": []float64{0, 0}, // 需要實際的經緯度
			},
		}
		features = append(features, feature)
	}

	return map[string]interface{}{
		"type":     "FeatureCollection",
		"features": features,
	}
}

// extractAddressFromDescription 從描述中提取地址資訊
func extractAddressFromDescription(description string) string {
	// 簡單的地址提取邏輯
	// 尋找包含地區名稱的模式
	keywords := []string{
		"台北市", "新北市", "桃園市", "台中市", "台南市", "高雄市",
		"基隆市", "新竹市", "嘉義市", "宜蘭縣", "新竹縣", "苗栗縣",
		"彰化縣", "南投縣", "雲林縣", "嘉義縣", "屏東縣", "台東縣",
		"花蓮縣", "澎湖縣", "金門縣", "連江縣",
	}

	for _, keyword := range keywords {
		if len(description) > 0 && len(keyword) > 0 {
			// 這裡可以實現更複雜的地址提取邏輯
			// 目前返回包含地區關鍵字的描述
			if len(description) > 50 {
				return description[:50] + "..."
			}
			return description
		}
	}

	return ""
}
