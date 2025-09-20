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
// @Summary 獲取諮商所列表
// @Description 獲取諮商所列表，支援分頁和搜索
// @Tags counseling-centers
// @Accept json
// @Produce json
// @Param page query int false "頁碼" default(1)
// @Param page_size query int false "每頁數量" default(10)
// @Param search query string false "搜索關鍵字"
// @Param online_only query boolean false "僅顯示線上諮商"
// @Success 200 {object} dto.CounselingCenterListResponse
// @Failure 500 {object} vo.ErrorResponse
// @Router /counseling-centers [get]
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
// @Summary 獲取單個諮商所
// @Description 根據ID獲取諮商所詳細資訊
// @Tags counseling-centers
// @Accept json
// @Produce json
// @Param id path string true "諮商所ID"
// @Success 200 {object} dto.CounselingCenterResponse
// @Failure 400 {object} vo.ErrorResponse
// @Failure 404 {object} vo.ErrorResponse
// @Router /counseling-centers/{id} [get]
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
// @Summary 創建諮商所
// @Description 創建新的諮商所記錄
// @Tags counseling-centers
// @Accept json
// @Produce json
// @Param center body dto.CounselingCenterRequest true "諮商所資訊"
// @Success 201 {object} dto.CounselingCenterResponse
// @Failure 400 {object} vo.ErrorResponse
// @Failure 500 {object} vo.ErrorResponse
// @Router /admin/counseling-centers [post]
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
// @Summary 更新諮商所
// @Description 更新諮商所資訊
// @Tags counseling-centers
// @Accept json
// @Produce json
// @Param id path string true "諮商所ID"
// @Param center body dto.CounselingCenterRequest true "諮商所資訊"
// @Success 200 {object} dto.CounselingCenterResponse
// @Failure 400 {object} vo.ErrorResponse
// @Failure 404 {object} vo.ErrorResponse
// @Failure 500 {object} vo.ErrorResponse
// @Router /admin/counseling-centers/{id} [put]
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
// @Summary 刪除諮商所
// @Description 刪除諮商所記錄
// @Tags counseling-centers
// @Accept json
// @Produce json
// @Param id path string true "諮商所ID"
// @Success 200 {object} vo.Response
// @Failure 400 {object} vo.ErrorResponse
// @Failure 500 {object} vo.ErrorResponse
// @Router /admin/counseling-centers/{id} [delete]
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
