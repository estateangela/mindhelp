package handlers

import (
	"net/http"
	"strconv"

	"mindhelp-backend/internal/database"
	"mindhelp-backend/internal/dto"
	"mindhelp-backend/internal/models"
	"mindhelp-backend/internal/vo"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GetCounselingCenters 獲取諮商所列表
func GetCounselingCenters(c *gin.Context) {
	// 解析查詢參數
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	search := c.Query("search")
	onlineOnly := c.Query("online_only") == "true"

	// 確保分頁參數合理
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	// 構建查詢
	query := database.GetDB().Model(&models.CounselingCenter{})

	// 添加搜索條件
	if search != "" {
		query = query.Where("name ILIKE ? OR address ILIKE ?", "%"+search+"%", "%"+search+"%")
	}
	if onlineOnly {
		query = query.Where("online_counseling = ?", true)
	}

	// 獲取總數
	var total int64
	if err := query.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, vo.ErrorResponse{
			Code:    "INTERNAL_SERVER_ERROR",
			Message: "Failed to count counseling centers",
			Error:   err.Error(),
		})
		return
	}

	// 獲取諮商所列表
	var centers []models.CounselingCenter
	if err := query.Offset(offset).Limit(pageSize).Find(&centers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, vo.ErrorResponse{
			Code:    "INTERNAL_SERVER_ERROR",
			Message: "Failed to fetch counseling centers",
			Error:   err.Error(),
		})
		return
	}

	// 轉換為回應格式
	var centerResponses []dto.CounselingCenterResponse
	for _, center := range centers {
		centerResponses = append(centerResponses, dto.CounselingCenterResponse{
			ID:               center.ID.String(),
			Name:             center.Name,
			Address:          center.Address,
			Phone:            center.Phone,
			OnlineCounseling: center.OnlineCounseling,
			CreatedAt:        center.CreatedAt,
			UpdatedAt:        center.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, dto.CounselingCenterListResponse{
		CounselingCenters: centerResponses,
		Total:             total,
		Page:              page,
		PageSize:          pageSize,
	})
}

// GetCounselingCenter 獲取單個諮商所
func GetCounselingCenter(c *gin.Context) {
	id := c.Param("id")
	centerID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, vo.ErrorResponse{
			Code:    "BAD_REQUEST",
			Message: "Invalid counseling center ID",
			Error:   err.Error(),
		})
		return
	}

	var center models.CounselingCenter
	if err := database.GetDB().First(&center, "id = ?", centerID).Error; err != nil {
		c.JSON(http.StatusNotFound, vo.ErrorResponse{
			Code:    "NOT_FOUND",
			Message: "Counseling center not found",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.CounselingCenterResponse{
		ID:               center.ID.String(),
		Name:             center.Name,
		Address:          center.Address,
		Phone:            center.Phone,
		OnlineCounseling: center.OnlineCounseling,
		CreatedAt:        center.CreatedAt,
		UpdatedAt:        center.UpdatedAt,
	})
}

// CreateCounselingCenter 創建諮商所
func CreateCounselingCenter(c *gin.Context) {
	var req dto.CounselingCenterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, vo.ErrorResponse{
			Code:    "BAD_REQUEST",
			Message: "Invalid request data",
			Error:   err.Error(),
		})
		return
	}

	center := models.CounselingCenter{
		Name:             req.Name,
		Address:          req.Address,
		Phone:            req.Phone,
		OnlineCounseling: req.OnlineCounseling,
	}

	if err := database.GetDB().Create(&center).Error; err != nil {
		c.JSON(http.StatusInternalServerError, vo.ErrorResponse{
			Code:    "INTERNAL_SERVER_ERROR",
			Message: "Failed to create counseling center",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, dto.CounselingCenterResponse{
		ID:               center.ID.String(),
		Name:             center.Name,
		Address:          center.Address,
		Phone:            center.Phone,
		OnlineCounseling: center.OnlineCounseling,
		CreatedAt:        center.CreatedAt,
		UpdatedAt:        center.UpdatedAt,
	})
}

// UpdateCounselingCenter 更新諮商所
func UpdateCounselingCenter(c *gin.Context) {
	id := c.Param("id")
	centerID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, vo.ErrorResponse{
			Code:    "BAD_REQUEST",
			Message: "Invalid counseling center ID",
			Error:   err.Error(),
		})
		return
	}

	var req dto.CounselingCenterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, vo.ErrorResponse{
			Code:    "BAD_REQUEST",
			Message: "Invalid request data",
			Error:   err.Error(),
		})
		return
	}

	var center models.CounselingCenter
	if err := database.GetDB().First(&center, "id = ?", centerID).Error; err != nil {
		c.JSON(http.StatusNotFound, vo.ErrorResponse{
			Code:    "NOT_FOUND",
			Message: "Counseling center not found",
			Error:   err.Error(),
		})
		return
	}

	// 更新欄位
	center.Name = req.Name
	center.Address = req.Address
	center.Phone = req.Phone
	center.OnlineCounseling = req.OnlineCounseling

	if err := database.GetDB().Save(&center).Error; err != nil {
		c.JSON(http.StatusInternalServerError, vo.ErrorResponse{
			Code:    "INTERNAL_SERVER_ERROR",
			Message: "Failed to update counseling center",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.CounselingCenterResponse{
		ID:               center.ID.String(),
		Name:             center.Name,
		Address:          center.Address,
		Phone:            center.Phone,
		OnlineCounseling: center.OnlineCounseling,
		CreatedAt:        center.CreatedAt,
		UpdatedAt:        center.UpdatedAt,
	})
}

// DeleteCounselingCenter 刪除諮商所
func DeleteCounselingCenter(c *gin.Context) {
	id := c.Param("id")
	centerID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, vo.ErrorResponse{
			Code:    "BAD_REQUEST",
			Message: "Invalid counseling center ID",
			Error:   err.Error(),
		})
		return
	}

	if err := database.GetDB().Delete(&models.CounselingCenter{}, "id = ?", centerID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, vo.ErrorResponse{
			Code:    "INTERNAL_SERVER_ERROR",
			Message: "Failed to delete counseling center",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, vo.SuccessResponse(nil, "Counseling center deleted successfully"))
}
