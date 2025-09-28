package handlers

import (
	"net/http"
	"strconv"

	"mindhelp-backend/internal/database"
	"mindhelp-backend/internal/dto"
	"mindhelp-backend/internal/middleware"
	"mindhelp-backend/internal/models"
	"mindhelp-backend/internal/vo"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// LocationHandler 位置處理器
type LocationHandler struct{}

// NewLocationHandler 創建新的位置處理器
func NewLocationHandler() *LocationHandler {
	return &LocationHandler{}
}

// getDB 獲取資料庫連接，如果失敗會向客戶端返回錯誤
func (h *LocationHandler) getDB(c *gin.Context) (*gorm.DB, bool) {
	db, err := database.GetDBSafely()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, vo.NewErrorResponse(
			"database_unavailable",
			"資料庫暫時無法使用，請稍後再試",
			"DATABASE_UNAVAILABLE",
			nil,
			c.Request.URL.Path,
		))
		return nil, false
	}
	return db, true
}

// CreateLocation 創建位置
// @Summary 創建位置
// @Description 創建新的心理健康資源位置
// @Tags location

// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.LocationRequest true "位置資訊"
// @Success 201 {object} vo.Response{data=dto.LocationResponse}
// @Failure 400 {object} vo.ErrorResponse
// @Failure 401 {object} vo.ErrorResponse
// @Router /locations [post]
func (h *LocationHandler) CreateLocation(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, vo.NewErrorResponse(
			"unauthorized",
			"User not authenticated",
			"UNAUTHORIZED",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	var req dto.LocationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, vo.NewErrorResponse(
			"bad_request",
			"Invalid request data",
			"VALIDATION_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	// 驗證請求資料
	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, vo.NewErrorResponse(
			"bad_request",
			"Validation failed",
			"VALIDATION_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	// 創建位置
	location := models.Location{
		UserID:      uuid.MustParse(userID),
		Name:        req.Name,
		Description: req.Description,
		Address:     req.Address,
		Latitude:    req.Latitude,
		Longitude:   req.Longitude,
		Category:    req.Category,
		Phone:       req.Phone,
		Website:     req.Website,
		Rating:      req.Rating,
		IsPublic:    req.IsPublic,
	}

	// 獲取資料庫連接
	db, ok := h.getDB(c)
	if !ok {
		return
	}

	if err := db.Create(&location).Error; err != nil {
		c.JSON(http.StatusInternalServerError, vo.NewErrorResponse(
			"internal_error",
			"Failed to create location",
			"INTERNAL_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	// 構建回應
	response := dto.LocationResponse{
		ID:          location.ID.String(),
		UserID:      location.UserID.String(),
		Name:        location.Name,
		Description: location.Description,
		Address:     location.Address,
		Latitude:    location.Latitude,
		Longitude:   location.Longitude,
		Category:    location.Category,
		Phone:       location.Phone,
		Website:     location.Website,
		Rating:      location.Rating,
		IsPublic:    location.IsPublic,
		CreatedAt:   location.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   location.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	c.JSON(http.StatusCreated, vo.SuccessResponse(response, "Location created successfully"))
}

// SearchLocations 搜尋位置
// @Summary 搜尋位置
// @Description 搜尋附近的心理健康資源位置
// @Tags location
// @Accept json
// @Produce json
// @Param query query string false "搜尋關鍵字"
// @Param latitude query float64 false "緯度"
// @Param longitude query float64 false "經度"
// @Param radius query float64 false "搜尋半徑(公里)" default(10)
// @Param category query string false "類別"
// @Param page query int false "頁碼" default(1)
// @Param page_size query int false "每頁數量" default(20)
// @Success 200 {object} vo.Response{data=dto.LocationSearchResponse}
// @Failure 400 {object} vo.ErrorResponse
// @Router /locations/search [get]
func (h *LocationHandler) SearchLocations(c *gin.Context) {
	// 獲取查詢參數
	query := c.Query("query")
	latitudeStr := c.Query("latitude")
	longitudeStr := c.Query("longitude")
	radiusStr := c.DefaultQuery("radius", "10")
	category := c.Query("category")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	radius, _ := strconv.ParseFloat(radiusStr, 64)
	if radius <= 0 {
		radius = 10
	}

	offset := (page - 1) * pageSize

	// 獲取資料庫連接
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

	// 構建查詢
	dbQuery := db.Model(&models.Location{}).Where("is_public = ?", true)

	// 添加關鍵字搜尋
	if query != "" {
		dbQuery = dbQuery.Where("name ILIKE ? OR description ILIKE ?", "%"+query+"%", "%"+query+"%")
	}

	// 添加類別篩選
	if category != "" {
		dbQuery = dbQuery.Where("category = ?", category)
	}

	// 如果有座標，添加距離計算
	if latitudeStr != "" && longitudeStr != "" {
		latitude, _ := strconv.ParseFloat(latitudeStr, 64)
		longitude, _ := strconv.ParseFloat(longitudeStr, 64)

		// 使用 Haversine 公式計算距離
		// 這裡簡化處理，實際應用中可能需要更精確的距離計算
		dbQuery = dbQuery.Where(
			"SQRT(POWER(latitude - ?, 2) + POWER(longitude - ?, 2)) <= ?",
			latitude, longitude, radius/111.0, // 粗略轉換：1度約等於111公里
		)
	}

	// 獲取總數
	var total int64
	if err := dbQuery.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, vo.NewErrorResponse(
			"internal_error",
			"Failed to get location count",
			"INTERNAL_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	// 獲取分頁資料
	var locations []models.Location
	if err := dbQuery.Offset(offset).Limit(pageSize).Find(&locations).Error; err != nil {
		c.JSON(http.StatusInternalServerError, vo.NewErrorResponse(
			"internal_error",
			"Failed to get locations",
			"INTERNAL_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	// 轉換為 DTO
	var locationResponses []dto.LocationResponse
	for _, loc := range locations {
		response := dto.LocationResponse{
			ID:          loc.ID.String(),
			UserID:      loc.UserID.String(),
			Name:        loc.Name,
			Description: loc.Description,
			Address:     loc.Address,
			Latitude:    loc.Latitude,
			Longitude:   loc.Longitude,
			Category:    loc.Category,
			Phone:       loc.Phone,
			Website:     loc.Website,
			Rating:      loc.Rating,
			IsPublic:    loc.IsPublic,
			CreatedAt:   loc.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:   loc.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
		locationResponses = append(locationResponses, response)
	}

	// 構建分頁回應
	searchResponse := dto.LocationSearchResponse{
		Locations: locationResponses,
		Total:     total,
		Page:      page,
		PageSize:  pageSize,
		HasMore:   offset+pageSize < int(total),
	}

	c.JSON(http.StatusOK, vo.SuccessResponse(searchResponse, "Locations found successfully"))
}

// GetLocation 獲取位置詳情
// @Summary 獲取位置詳情
// @Description 獲取特定位置的詳細資訊
// @Tags location
// @Accept json
// @Produce json
// @Param id path string true "位置ID"
// @Success 200 {object} vo.Response{data=dto.LocationResponse}
// @Failure 400 {object} vo.ErrorResponse
// @Failure 404 {object} vo.ErrorResponse
// @Router /locations/{id} [get]
func (h *LocationHandler) GetLocation(c *gin.Context) {
	locationID := c.Param("id")
	if locationID == "" {
		c.JSON(http.StatusBadRequest, vo.NewErrorResponse(
			"bad_request",
			"Location ID is required",
			"VALIDATION_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	// 驗證 UUID 格式
	parsedID, err := uuid.Parse(locationID)
	if err != nil {
		c.JSON(http.StatusBadRequest, vo.NewErrorResponse(
			"bad_request",
			"Invalid location ID format",
			"VALIDATION_ERROR",
			nil,
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
			nil,
			c.Request.URL.Path,
		))
		return
	}

	var location models.Location
	if err := db.Where("id = ? AND is_public = ?", parsedID, true).First(&location).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, vo.NewErrorResponse(
				"not_found",
				"Location not found",
				"NOT_FOUND",
				nil,
				c.Request.URL.Path,
			))
			return
		}
		c.JSON(http.StatusInternalServerError, vo.NewErrorResponse(
			"internal_error",
			"Failed to get location",
			"INTERNAL_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	// 構建回應
	response := dto.LocationResponse{
		ID:          location.ID.String(),
		UserID:      location.UserID.String(),
		Name:        location.Name,
		Description: location.Description,
		Address:     location.Address,
		Latitude:    location.Latitude,
		Longitude:   location.Longitude,
		Category:    location.Category,
		Phone:       location.Phone,
		Website:     location.Website,
		Rating:      location.Rating,
		IsPublic:    location.IsPublic,
		CreatedAt:   location.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   location.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	c.JSON(http.StatusOK, vo.SuccessResponse(response, "Location retrieved successfully"))
}

// UpdateLocation 更新位置
// @Summary 更新位置
// @Description 更新現有位置資訊
// @Tags location
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "位置ID"
// @Param request body dto.LocationUpdateRequest true "更新資訊"
// @Success 200 {object} vo.Response{data=dto.LocationResponse}
// @Failure 400 {object} vo.ErrorResponse
// @Failure 401 {object} vo.ErrorResponse
// @Failure 403 {object} vo.ErrorResponse
// @Failure 404 {object} vo.ErrorResponse
// @Router /locations/{id} [put]
func (h *LocationHandler) UpdateLocation(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, vo.NewErrorResponse(
			"unauthorized",
			"User not authenticated",
			"UNAUTHORIZED",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	locationID := c.Param("id")
	if locationID == "" {
		c.JSON(http.StatusBadRequest, vo.NewErrorResponse(
			"bad_request",
			"Location ID is required",
			"VALIDATION_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	// 驗證 UUID 格式
	parsedID, err := uuid.Parse(locationID)
	if err != nil {
		c.JSON(http.StatusBadRequest, vo.NewErrorResponse(
			"bad_request",
			"Invalid location ID format",
			"VALIDATION_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	var req dto.LocationUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, vo.NewErrorResponse(
			"bad_request",
			"Invalid request data",
			"VALIDATION_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	// 驗證請求資料
	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, vo.NewErrorResponse(
			"bad_request",
			"Validation failed",
			"VALIDATION_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	// 查找位置
	var location models.Location
	// 獲取資料庫連接
	db, ok := h.getDB(c)
	if !ok {
		return
	}

	if err := db.Where("id = ?", parsedID).First(&location).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, vo.NewErrorResponse(
				"not_found",
				"Location not found",
				"NOT_FOUND",
				nil,
				c.Request.URL.Path,
			))
			return
		}
		c.JSON(http.StatusInternalServerError, vo.NewErrorResponse(
			"internal_error",
			"Failed to get location",
			"INTERNAL_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	// 檢查權限
	if location.UserID.String() != userID {
		c.JSON(http.StatusForbidden, vo.NewErrorResponse(
			"forbidden",
			"You can only update your own locations",
			"FORBIDDEN",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	// 更新欄位
	updates := make(map[string]interface{})
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.Address != nil {
		updates["address"] = *req.Address
	}
	if req.Latitude != nil {
		updates["latitude"] = *req.Latitude
	}
	if req.Longitude != nil {
		updates["longitude"] = *req.Longitude
	}
	if req.Category != nil {
		updates["category"] = *req.Category
	}
	if req.Phone != nil {
		updates["phone"] = *req.Phone
	}
	if req.Website != nil {
		updates["website"] = *req.Website
	}
	if req.Rating != nil {
		updates["rating"] = *req.Rating
	}
	if req.IsPublic != nil {
		updates["is_public"] = *req.IsPublic
	}

	// 執行更新
	if err := db.Model(&location).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, vo.NewErrorResponse(
			"internal_error",
			"Failed to update location",
			"INTERNAL_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	// 重新獲取更新後的位置（防止 SQL Injection，明確指定主鍵查詢）
	if err := db.Where("id = ?", locationID).First(&location).Error; err != nil {
		c.JSON(http.StatusInternalServerError, vo.NewErrorResponse(
			"internal_error",
			"Failed to get updated location",
			"INTERNAL_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	// 構建回應
	response := dto.LocationResponse{
		ID:          location.ID.String(),
		UserID:      location.UserID.String(),
		Name:        location.Name,
		Description: location.Description,
		Address:     location.Address,
		Latitude:    location.Latitude,
		Longitude:   location.Longitude,
		Category:    location.Category,
		Phone:       location.Phone,
		Website:     location.Website,
		Rating:      location.Rating,
		IsPublic:    location.IsPublic,
		CreatedAt:   location.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   location.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	c.JSON(http.StatusOK, vo.SuccessResponse(response, "Location updated successfully"))
}

// DeleteLocation 刪除位置
// @Summary 刪除位置
// @Description 刪除指定的位置
// @Tags location
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "位置ID"
// @Success 200 {object} vo.Response
// @Failure 400 {object} vo.ErrorResponse
// @Failure 401 {object} vo.ErrorResponse
// @Failure 403 {object} vo.ErrorResponse
// @Failure 404 {object} vo.ErrorResponse
// @Router /locations/{id} [delete]
func (h *LocationHandler) DeleteLocation(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, vo.NewErrorResponse(
			"unauthorized",
			"User not authenticated",
			"UNAUTHORIZED",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	locationID := c.Param("id")
	if locationID == "" {
		c.JSON(http.StatusBadRequest, vo.NewErrorResponse(
			"bad_request",
			"Location ID is required",
			"VALIDATION_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	// 驗證 UUID 格式
	parsedID, err := uuid.Parse(locationID)
	if err != nil {
		c.JSON(http.StatusBadRequest, vo.NewErrorResponse(
			"bad_request",
			"Invalid location ID format",
			"VALIDATION_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	// 查找位置
	var location models.Location
	// 獲取資料庫連接
	db, ok := h.getDB(c)
	if !ok {
		return
	}

	if err := db.Where("id = ?", parsedID).First(&location).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, vo.NewErrorResponse(
				"not_found",
				"Location not found",
				"NOT_FOUND",
				nil,
				c.Request.URL.Path,
			))
			return
		}
		c.JSON(http.StatusInternalServerError, vo.NewErrorResponse(
			"internal_error",
			"Failed to get location",
			"INTERNAL_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	// 檢查權限
	if location.UserID.String() != userID {
		c.JSON(http.StatusForbidden, vo.NewErrorResponse(
			"forbidden",
			"You can only delete your own locations",
			"FORBIDDEN",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	// 軟刪除位置
	if err := db.Delete(&location).Error; err != nil {
		c.JSON(http.StatusInternalServerError, vo.NewErrorResponse(
			"internal_error",
			"Failed to delete location",
			"INTERNAL_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	c.JSON(http.StatusOK, vo.SuccessResponse(nil, "Location deleted successfully"))
}
