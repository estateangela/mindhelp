package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"mindhelp-backend/internal/database"
	"mindhelp-backend/internal/models"
	"mindhelp-backend/internal/vo"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// GetRecommendedDoctors 獲取推薦醫師列表
func GetRecommendedDoctors(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	sortBy := c.DefaultQuery("sort_by", "experience_count")
	sortOrder := c.DefaultQuery("sort_order", "desc")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	// 獲取資料庫連接
	db, err := database.GetDBSafely()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, vo.NewErrorResponse(
			"database_unavailable",
			"資料庫暫時無法使用，請稍後再試",
			"DATABASE_UNAVAILABLE",
			[]string{err.Error()},
			c.Request.URL.Path,
		))
		return
	}

	// 計算總數 - 只統計有 name 的記錄
	var total int64
	if err := db.Model(&models.RecommendedDoctor{}).
		Where("name IS NOT NULL AND name != ''").
		Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, vo.NewErrorResponse(
			"internal_error",
			"Failed to get recommended doctors count",
			"INTERNAL_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	// 構建查詢
	query := db.Where("name IS NOT NULL AND name != ''")

	// 排序
	if sortBy == "experience_count" && (sortOrder == "asc" || sortOrder == "desc") {
		query = query.Order("experience_count " + strings.ToUpper(sortOrder))
	} else {
		query = query.Order("updated_at DESC")
	}

	// 獲取分頁資料
	var doctors []models.RecommendedDoctor
	if err := query.Offset(offset).Limit(pageSize).Find(&doctors).Error; err != nil {
		c.JSON(http.StatusInternalServerError, vo.NewErrorResponse(
			"internal_error",
			"Failed to get recommended doctors",
			"INTERNAL_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	c.JSON(http.StatusOK, vo.NewResponse("success", gin.H{
		"recommended_doctors": doctors,
		"total":               total,
		"page":                page,
		"page_size":           pageSize,
		"total_pages":         (total + int64(pageSize) - 1) / int64(pageSize),
	}))
}

// GetRecommendedDoctor 獲取單個推薦醫師
func GetRecommendedDoctor(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, vo.NewErrorResponse(
			"invalid_id",
			"Invalid recommended doctor ID format",
			"VALIDATION_ERROR",
			[]string{err.Error()},
			c.Request.URL.Path,
		))
		return
	}

	// 獲取資料庫連接
	db, err := database.GetDBSafely()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, vo.NewErrorResponse(
			"database_unavailable",
			"資料庫暫時無法使用，請稍後再試",
			"DATABASE_UNAVAILABLE",
			[]string{err.Error()},
			c.Request.URL.Path,
		))
		return
	}

	var doctor models.RecommendedDoctor
	if err := db.Where("id = ?", id).First(&doctor).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, vo.NewErrorResponse(
				"not_found",
				"Recommended doctor not found",
				"NOT_FOUND",
				nil,
				c.Request.URL.Path,
			))
			return
		}
		c.JSON(http.StatusInternalServerError, vo.NewErrorResponse(
			"internal_error",
			"Failed to get recommended doctor",
			"INTERNAL_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	c.JSON(http.StatusOK, vo.NewResponse("success", doctor))
}

// CreateRecommendedDoctor 創建推薦醫師
func CreateRecommendedDoctor(c *gin.Context) {
	var doctor models.RecommendedDoctor
	if err := c.ShouldBindJSON(&doctor); err != nil {
		c.JSON(http.StatusBadRequest, vo.NewErrorResponse(
			"invalid_input",
			"Invalid request format",
			"VALIDATION_ERROR",
			[]string{err.Error()},
			c.Request.URL.Path,
		))
		return
	}

	// 獲取資料庫連接
	db, err := database.GetDBSafely()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, vo.NewErrorResponse(
			"database_unavailable",
			"資料庫暫時無法使用，請稍後再試",
			"DATABASE_UNAVAILABLE",
			[]string{err.Error()},
			c.Request.URL.Path,
		))
		return
	}

	if err := db.Create(&doctor).Error; err != nil {
		c.JSON(http.StatusInternalServerError, vo.NewErrorResponse(
			"internal_error",
			"Failed to create recommended doctor",
			"INTERNAL_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	c.JSON(http.StatusCreated, vo.NewResponse("success", doctor))
}

// UpdateRecommendedDoctor 更新推薦醫師
func UpdateRecommendedDoctor(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, vo.NewErrorResponse(
			"invalid_id",
			"Invalid recommended doctor ID format",
			"VALIDATION_ERROR",
			[]string{err.Error()},
			c.Request.URL.Path,
		))
		return
	}

	// 獲取資料庫連接
	db, err := database.GetDBSafely()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, vo.NewErrorResponse(
			"database_unavailable",
			"資料庫暫時無法使用，請稍後再試",
			"DATABASE_UNAVAILABLE",
			[]string{err.Error()},
			c.Request.URL.Path,
		))
		return
	}

	var doctor models.RecommendedDoctor
	if err := db.Where("id = ?", id).First(&doctor).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, vo.NewErrorResponse(
				"not_found",
				"Recommended doctor not found",
				"NOT_FOUND",
				nil,
				c.Request.URL.Path,
			))
			return
		}
		c.JSON(http.StatusInternalServerError, vo.NewErrorResponse(
			"internal_error",
			"Failed to get recommended doctor",
			"INTERNAL_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	if err := c.ShouldBindJSON(&doctor); err != nil {
		c.JSON(http.StatusBadRequest, vo.NewErrorResponse(
			"invalid_input",
			"Invalid request format",
			"VALIDATION_ERROR",
			[]string{err.Error()},
			c.Request.URL.Path,
		))
		return
	}

	if err := db.Save(&doctor).Error; err != nil {
		c.JSON(http.StatusInternalServerError, vo.NewErrorResponse(
			"internal_error",
			"Failed to update recommended doctor",
			"INTERNAL_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	c.JSON(http.StatusOK, vo.NewResponse("success", doctor))
}

// DeleteRecommendedDoctor 刪除推薦醫師
func DeleteRecommendedDoctor(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, vo.NewErrorResponse(
			"invalid_id",
			"Invalid recommended doctor ID format",
			"VALIDATION_ERROR",
			[]string{err.Error()},
			c.Request.URL.Path,
		))
		return
	}

	// 獲取資料庫連接
	db, err := database.GetDBSafely()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, vo.NewErrorResponse(
			"database_unavailable",
			"資料庫暫時無法使用，請稍後再試",
			"DATABASE_UNAVAILABLE",
			[]string{err.Error()},
			c.Request.URL.Path,
		))
		return
	}

	var doctor models.RecommendedDoctor
	if err := db.Where("id = ?", id).First(&doctor).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, vo.NewErrorResponse(
				"not_found",
				"Recommended doctor not found",
				"NOT_FOUND",
				nil,
				c.Request.URL.Path,
			))
			return
		}
		c.JSON(http.StatusInternalServerError, vo.NewErrorResponse(
			"internal_error",
			"Failed to get recommended doctor",
			"INTERNAL_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	if err := db.Delete(&doctor).Error; err != nil {
		c.JSON(http.StatusInternalServerError, vo.NewErrorResponse(
			"internal_error",
			"Failed to delete recommended doctor",
			"INTERNAL_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	c.JSON(http.StatusOK, vo.NewResponse("success", gin.H{
		"message": "Recommended doctor deleted successfully",
	}))
}
