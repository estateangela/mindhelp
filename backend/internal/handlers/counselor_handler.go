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

// GetCounselors 獲取諮商師列表
// @Summary 獲取諮商師列表
// @Description 獲取諮商師列表，支援分頁和搜索
// @Tags counselors
// @Accept json
// @Produce json
// @Param page query int false "頁碼" default(1)
// @Param page_size query int false "每頁數量" default(10)
// @Param search query string false "搜索關鍵字"
// @Param work_location query string false "工作地點"
// @Param specialty query string false "專業領域"
// @Success 200 {object} dto.CounselorListResponse
// @Failure 500 {object} vo.ErrorResponse
// @Router /counselors [get]
func GetCounselors(c *gin.Context) {
	// 解析查詢參數
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	search := c.Query("search")
	workLocation := c.Query("work_location")
	specialty := c.Query("specialty")

	// 確保分頁參數合理
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
			"Database service is currently unavailable",
			"SERVICE_UNAVAILABLE",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	// 構建查詢
	query := db.Model(&models.Counselor{})

	// 添加搜索條件
	if search != "" {
		query = query.Where("name ILIKE ? OR license_number ILIKE ? OR specialties ILIKE ?", 
			"%"+search+"%", "%"+search+"%", "%"+search+"%")
	}
	if workLocation != "" {
		query = query.Where("work_location ILIKE ?", "%"+workLocation+"%")
	}
	if specialty != "" {
		query = query.Where("specialties ILIKE ?", "%"+specialty+"%")
	}

	// 獲取總數
	var total int64
	if err := query.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, vo.ErrorResponse{
			Code:    "INTERNAL_SERVER_ERROR",
			Message: "Failed to count counselors",
			Error:   err.Error(),
		})
		return
	}

	// 獲取諮商師列表
	var counselors []models.Counselor
	if err := query.Offset(offset).Limit(pageSize).Find(&counselors).Error; err != nil {
		c.JSON(http.StatusInternalServerError, vo.ErrorResponse{
			Code:    "INTERNAL_SERVER_ERROR",
			Message: "Failed to fetch counselors",
			Error:   err.Error(),
		})
		return
	}

	// 轉換為回應格式
	var counselorResponses []dto.CounselorResponse
	for _, counselor := range counselors {
		counselorResponses = append(counselorResponses, dto.CounselorResponse{
			ID:               counselor.ID.String(),
			Name:             counselor.Name,
			LicenseNumber:    counselor.LicenseNumber,
			Gender:           counselor.Gender,
			Specialties:      counselor.Specialties,
			LanguageSkills:   counselor.LanguageSkills,
			WorkLocation:     counselor.WorkLocation,
			WorkUnit:         counselor.WorkUnit,
			InstitutionCode:  counselor.InstitutionCode,
			PsychologySchool: counselor.PsychologySchool,
			TreatmentMethods: counselor.TreatmentMethods,
			CreatedAt:        counselor.CreatedAt,
			UpdatedAt:        counselor.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, dto.CounselorListResponse{
		Counselors: counselorResponses,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
	})
}

// GetCounselor 獲取單個諮商師
// @Summary 獲取單個諮商師
// @Description 根據ID獲取諮商師詳細資訊
// @Tags counselors
// @Accept json
// @Produce json
// @Param id path string true "諮商師ID"
// @Success 200 {object} dto.CounselorResponse
// @Failure 400 {object} vo.ErrorResponse
// @Failure 404 {object} vo.ErrorResponse
// @Router /counselors/{id} [get]
func GetCounselor(c *gin.Context) {
	id := c.Param("id")
	counselorID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, vo.ErrorResponse{
			Code:    "BAD_REQUEST",
			Message: "Invalid counselor ID",
			Error:   err.Error(),
		})
		return
	}

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

	var counselor models.Counselor
	if err := db.First(&counselor, "id = ?", counselorID).Error; err != nil {
		c.JSON(http.StatusNotFound, vo.ErrorResponse{
			Code:    "NOT_FOUND",
			Message: "Counselor not found",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.CounselorResponse{
		ID:               counselor.ID.String(),
		Name:             counselor.Name,
		LicenseNumber:    counselor.LicenseNumber,
		Gender:           counselor.Gender,
		Specialties:      counselor.Specialties,
		LanguageSkills:   counselor.LanguageSkills,
		WorkLocation:     counselor.WorkLocation,
		WorkUnit:         counselor.WorkUnit,
		InstitutionCode:  counselor.InstitutionCode,
		PsychologySchool: counselor.PsychologySchool,
		TreatmentMethods: counselor.TreatmentMethods,
		CreatedAt:        counselor.CreatedAt,
		UpdatedAt:        counselor.UpdatedAt,
	})
}

// CreateCounselor 創建諮商師
// @Summary 創建諮商師
// @Description 創建新的諮商師記錄
// @Tags counselors
// @Accept json
// @Produce json
// @Param counselor body dto.CounselorRequest true "諮商師資訊"
// @Success 201 {object} dto.CounselorResponse
// @Failure 400 {object} vo.ErrorResponse
// @Failure 500 {object} vo.ErrorResponse
// @Router /admin/counselors [post]
func CreateCounselor(c *gin.Context) {
	var req dto.CounselorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, vo.ErrorResponse{
			Code:    "BAD_REQUEST",
			Message: "Invalid request data",
			Error:   err.Error(),
		})
		return
	}

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

	counselor := models.Counselor{
		Name:             req.Name,
		LicenseNumber:    req.LicenseNumber,
		Gender:           req.Gender,
		Specialties:      req.Specialties,
		LanguageSkills:   req.LanguageSkills,
		WorkLocation:     req.WorkLocation,
		WorkUnit:         req.WorkUnit,
		InstitutionCode:  req.InstitutionCode,
		PsychologySchool: req.PsychologySchool,
		TreatmentMethods: req.TreatmentMethods,
	}

	if err := db.Create(&counselor).Error; err != nil {
		c.JSON(http.StatusInternalServerError, vo.ErrorResponse{
			Code:    "INTERNAL_SERVER_ERROR",
			Message: "Failed to create counselor",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, dto.CounselorResponse{
		ID:               counselor.ID.String(),
		Name:             counselor.Name,
		LicenseNumber:    counselor.LicenseNumber,
		Gender:           counselor.Gender,
		Specialties:      counselor.Specialties,
		LanguageSkills:   counselor.LanguageSkills,
		WorkLocation:     counselor.WorkLocation,
		WorkUnit:         counselor.WorkUnit,
		InstitutionCode:  counselor.InstitutionCode,
		PsychologySchool: counselor.PsychologySchool,
		TreatmentMethods: counselor.TreatmentMethods,
		CreatedAt:        counselor.CreatedAt,
		UpdatedAt:        counselor.UpdatedAt,
	})
}

// UpdateCounselor 更新諮商師
// @Summary 更新諮商師
// @Description 更新諮商師資訊
// @Tags counselors
// @Accept json
// @Produce json
// @Param id path string true "諮商師ID"
// @Param counselor body dto.CounselorRequest true "諮商師資訊"
// @Success 200 {object} dto.CounselorResponse
// @Failure 400 {object} vo.ErrorResponse
// @Failure 404 {object} vo.ErrorResponse
// @Failure 500 {object} vo.ErrorResponse
// @Router /admin/counselors/{id} [put]
func UpdateCounselor(c *gin.Context) {
	id := c.Param("id")
	counselorID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, vo.ErrorResponse{
			Code:    "BAD_REQUEST",
			Message: "Invalid counselor ID",
			Error:   err.Error(),
		})
		return
	}

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

	var req dto.CounselorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, vo.ErrorResponse{
			Code:    "BAD_REQUEST",
			Message: "Invalid request data",
			Error:   err.Error(),
		})
		return
	}

	var counselor models.Counselor
	if err := db.First(&counselor, "id = ?", counselorID).Error; err != nil {
		c.JSON(http.StatusNotFound, vo.ErrorResponse{
			Code:    "NOT_FOUND",
			Message: "Counselor not found",
			Error:   err.Error(),
		})
		return
	}

	// 更新欄位
	counselor.Name = req.Name
	counselor.LicenseNumber = req.LicenseNumber
	counselor.Gender = req.Gender
	counselor.Specialties = req.Specialties
	counselor.LanguageSkills = req.LanguageSkills
	counselor.WorkLocation = req.WorkLocation
	counselor.WorkUnit = req.WorkUnit
	counselor.InstitutionCode = req.InstitutionCode
	counselor.PsychologySchool = req.PsychologySchool
	counselor.TreatmentMethods = req.TreatmentMethods

	if err := db.Save(&counselor).Error; err != nil {
		c.JSON(http.StatusInternalServerError, vo.ErrorResponse{
			Code:    "INTERNAL_SERVER_ERROR",
			Message: "Failed to update counselor",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.CounselorResponse{
		ID:               counselor.ID.String(),
		Name:             counselor.Name,
		LicenseNumber:    counselor.LicenseNumber,
		Gender:           counselor.Gender,
		Specialties:      counselor.Specialties,
		LanguageSkills:   counselor.LanguageSkills,
		WorkLocation:     counselor.WorkLocation,
		WorkUnit:         counselor.WorkUnit,
		InstitutionCode:  counselor.InstitutionCode,
		PsychologySchool: counselor.PsychologySchool,
		TreatmentMethods: counselor.TreatmentMethods,
		CreatedAt:        counselor.CreatedAt,
		UpdatedAt:        counselor.UpdatedAt,
	})
}

// DeleteCounselor 刪除諮商師
// @Summary 刪除諮商師
// @Description 刪除諮商師記錄
// @Tags counselors
// @Accept json
// @Produce json
// @Param id path string true "諮商師ID"
// @Success 200 {object} vo.Response
// @Failure 400 {object} vo.ErrorResponse
// @Failure 500 {object} vo.ErrorResponse
// @Router /admin/counselors/{id} [delete]
func DeleteCounselor(c *gin.Context) {
	id := c.Param("id")
	counselorID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, vo.ErrorResponse{
			Code:    "BAD_REQUEST",
			Message: "Invalid counselor ID",
			Error:   err.Error(),
		})
		return
	}

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

	if err := db.Delete(&models.Counselor{}, "id = ?", counselorID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, vo.ErrorResponse{
			Code:    "INTERNAL_SERVER_ERROR",
			Message: "Failed to delete counselor",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, vo.SuccessResponse(nil, "Counselor deleted successfully"))
}
