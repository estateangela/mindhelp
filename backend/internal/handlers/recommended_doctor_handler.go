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

// GetRecommendedDoctors 獲取推薦醫師列表
// @Summary 獲取推薦醫師列表
// @Description 獲取推薦醫師列表，支援分頁、搜索和排序
// @Tags recommended-doctors
// @Accept json
// @Produce json
// @Param page query int false "頁碼" default(1)
// @Param page_size query int false "每頁數量" default(10)
// @Param search query string false "搜索關鍵字"
// @Param sort_by query string false "排序欄位" default(experience_count)
// @Param sort_order query string false "排序順序" default(desc)
// @Success 200 {object} dto.RecommendedDoctorListResponse
// @Failure 500 {object} vo.ErrorResponse
// @Router /recommended-doctors [get]
func GetRecommendedDoctors(c *gin.Context) {
	// 解析查詢參數
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	search := c.Query("search")
	sortBy := c.DefaultQuery("sort_by", "experience_count")
	sortOrder := c.DefaultQuery("sort_order", "desc")

	// 確保分頁參數合理
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	// 構建查詢
	query := database.GetDB().Model(&models.RecommendedDoctor{})

	// 添加搜索條件
	if search != "" {
		query = query.Where("name ILIKE ? OR description ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	// 添加排序
	orderClause := sortBy + " " + sortOrder
	if sortOrder != "asc" && sortOrder != "desc" {
		orderClause = "experience_count desc"
	}
	query = query.Order(orderClause)

	// 獲取總數
	var total int64
	if err := query.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, vo.ErrorResponse{
			Code:    "INTERNAL_SERVER_ERROR",
			Message: "Failed to count recommended doctors",
			Error:   err.Error(),
		})
		return
	}

	// 獲取推薦醫師列表
	var doctors []models.RecommendedDoctor
	if err := query.Offset(offset).Limit(pageSize).Find(&doctors).Error; err != nil {
		c.JSON(http.StatusInternalServerError, vo.ErrorResponse{
			Code:    "INTERNAL_SERVER_ERROR",
			Message: "Failed to fetch recommended doctors",
			Error:   err.Error(),
		})
		return
	}

	// 轉換為回應格式
	var doctorResponses []dto.RecommendedDoctorResponse
	for _, doctor := range doctors {
		doctorResponses = append(doctorResponses, dto.RecommendedDoctorResponse{
			ID:              doctor.ID.String(),
			Name:            doctor.Name,
			Description:     doctor.Description,
			ExperienceCount: doctor.ExperienceCount,
			CreatedAt:       doctor.CreatedAt,
			UpdatedAt:       doctor.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, dto.RecommendedDoctorListResponse{
		RecommendedDoctors: doctorResponses,
		Total:              total,
		Page:               page,
		PageSize:           pageSize,
	})
}

// GetRecommendedDoctor 獲取單個推薦醫師
// @Summary 獲取單個推薦醫師
// @Description 根據ID獲取推薦醫師詳細資訊
// @Tags recommended-doctors
// @Accept json
// @Produce json
// @Param id path string true "推薦醫師ID"
// @Success 200 {object} dto.RecommendedDoctorResponse
// @Failure 400 {object} vo.ErrorResponse
// @Failure 404 {object} vo.ErrorResponse
// @Router /recommended-doctors/{id} [get]
func GetRecommendedDoctor(c *gin.Context) {
	id := c.Param("id")
	doctorID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, vo.ErrorResponse{
			Code:    "BAD_REQUEST",
			Message: "Invalid recommended doctor ID",
			Error:   err.Error(),
		})
		return
	}

	var doctor models.RecommendedDoctor
	if err := database.GetDB().First(&doctor, "id = ?", doctorID).Error; err != nil {
		c.JSON(http.StatusNotFound, vo.ErrorResponse{
			Code:    "NOT_FOUND",
			Message: "Recommended doctor not found",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.RecommendedDoctorResponse{
		ID:              doctor.ID.String(),
		Name:            doctor.Name,
		Description:     doctor.Description,
		ExperienceCount: doctor.ExperienceCount,
		CreatedAt:       doctor.CreatedAt,
		UpdatedAt:       doctor.UpdatedAt,
	})
}

// CreateRecommendedDoctor 創建推薦醫師
// @Summary 創建推薦醫師
// @Description 創建新的推薦醫師記錄
// @Tags recommended-doctors
// @Accept json
// @Produce json
// @Param doctor body dto.RecommendedDoctorRequest true "推薦醫師資訊"
// @Success 201 {object} dto.RecommendedDoctorResponse
// @Failure 400 {object} vo.ErrorResponse
// @Failure 500 {object} vo.ErrorResponse
// @Router /admin/recommended-doctors [post]
func CreateRecommendedDoctor(c *gin.Context) {
	var req dto.RecommendedDoctorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, vo.ErrorResponse{
			Code:    "BAD_REQUEST",
			Message: "Invalid request data",
			Error:   err.Error(),
		})
		return
	}

	doctor := models.RecommendedDoctor{
		Name:            req.Name,
		Description:     req.Description,
		ExperienceCount: req.ExperienceCount,
	}

	if err := database.GetDB().Create(&doctor).Error; err != nil {
		c.JSON(http.StatusInternalServerError, vo.ErrorResponse{
			Code:    "INTERNAL_SERVER_ERROR",
			Message: "Failed to create recommended doctor",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, dto.RecommendedDoctorResponse{
		ID:              doctor.ID.String(),
		Name:            doctor.Name,
		Description:     doctor.Description,
		ExperienceCount: doctor.ExperienceCount,
		CreatedAt:       doctor.CreatedAt,
		UpdatedAt:       doctor.UpdatedAt,
	})
}

// UpdateRecommendedDoctor 更新推薦醫師
// @Summary 更新推薦醫師
// @Description 更新推薦醫師資訊
// @Tags recommended-doctors
// @Accept json
// @Produce json
// @Param id path string true "推薦醫師ID"
// @Param doctor body dto.RecommendedDoctorRequest true "推薦醫師資訊"
// @Success 200 {object} dto.RecommendedDoctorResponse
// @Failure 400 {object} vo.ErrorResponse
// @Failure 404 {object} vo.ErrorResponse
// @Failure 500 {object} vo.ErrorResponse
// @Router /admin/recommended-doctors/{id} [put]
func UpdateRecommendedDoctor(c *gin.Context) {
	id := c.Param("id")
	doctorID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, vo.ErrorResponse{
			Code:    "BAD_REQUEST",
			Message: "Invalid recommended doctor ID",
			Error:   err.Error(),
		})
		return
	}

	var req dto.RecommendedDoctorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, vo.ErrorResponse{
			Code:    "BAD_REQUEST",
			Message: "Invalid request data",
			Error:   err.Error(),
		})
		return
	}

	var doctor models.RecommendedDoctor
	if err := database.GetDB().First(&doctor, "id = ?", doctorID).Error; err != nil {
		c.JSON(http.StatusNotFound, vo.ErrorResponse{
			Code:    "NOT_FOUND",
			Message: "Recommended doctor not found",
			Error:   err.Error(),
		})
		return
	}

	// 更新欄位
	doctor.Name = req.Name
	doctor.Description = req.Description
	doctor.ExperienceCount = req.ExperienceCount

	if err := database.GetDB().Save(&doctor).Error; err != nil {
		c.JSON(http.StatusInternalServerError, vo.ErrorResponse{
			Code:    "INTERNAL_SERVER_ERROR",
			Message: "Failed to update recommended doctor",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.RecommendedDoctorResponse{
		ID:              doctor.ID.String(),
		Name:            doctor.Name,
		Description:     doctor.Description,
		ExperienceCount: doctor.ExperienceCount,
		CreatedAt:       doctor.CreatedAt,
		UpdatedAt:       doctor.UpdatedAt,
	})
}

// DeleteRecommendedDoctor 刪除推薦醫師
// @Summary 刪除推薦醫師
// @Description 刪除推薦醫師記錄
// @Tags recommended-doctors
// @Accept json
// @Produce json
// @Param id path string true "推薦醫師ID"
// @Success 200 {object} vo.Response
// @Failure 400 {object} vo.ErrorResponse
// @Failure 500 {object} vo.ErrorResponse
// @Router /admin/recommended-doctors/{id} [delete]
func DeleteRecommendedDoctor(c *gin.Context) {
	id := c.Param("id")
	doctorID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, vo.ErrorResponse{
			Code:    "BAD_REQUEST",
			Message: "Invalid recommended doctor ID",
			Error:   err.Error(),
		})
		return
	}

	if err := database.GetDB().Delete(&models.RecommendedDoctor{}, "id = ?", doctorID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, vo.ErrorResponse{
			Code:    "INTERNAL_SERVER_ERROR",
			Message: "Failed to delete recommended doctor",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, vo.SuccessResponse(nil, "Recommended doctor deleted successfully"))
}
