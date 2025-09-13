package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"mindhelp-backend/internal/config"
	"mindhelp-backend/internal/database"
	"mindhelp-backend/internal/dto"
	"mindhelp-backend/internal/middleware"
	"mindhelp-backend/internal/models"
	"mindhelp-backend/internal/vo"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ChatHandler 聊天處理器
type ChatHandler struct {
	db  *gorm.DB
	cfg *config.Config
}

// NewChatHandler 創建新的聊天處理器
func NewChatHandler(cfg *config.Config) *ChatHandler {
	return &ChatHandler{
		db:  database.GetDB(),
		cfg: cfg,
	}
}

// SendMessage 發送聊天訊息 (舊版兼容)
// @Summary 發送聊天訊息 (舊版)
// @Description 發送訊息給 AI 並獲取回應 (建議使用 session-based 端點)
// @Tags chat
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.ChatMessageRequest true "聊天訊息"
// @Success 200 {object} vo.Response{data=dto.ChatMessageResponse}
// @Failure 400 {object} vo.ErrorResponse
// @Failure 401 {object} vo.ErrorResponse
// @Failure 500 {object} vo.ErrorResponse
// @Router /chat/send [post]
// @Deprecated
func (h *ChatHandler) SendMessage(c *gin.Context) {
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

	var req dto.ChatMessageRequest
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

	var sessionID *uuid.UUID

	// 如果有提供 SessionID，使用指定的 session
	if req.SessionID != nil && *req.SessionID != "" {
		parsedSessionID, err := uuid.Parse(*req.SessionID)
		if err != nil {
			c.JSON(http.StatusBadRequest, vo.NewErrorResponse(
				"bad_request",
				"Invalid session ID format",
				"VALIDATION_ERROR",
				nil,
				c.Request.URL.Path,
			))
			return
		}
		
		// 驗證 session 是否屬於當前使用者
		var session models.ChatSession
		if err := h.db.Where("id = ? AND user_id = ?", parsedSessionID, userID).First(&session).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, vo.NewErrorResponse(
					"not_found",
					"Chat session not found",
					"SESSION_NOT_FOUND",
					nil,
					c.Request.URL.Path,
				))
				return
			}
			c.JSON(http.StatusInternalServerError, vo.NewErrorResponse(
				"internal_error",
				"Failed to check session",
				"INTERNAL_ERROR",
				nil,
				c.Request.URL.Path,
			))
			return
		}
		sessionID = &parsedSessionID
	} else {
		// 創建新的 session
		newSession := models.ChatSession{
			UserID:                 uuid.MustParse(userID),
			FirstMessageSnippet:    req.Content[:min(len(req.Content), 100)],
			LastUpdatedAt:          time.Now(),
			MessageCount:          0,
			IsActive:              true,
		}
		
		if err := h.db.Create(&newSession).Error; err != nil {
			c.JSON(http.StatusInternalServerError, vo.NewErrorResponse(
				"internal_error",
				"Failed to create chat session",
				"INTERNAL_ERROR",
				nil,
				c.Request.URL.Path,
			))
			return
		}
		sessionID = &newSession.ID
	}

	// 保存使用者訊息
	userMessage := models.ChatMessage{
		UserID:    uuid.MustParse(userID),
		SessionID: sessionID,
		Role:      "user",
		Content:   req.Content,
		Timestamp: time.Now().UnixMilli(),
		Model:     req.Model,
	}

	if err := h.db.Create(&userMessage).Error; err != nil {
		c.JSON(http.StatusInternalServerError, vo.NewErrorResponse(
			"internal_error",
			"Failed to save user message",
			"INTERNAL_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	// 調用 OpenRouter API
	aiResponse, err := h.callOpenRouterAPI(req.Content, req.Model)
	if err != nil {
		c.JSON(http.StatusInternalServerError, vo.NewErrorResponse(
			"internal_error",
			"Failed to get AI response",
			"INTERNAL_ERROR",
			[]string{err.Error()},
			c.Request.URL.Path,
		))
		return
	}

	// 保存 AI 回應
	botMessage := models.ChatMessage{
		UserID:    uuid.MustParse(userID),
		SessionID: sessionID,
		Role:      "bot",
		Content:   aiResponse.Choices[0].Message.Content,
		Timestamp: time.Now().UnixMilli(),
		Model:     req.Model,
		Tokens:    aiResponse.Usage.TotalTokens,
	}

	if err := h.db.Create(&botMessage).Error; err != nil {
		c.JSON(http.StatusInternalServerError, vo.NewErrorResponse(
			"internal_error",
			"Failed to save bot message",
			"INTERNAL_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	// 更新 session 統計
	h.db.Model(&models.ChatSession{}).Where("id = ?", sessionID).
		Updates(map[string]interface{}{
			"last_updated_at": time.Now(),
			"message_count":   gorm.Expr("message_count + ?", 2), // user + bot message
		})

	// 構建回應
	sessionIDStr := ""
	if sessionID != nil {
		sessionIDStr = sessionID.String()
	}
	
	response := dto.ChatMessageResponse{
		ID:        botMessage.ID.String(),
		UserID:    botMessage.UserID.String(),
		SessionID: &sessionIDStr,
		Role:      botMessage.Role,
		Content:   botMessage.Content,
		Timestamp: botMessage.Timestamp,
		Model:     botMessage.Model,
		Tokens:    botMessage.Tokens,
		CreatedAt: botMessage.CreatedAt.Format(time.RFC3339),
	}

	c.JSON(http.StatusOK, vo.SuccessResponse(response, "Message sent successfully"))
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// GetChatHistory 獲取聊天歷史
// @Summary 獲取聊天歷史
// @Description 獲取使用者的聊天歷史記錄
// @Tags chat
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "頁碼" default(1)
// @Param page_size query int false "每頁數量" default(20)
// @Success 200 {object} vo.Response{data=dto.ChatHistoryResponse}
// @Failure 400 {object} vo.ErrorResponse
// @Failure 401 {object} vo.ErrorResponse
// @Router /chat/history [get]
func (h *ChatHandler) GetChatHistory(c *gin.Context) {
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

	// 獲取查詢參數
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	// 查詢聊天記錄
	var messages []models.ChatMessage
	var total int64

	// 獲取總數
	if err := h.db.Model(&models.ChatMessage{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, vo.NewErrorResponse(
			"internal_error",
			"Failed to get message count",
			"INTERNAL_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	// 獲取分頁資料
	if err := h.db.Where("user_id = ?", userID).
		Order("timestamp DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&messages).Error; err != nil {
		c.JSON(http.StatusInternalServerError, vo.NewErrorResponse(
			"internal_error",
			"Failed to get messages",
			"INTERNAL_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	// 轉換為 DTO
	var messageResponses []dto.ChatMessageResponse
	for _, msg := range messages {
		response := dto.ChatMessageResponse{
			ID:        msg.ID.String(),
			UserID:    msg.UserID.String(),
			Role:      msg.Role,
			Content:   msg.Content,
			Timestamp: msg.Timestamp,
			Model:     msg.Model,
			Tokens:    msg.Tokens,
			CreatedAt: msg.CreatedAt.Format(time.RFC3339),
		}
		messageResponses = append(messageResponses, response)
	}

	// 構建分頁回應
	chatHistoryResponse := dto.ChatHistoryResponse{
		Messages: messageResponses,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		HasMore:  offset+pageSize < int(total),
	}

	c.JSON(http.StatusOK, vo.SuccessResponse(chatHistoryResponse, "Chat history retrieved successfully"))
}

// callOpenRouterAPI 調用 OpenRouter API
func (h *ChatHandler) callOpenRouterAPI(content, model string) (*dto.OpenRouterResponse, error) {
	// 構建請求
	request := dto.OpenRouterRequest{
		Model:       model,
		Messages:    []dto.Message{{Role: "user", Content: content}},
		Temperature: 0.7,
		MaxTokens:   512,
	}

	// 如果沒有指定模型，使用預設模型
	if model == "" {
		request.Model = "google/gemma-3n-e4b-it:free"
	}

	// 序列化請求
	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// 創建 HTTP 請求
	req, err := http.NewRequest("POST", h.cfg.OpenRouter.BaseURL+"/chat/completions", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+h.cfg.OpenRouter.APIKey)

	// 發送請求
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// 檢查回應狀態
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("OpenRouter API returned status: %d", resp.StatusCode)
	}

	// 解析回應
	var response dto.OpenRouterResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response, nil
}
