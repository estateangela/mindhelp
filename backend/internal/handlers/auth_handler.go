package handlers

import (
	"net/http"

	"mindhelp-backend/internal/config"
	"mindhelp-backend/internal/vo"

	"github.com/gin-gonic/gin"
)

// AuthHandler 認證處理器
type AuthHandler struct {
	cfg *config.Config
}

// NewAuthHandler 創建新的認證處理器
func NewAuthHandler(cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		cfg: cfg,
	}
}

// Register 使用者註冊
// @Summary 使用者註冊
// @Description 創建新的使用者帳戶
// @Tags auth
// @Accept json
// @Produce json
// @Success 201 {object} vo.Response
// @Failure 400 {object} vo.ErrorResponse
// @Failure 409 {object} vo.ErrorResponse
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, vo.NewErrorResponse(
		"not_implemented",
		"此功能尚未實現",
		"NOT_IMPLEMENTED",
		nil,
		c.Request.URL.Path,
	))
}

// Login 使用者登入
// @Summary 使用者登入
// @Description 使用者登入驗證
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} vo.Response
// @Failure 400 {object} vo.ErrorResponse
// @Failure 401 {object} vo.ErrorResponse
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, vo.NewErrorResponse(
		"not_implemented",
		"此功能尚未實現",
		"NOT_IMPLEMENTED",
		nil,
		c.Request.URL.Path,
	))
}

// RefreshToken 重新整理 Token
// @Summary 重新整理 Token
// @Description 使用 refresh token 獲取新的 access token
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} vo.Response
// @Failure 400 {object} vo.ErrorResponse
// @Failure 401 {object} vo.ErrorResponse
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, vo.NewErrorResponse(
		"not_implemented",
		"此功能尚未實現",
		"NOT_IMPLEMENTED",
		nil,
		c.Request.URL.Path,
	))
}