package handlers

import (
	"net/http"
	"strconv"

	"mindhelp-backend/internal/database"
	"mindhelp-backend/internal/models"
	"mindhelp-backend/internal/vo"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// GetCounselors 獲取諮商師列表
func GetCounselors(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

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

	// 計算總數
	var total int64
	if err := db.Model(&models.Counselor{}).Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, vo.NewErrorResponse(
			"internal_error",
			"Failed to get counselors count",
			"INTERNAL_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	// 獲取分頁資料
	var counselors []models.Counselor
	if err := db.Offset(offset).Limit(pageSize).Find(&counselors).Error; err != nil {
		c.JSON(http.StatusInternalServerError, vo.NewErrorResponse(
			"internal_error",
			"Failed to get counselors",
			"INTERNAL_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	c.JSON(http.StatusOK, vo.NewResponse("success", gin.H{
		"counselors":  counselors,
		"total":       total,
		"page":        page,
		"page_size":   pageSize,
		"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
	}))
}

// GetCounselor 獲取單個諮商師
func GetCounselor(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, vo.NewErrorResponse(
			"invalid_id",
			"Invalid counselor ID format",
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

	var counselor models.Counselor
	if err := db.Where("id = ?", id).First(&counselor).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, vo.NewErrorResponse(
				"not_found",
				"Counselor not found",
				"NOT_FOUND",
				nil,
				c.Request.URL.Path,
			))
			return
		}
		c.JSON(http.StatusInternalServerError, vo.NewErrorResponse(
			"internal_error",
			"Failed to get counselor",
			"INTERNAL_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	c.JSON(http.StatusOK, vo.NewResponse("success", counselor))
}

// CreateCounselor 創建諮商師
func CreateCounselor(c *gin.Context) {
	var counselor models.Counselor
	if err := c.ShouldBindJSON(&counselor); err != nil {
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

	if err := db.Create(&counselor).Error; err != nil {
		c.JSON(http.StatusInternalServerError, vo.NewErrorResponse(
			"internal_error",
			"Failed to create counselor",
			"INTERNAL_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	c.JSON(http.StatusCreated, vo.NewResponse("success", counselor))
}

// UpdateCounselor 更新諮商師
func UpdateCounselor(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, vo.NewErrorResponse(
			"invalid_id",
			"Invalid counselor ID format",
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

	var counselor models.Counselor
	if err := db.Where("id = ?", id).First(&counselor).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, vo.NewErrorResponse(
				"not_found",
				"Counselor not found",
				"NOT_FOUND",
				nil,
				c.Request.URL.Path,
			))
			return
		}
		c.JSON(http.StatusInternalServerError, vo.NewErrorResponse(
			"internal_error",
			"Failed to get counselor",
			"INTERNAL_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	if err := c.ShouldBindJSON(&counselor); err != nil {
		c.JSON(http.StatusBadRequest, vo.NewErrorResponse(
			"invalid_input",
			"Invalid request format",
			"VALIDATION_ERROR",
			[]string{err.Error()},
			c.Request.URL.Path,
		))
		return
	}

	if err := db.Save(&counselor).Error; err != nil {
		c.JSON(http.StatusInternalServerError, vo.NewErrorResponse(
			"internal_error",
			"Failed to update counselor",
			"INTERNAL_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	c.JSON(http.StatusOK, vo.NewResponse("success", counselor))
}

// DeleteCounselor 刪除諮商師
func DeleteCounselor(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, vo.NewErrorResponse(
			"invalid_id",
			"Invalid counselor ID format",
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

	var counselor models.Counselor
	if err := db.Where("id = ?", id).First(&counselor).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, vo.NewErrorResponse(
				"not_found",
				"Counselor not found",
				"NOT_FOUND",
				nil,
				c.Request.URL.Path,
			))
			return
		}
		c.JSON(http.StatusInternalServerError, vo.NewErrorResponse(
			"internal_error",
			"Failed to get counselor",
			"INTERNAL_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	if err := db.Delete(&counselor).Error; err != nil {
		c.JSON(http.StatusInternalServerError, vo.NewErrorResponse(
			"internal_error",
			"Failed to delete counselor",
			"INTERNAL_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	c.JSON(http.StatusOK, vo.NewResponse("success", gin.H{
		"message": "Counselor deleted successfully",
	}))
}
