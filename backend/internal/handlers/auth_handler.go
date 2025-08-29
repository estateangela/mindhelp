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
	"github.com/golang-jwt/jwt/v5"

	//"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// AuthHandler 認證處理器
type AuthHandler struct {
	db  *gorm.DB
	cfg *config.Config
}

// NewAuthHandler 創建新的認證處理器
func NewAuthHandler(cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		db:  database.GetDB(),
		cfg: cfg,
	}
}

// Register 使用者註冊
// @Summary 使用者註冊
// @Description 創建新的使用者帳戶
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "註冊資訊"
// @Success 201 {object} vo.Response{data=dto.AuthResponse}
// @Failure 400 {object} vo.ErrorResponse
// @Failure 409 {object} vo.ErrorResponse
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
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

	// 檢查 email 是否已存在
	var existingUser models.User
	if err := h.db.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, vo.NewErrorResponse(
			"conflict",
			"Email already exists",
			"EMAIL_EXISTS",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	// 檢查 username 是否已存在
	if err := h.db.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, vo.NewErrorResponse(
			"conflict",
			"Username already exists",
			"USERNAME_EXISTS",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	// 加密密碼
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, vo.NewErrorResponse(
			"internal_error",
			"Failed to process password",
			"INTERNAL_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	// 創建使用者
	user := models.User{
		Email:    req.Email,
		Password: string(hashedPassword),
		Username: req.Username,
		FullName: req.FullName,
		Phone:    req.Phone,
		Avatar:   req.Avatar,
		IsActive: true,
	}

	if err := h.db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, vo.NewErrorResponse(
			"internal_error",
			"Failed to create user",
			"INTERNAL_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	// 生成 JWT token
	accessToken, err := middleware.GenerateToken(h.cfg, user.ID.String(), user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, vo.NewErrorResponse(
			"internal_error",
			"Failed to generate token",
			"INTERNAL_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	refreshToken, err := middleware.GenerateRefreshToken(h.cfg, user.ID.String(), user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, vo.NewErrorResponse(
			"internal_error",
			"Failed to generate refresh token",
			"INTERNAL_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	// 構建回應
	authResponse := dto.AuthResponse{
		User: dto.UserResponse{
			ID:        user.ID.String(),
			Email:     user.Email,
			Username:  user.Username,
			FullName:  user.FullName,
			Phone:     user.Phone,
			Avatar:    user.Avatar,
			IsActive:  user.IsActive,
			CreatedAt: user.CreatedAt.Format(time.RFC3339),
		},
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(h.cfg.JWT.Expiry.Seconds()),
	}

	c.JSON(http.StatusCreated, vo.SuccessResponse(authResponse, "User registered successfully"))
}

// Login 使用者登入
// @Summary 使用者登入
// @Description 使用者登入並獲取 JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "登入資訊"
// @Success 200 {object} vo.Response{data=dto.AuthResponse}
// @Failure 400 {object} vo.ErrorResponse
// @Failure 401 {object} vo.ErrorResponse
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
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
	if err := h.db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, vo.NewErrorResponse(
			"unauthorized",
			"Invalid email or password",
			"INVALID_CREDENTIALS",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	// 檢查使用者是否啟用
	if !user.IsActive {
		c.JSON(http.StatusUnauthorized, vo.NewErrorResponse(
			"unauthorized",
			"Account is deactivated",
			"ACCOUNT_DEACTIVATED",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	// 驗證密碼
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, vo.NewErrorResponse(
			"unauthorized",
			"Invalid email or password",
			"INVALID_CREDENTIALS",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	// 更新最後登入時間
	now := time.Now()
	user.LastLogin = &now
	h.db.Save(&user)

	// 生成 JWT token
	accessToken, err := middleware.GenerateToken(h.cfg, user.ID.String(), user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, vo.NewErrorResponse(
			"internal_error",
			"Failed to generate token",
			"INTERNAL_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	refreshToken, err := middleware.GenerateRefreshToken(h.cfg, user.ID.String(), user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, vo.NewErrorResponse(
			"internal_error",
			"Failed to generate refresh token",
			"INTERNAL_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	// 構建回應
	authResponse := dto.AuthResponse{
		User: dto.UserResponse{
			ID:       user.ID.String(),
			Email:    user.Email,
			Username: user.Username,
			FullName: user.FullName,
			Phone:    user.Phone,
			Avatar:   user.Avatar,
			IsActive: user.IsActive,
			//
			// backend/internal/handlers/auth_handler.go (286-292)

			LastLogin: func() *string {
				if user.LastLogin != nil {
					formatted := user.LastLogin.Format(time.RFC3339)
					return &formatted
				}
				return nil
			}(),

			//
			CreatedAt: user.CreatedAt.Format(time.RFC3339),
		},
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(h.cfg.JWT.Expiry.Seconds()),
	}

	c.JSON(http.StatusOK, vo.SuccessResponse(authResponse, "Login successful"))
}

// RefreshToken 刷新 JWT token
// @Summary 刷新 JWT token
// @Description 使用 refresh token 獲取新的 access token
// @Tags auth
// @Accept json
// @Produce json
// @Param refresh_token body string true "Refresh token"
// @Success 200 {object} vo.Response{data=dto.AuthResponse}
// @Failure 400 {object} vo.ErrorResponse
// @Failure 401 {object} vo.ErrorResponse
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, vo.NewErrorResponse(
			"bad_request",
			"Refresh token is required",
			"VALIDATION_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	// 解析 refresh token
	token, err := jwt.ParseWithClaims(req.RefreshToken, &middleware.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(h.cfg.JWT.Secret), nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, vo.NewErrorResponse(
			"unauthorized",
			"Invalid refresh token",
			"INVALID_TOKEN",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	claims, ok := token.Claims.(*middleware.Claims)
	if !ok {
		c.JSON(http.StatusUnauthorized, vo.NewErrorResponse(
			"unauthorized",
			"Invalid token claims",
			"INVALID_TOKEN",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	// 生成新的 access token
	accessToken, err := middleware.GenerateToken(h.cfg, claims.UserID, claims.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, vo.NewErrorResponse(
			"internal_error",
			"Failed to generate token",
			"INTERNAL_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	// 構建回應
	authResponse := dto.AuthResponse{
		AccessToken: accessToken,
		ExpiresIn:   int64(h.cfg.JWT.Expiry.Seconds()),
	}

	c.JSON(http.StatusOK, vo.SuccessResponse(authResponse, "Token refreshed successfully"))
}
