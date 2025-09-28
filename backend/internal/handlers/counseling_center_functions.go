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

// GetCounselingCenters 獲取諮商所列表
func GetCounselingCenters(c *gin.Context) {
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
			err.Error(),
			c.Request.URL.Path,
		))
		return
	}

	// 計算總數
	var total int64
	if err := db.Model(&models.CounselingCenter{}).Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, vo.NewErrorResponse(
			"internal_error",
			"Failed to get counseling centers count",
			"INTERNAL_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	// 獲取分頁資料
	var centers []models.CounselingCenter
	if err := db.Offset(offset).Limit(pageSize).Find(&centers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, vo.NewErrorResponse(
			"internal_error",
			"Failed to get counseling centers",
			"INTERNAL_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	c.JSON(http.StatusOK, vo.NewResponse("success", gin.H{
		"counseling_centers": centers,
		"total":              total,
		"page":               page,
		"page_size":          pageSize,
		"total_pages":        (total + int64(pageSize) - 1) / int64(pageSize),
	}))
}

// GetCounselingCenter 獲取單個諮商所
func GetCounselingCenter(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, vo.NewErrorResponse(
			"invalid_id",
			"Invalid counseling center ID format",
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
			err.Error(),
			c.Request.URL.Path,
		))
		return
	}

	var center models.CounselingCenter
	if err := db.Where("id = ?", id).First(&center).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, vo.NewErrorResponse(
				"not_found",
				"Counseling center not found",
				"NOT_FOUND",
				nil,
				c.Request.URL.Path,
			))
			return
		}
		c.JSON(http.StatusInternalServerError, vo.NewErrorResponse(
			"internal_error",
			"Failed to get counseling center",
			"INTERNAL_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	c.JSON(http.StatusOK, vo.NewResponse("success", center))
}

// CreateCounselingCenter 創建諮商所
func CreateCounselingCenter(c *gin.Context) {
	var center models.CounselingCenter
	if err := c.ShouldBindJSON(&center); err != nil {
		c.JSON(http.StatusBadRequest, vo.NewErrorResponse(
			"invalid_input",
			"Invalid request format",
			"VALIDATION_ERROR",
			err.Error(),
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
			err.Error(),
			c.Request.URL.Path,
		))
		return
	}

	if err := db.Create(&center).Error; err != nil {
		c.JSON(http.StatusInternalServerError, vo.NewErrorResponse(
			"internal_error",
			"Failed to create counseling center",
			"INTERNAL_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	c.JSON(http.StatusCreated, vo.NewResponse("success", center))
}

// UpdateCounselingCenter 更新諮商所
func UpdateCounselingCenter(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, vo.NewErrorResponse(
			"invalid_id",
			"Invalid counseling center ID format",
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
			err.Error(),
			c.Request.URL.Path,
		))
		return
	}

	var center models.CounselingCenter
	if err := db.Where("id = ?", id).First(&center).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, vo.NewErrorResponse(
				"not_found",
				"Counseling center not found",
				"NOT_FOUND",
				nil,
				c.Request.URL.Path,
			))
			return
		}
		c.JSON(http.StatusInternalServerError, vo.NewErrorResponse(
			"internal_error",
			"Failed to get counseling center",
			"INTERNAL_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	if err := c.ShouldBindJSON(&center); err != nil {
		c.JSON(http.StatusBadRequest, vo.NewErrorResponse(
			"invalid_input",
			"Invalid request format",
			"VALIDATION_ERROR",
			err.Error(),
			c.Request.URL.Path,
		))
		return
	}

	if err := db.Save(&center).Error; err != nil {
		c.JSON(http.StatusInternalServerError, vo.NewErrorResponse(
			"internal_error",
			"Failed to update counseling center",
			"INTERNAL_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	c.JSON(http.StatusOK, vo.NewResponse("success", center))
}

// DeleteCounselingCenter 刪除諮商所
func DeleteCounselingCenter(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, vo.NewErrorResponse(
			"invalid_id",
			"Invalid counseling center ID format",
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
			err.Error(),
			c.Request.URL.Path,
		))
		return
	}

	var center models.CounselingCenter
	if err := db.Where("id = ?", id).First(&center).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, vo.NewErrorResponse(
				"not_found",
				"Counseling center not found",
				"NOT_FOUND",
				nil,
				c.Request.URL.Path,
			))
			return
		}
		c.JSON(http.StatusInternalServerError, vo.NewErrorResponse(
			"internal_error",
			"Failed to get counseling center",
			"INTERNAL_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	if err := db.Delete(&center).Error; err != nil {
		c.JSON(http.StatusInternalServerError, vo.NewErrorResponse(
			"internal_error",
			"Failed to delete counseling center",
			"INTERNAL_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	c.JSON(http.StatusOK, vo.NewResponse("success", gin.H{
		"message": "Counseling center deleted successfully",
	}))
}
