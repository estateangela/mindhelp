package handlers

import (
	"net/http"
	"time"

	"mindhelp-backend/internal/config"
	"mindhelp-backend/internal/database"
	"mindhelp-backend/internal/dto"
	"mindhelp-backend/internal/middleware"
	"mindhelp-backend/internal/models"
	"mindhelp-backend/internal/vo"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UserHandler 使用者處理器
type UserHandler struct {
	db  *gorm.DB
	cfg *config.Config
}

// NewUserHandler 創建新的使用者處理器
func NewUserHandler(cfg *config.Config) *UserHandler {
	return &UserHandler{
		db:  database.GetDB(),
		cfg: cfg,
	}
}

// GetProfile 獲取當前使用者資料
// @Summary 獲取使用者資料
// @Description 獲取當前登入使用者的詳細資訊
// @Tags user
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} vo.Response{data=dto.UserProfileResponse}
// @Failure 401 {object} vo.ErrorResponse
// @Failure 500 {object} vo.ErrorResponse
// @Router /users/me [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
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

	var user models.User
	if err := h.db.Where("id = ?", userID).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, vo.NewErrorResponse(
				"not_found",
				"User not found",
				"USER_NOT_FOUND",
				nil,
				c.Request.URL.Path,
			))
			return
		}
		c.JSON(http.StatusInternalServerError, vo.NewErrorResponse(
			"internal_error",
			"Failed to get user",
			"INTERNAL_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	// 構建回應
	response := dto.UserProfileResponse{
		ID:       user.ID.String(),
		Email:    user.Email,
		Username: user.Username,
		FullName: user.FullName,
		Phone:    user.Phone,
		Avatar:   user.Avatar,
		IsActive: user.IsActive,
		LastLogin: func() *string {
			if user.LastLogin != nil {
				formatted := user.LastLogin.Format(time.RFC3339)
				return &formatted
			}
			return nil
		}(),
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
	}

	c.JSON(http.StatusOK, vo.SuccessResponse(response, "User profile retrieved successfully"))
}

// UpdateProfile 更新使用者資料
// @Summary 更新使用者資料
// @Description 更新使用者的個人資料
// @Tags user
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.UpdateUserRequest true "更新資料"
// @Success 200 {object} vo.Response{data=dto.UserProfileResponse}
// @Failure 400 {object} vo.ErrorResponse
// @Failure 401 {object} vo.ErrorResponse
// @Failure 409 {object} vo.ErrorResponse
// @Router /users/me [put]
func (h *UserHandler) UpdateProfile(c *gin.Context) {
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

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, vo.NewErrorResponse(
			"bad_request",
			"Invalid request data",
			"VALIDATION_ERROR",
			[]string{err.Error()},
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
			[]string{err.Error()},
			c.Request.URL.Path,
		))
		return
	}

	// 查找使用者
	var user models.User
	if err := h.db.Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, vo.NewErrorResponse(
			"not_found",
			"User not found",
			"USER_NOT_FOUND",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	// 檢查 username 是否已存在 (如果要更新的話)
	if req.Username != nil && *req.Username != user.Username {
		var existingUser models.User
		if err := h.db.Where("username = ? AND id != ?", *req.Username, userID).First(&existingUser).Error; err == nil {
			c.JSON(http.StatusConflict, vo.NewErrorResponse(
				"conflict",
				"Username already exists",
				"USERNAME_EXISTS",
				nil,
				c.Request.URL.Path,
			))
			return
		}
	}

	// 更新欄位
	updates := make(map[string]interface{})
	if req.Username != nil {
		updates["username"] = *req.Username
	}
	if req.FullName != nil {
		updates["full_name"] = *req.FullName
	}
	if req.Phone != nil {
		updates["phone"] = *req.Phone
	}
	if req.Avatar != nil {
		updates["avatar"] = *req.Avatar
	}

	// 執行更新
	if len(updates) > 0 {
		if err := h.db.Model(&user).Updates(updates).Error; err != nil {
			c.JSON(http.StatusInternalServerError, vo.NewErrorResponse(
				"internal_error",
				"Failed to update user",
				"INTERNAL_ERROR",
				nil,
				c.Request.URL.Path,
			))
			return
		}
	}

	// 重新獲取更新後的使用者資料
	if err := h.db.Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, vo.NewErrorResponse(
			"internal_error",
			"Failed to get updated user",
			"INTERNAL_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	// 構建回應
	response := dto.UserProfileResponse{
		ID:       user.ID.String(),
		Email:    user.Email,
		Username: user.Username,
		FullName: user.FullName,
		Phone:    user.Phone,
		Avatar:   user.Avatar,
		IsActive: user.IsActive,
		LastLogin: func() *string {
			if user.LastLogin != nil {
				formatted := user.LastLogin.Format(time.RFC3339)
				return &formatted
			}
			return nil
		}(),
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
	}

	c.JSON(http.StatusOK, vo.SuccessResponse(response, "User profile updated successfully"))
}

// ChangePassword 變更密碼
// @Summary 變更密碼
// @Description 變更使用者密碼
// @Tags user
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.ChangePasswordRequest true "變更密碼資料"
// @Success 200 {object} vo.Response
// @Failure 400 {object} vo.ErrorResponse
// @Failure 401 {object} vo.ErrorResponse
// @Router /users/me/password [put]
func (h *UserHandler) ChangePassword(c *gin.Context) {
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

	var req dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, vo.NewErrorResponse(
			"bad_request",
			"Invalid request data",
			"VALIDATION_ERROR",
			[]string{err.Error()},
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
			[]string{err.Error()},
			c.Request.URL.Path,
		))
		return
	}

	// 查找使用者
	var user models.User
	if err := h.db.Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, vo.NewErrorResponse(
			"not_found",
			"User not found",
			"USER_NOT_FOUND",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	// 驗證當前密碼
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.CurrentPassword)); err != nil {
		c.JSON(http.StatusUnauthorized, vo.NewErrorResponse(
			"unauthorized",
			"Current password is incorrect",
			"INVALID_PASSWORD",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	// 加密新密碼
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, vo.NewErrorResponse(
			"internal_error",
			"Failed to process new password",
			"INTERNAL_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	// 更新密碼
	if err := h.db.Model(&user).Update("password", string(hashedPassword)).Error; err != nil {
		c.JSON(http.StatusInternalServerError, vo.NewErrorResponse(
			"internal_error",
			"Failed to update password",
			"INTERNAL_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	c.JSON(http.StatusOK, vo.SuccessResponse(nil, "Password changed successfully"))
}

// DeleteAccount 刪除帳號
// @Summary 刪除帳號
// @Description 永久刪除使用者帳號及所有相關資料
// @Tags user
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.DeleteAccountRequest true "刪除帳號資料"
// @Success 200 {object} vo.Response
// @Failure 400 {object} vo.ErrorResponse
// @Failure 401 {object} vo.ErrorResponse
// @Router /users/me [delete]
func (h *UserHandler) DeleteAccount(c *gin.Context) {
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

	var req dto.DeleteAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, vo.NewErrorResponse(
			"bad_request",
			"Invalid request data",
			"VALIDATION_ERROR",
			[]string{err.Error()},
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
			[]string{err.Error()},
			c.Request.URL.Path,
		))
		return
	}

	// 查找使用者
	var user models.User
	if err := h.db.Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, vo.NewErrorResponse(
			"not_found",
			"User not found",
			"USER_NOT_FOUND",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	// 驗證密碼
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, vo.NewErrorResponse(
			"unauthorized",
			"Password is incorrect",
			"INVALID_PASSWORD",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	// 軟刪除使用者 (GORM 會自動處理相關資料的級聯刪除)
	if err := h.db.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, vo.NewErrorResponse(
			"internal_error",
			"Failed to delete account",
			"INTERNAL_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	c.JSON(http.StatusOK, vo.SuccessResponse(nil, "Account deleted successfully"))
}

// GetStats 獲取使用者統計資訊
// @Summary 獲取使用者統計
// @Description 獲取使用者的各項統計資訊
// @Tags user
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} vo.Response{data=dto.UserStatsResponse}
// @Failure 401 {object} vo.ErrorResponse
// @Router /users/me/stats [get]
func (h *UserHandler) GetStats(c *gin.Context) {
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

	var stats dto.UserStatsResponse

	// 統計聊天會話數
	h.db.Model(&models.ChatSession{}).Where("user_id = ?", userID).Count(&stats.TotalChatSessions)

	// 統計測驗完成數
	h.db.Model(&models.QuizSubmission{}).Where("user_id = ?", userID).Count(&stats.TotalQuizzes)

	// 統計收藏數
	h.db.Model(&models.Bookmark{}).Where("user_id = ?", userID).Count(&stats.TotalBookmarks)

	// 統計評論數
	h.db.Model(&models.Review{}).Where("user_id = ?", userID).Count(&stats.TotalReviews)

	// 統計未讀通知數
	h.db.Model(&models.Notification{}).Where("user_id = ? AND is_read = ?", userID, false).Count(&stats.UnreadNotifications)

	c.JSON(http.StatusOK, vo.SuccessResponse(stats, "User stats retrieved successfully"))
}
