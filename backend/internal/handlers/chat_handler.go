package handlers

import (
	"net/http"

	"mindhelp-backend/internal/config"
	"mindhelp-backend/internal/vo"

	"github.com/gin-gonic/gin"
)

// ChatHandler 聊天處理器
type ChatHandler struct {
	cfg *config.Config
}

// NewChatHandler 創建新的聊天處理器
func NewChatHandler(cfg *config.Config) *ChatHandler {
	return &ChatHandler{
		cfg: cfg,
	}
}

// SendMessage 發送訊息 (舊版相容)
func (h *ChatHandler) SendMessage(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, vo.NewErrorResponse(
		"not_implemented",
		"此功能尚未實現",
		"NOT_IMPLEMENTED",
		nil,
		c.Request.URL.Path,
	))
}

// GetChatHistory 獲取聊天歷史 (舊版相容)
func (h *ChatHandler) GetChatHistory(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, vo.NewErrorResponse(
		"not_implemented",
		"此功能尚未實現",
		"NOT_IMPLEMENTED",
		nil,
		c.Request.URL.Path,
	))
}

// GetSessions 獲取聊天會話列表
func (h *ChatHandler) GetSessions(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, vo.NewErrorResponse(
		"not_implemented",
		"此功能尚未實現",
		"NOT_IMPLEMENTED",
		nil,
		c.Request.URL.Path,
	))
}

// CreateSession 創建新的聊天會話
func (h *ChatHandler) CreateSession(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, vo.NewErrorResponse(
		"not_implemented",
		"此功能尚未實現",
		"NOT_IMPLEMENTED",
		nil,
		c.Request.URL.Path,
	))
}

// GetSessionMessages 獲取會話訊息
func (h *ChatHandler) GetSessionMessages(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, vo.NewErrorResponse(
		"not_implemented",
		"此功能尚未實現",
		"NOT_IMPLEMENTED",
		nil,
		c.Request.URL.Path,
	))
}

// SendSessionMessage 在會話中發送訊息
func (h *ChatHandler) SendSessionMessage(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, vo.NewErrorResponse(
		"not_implemented",
		"此功能尚未實現",
		"NOT_IMPLEMENTED",
		nil,
		c.Request.URL.Path,
	))
}