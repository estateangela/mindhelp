package handlers

import (
	"encoding/json"
	"net/http"

	"mindhelp-backend/internal/database"
	"mindhelp-backend/internal/dto"
	"mindhelp-backend/internal/models"
	"mindhelp-backend/internal/vo"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ConfigHandler 配置處理器
type ConfigHandler struct {
	db *gorm.DB
}

// NewConfigHandler 創建新的配置處理器
func NewConfigHandler() *ConfigHandler {
	return &ConfigHandler{
		db: database.GetDB(),
	}
}

// GetConfig 獲取應用程式配置
// @Summary 獲取應用配置
// @Description 獲取 APP 的動態配置，包含功能開關和篩選選項
// @Tags config
// @Accept json
// @Produce json
// @Success 200 {object} vo.Response{data=dto.AppConfigResponse}
// @Failure 500 {object} vo.ErrorResponse
// @Router /config [get]
func (h *ConfigHandler) GetConfig(c *gin.Context) {
	// 獲取所有啟用的配置
	var configs []models.AppConfig
	if err := h.db.Where("is_active = ?", true).Find(&configs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, vo.NewErrorResponse(
			"internal_error",
			"Failed to get app config",
			"INTERNAL_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	// 解析配置
	response := dto.AppConfigResponse{
		Features: dto.FeatureConfig{
			EnableReviews:           true,  // 預設值
			EnableTherapistProfiles: false,
			EnableGroupChat:         false,
			EnableVideoConsult:      false,
			EnableSharing:          true,
		},
		Filters: dto.FilterConfig{
			ResourceTypes:  []dto.FilterOption{},
			Specialties:    []dto.FilterOption{},
			QuizCategories: []dto.FilterOption{},
		},
		SupportInfo: dto.SupportInfo{
			Email:        "support@mindhelp.com",
			Phone:        "+886-2-1234-5678",
			Website:      "https://mindhelp.com",
			WorkingHours: "週一至週五 09:00-18:00",
		},
	}

	// 處理各種配置
	for _, config := range configs {
		switch config.Key {
		case "features":
			var features dto.FeatureConfig
			if err := json.Unmarshal([]byte(config.Value), &features); err == nil {
				response.Features = features
			}
		
		case "resource_types":
			var resourceTypes []dto.FilterOption
			if err := json.Unmarshal([]byte(config.Value), &resourceTypes); err == nil {
				response.Filters.ResourceTypes = resourceTypes
			}
		
		case "specialties":
			var specialties []dto.FilterOption
			if err := json.Unmarshal([]byte(config.Value), &specialties); err == nil {
				response.Filters.Specialties = specialties
			}
		
		case "quiz_categories":
			var quizCategories []dto.FilterOption
			if err := json.Unmarshal([]byte(config.Value), &quizCategories); err == nil {
				response.Filters.QuizCategories = quizCategories
			}
		
		case "support_email":
			response.SupportInfo.Email = config.Value
		
		case "support_phone":
			response.SupportInfo.Phone = config.Value
		
		case "support_website":
			response.SupportInfo.Website = config.Value
		
		case "working_hours":
			response.SupportInfo.WorkingHours = config.Value
		}
	}

	// 如果沒有配置測驗類別，從資料庫中獲取
	if len(response.Filters.QuizCategories) == 0 {
		var categories []string
		h.db.Model(&models.Quiz{}).Where("is_active = ?", true).Distinct("category").Pluck("category", &categories)
		
		for _, category := range categories {
			displayName := h.getCategoryDisplayName(category)
			response.Filters.QuizCategories = append(response.Filters.QuizCategories, dto.FilterOption{
				Key:         category,
				DisplayName: displayName,
			})
		}
	}

	c.JSON(http.StatusOK, vo.SuccessResponse(response, "App config retrieved successfully"))
}

// getCategoryDisplayName 獲取類別的顯示名稱
func (h *ConfigHandler) getCategoryDisplayName(category string) string {
	categoryMap := map[string]string{
		"anxiety":    "焦慮症",
		"depression": "憂鬱症",
		"stress":     "壓力管理",
		"adhd":       "注意力不足過動症",
		"trauma":     "創傷治療",
		"sleep":      "睡眠障礙",
		"eating":     "飲食失調",
		"addiction":  "成癮問題",
	}
	
	if displayName, exists := categoryMap[category]; exists {
		return displayName
	}
	return category
}
