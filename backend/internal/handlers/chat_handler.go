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
			UserID:              uuid.MustParse(userID),
			FirstMessageSnippet: req.Content[:min(len(req.Content), 100)],
			LastUpdatedAt:       time.Now(),
			MessageCount:        0,
			IsActive:            true,
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

// GetSessions 獲取聊天會話列表
// @Summary 獲取聊天會話列表
// @Description 獲取使用者的聊天會話列表
// @Tags chat
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "頁碼" default(1)
// @Param limit query int false "每頁數量" default(20)
// @Success 200 {object} vo.Response{data=dto.ChatSessionListResponse}
// @Failure 401 {object} vo.ErrorResponse
// @Router /chat/sessions [get]
func (h *ChatHandler) GetSessions(c *gin.Context) {
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

	// 解析查詢參數
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	// 參數驗證
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit

	// 獲取總數
	var total int64
	if err := h.db.Model(&models.ChatSession{}).Where("user_id = ? AND is_active = ?", userID, true).Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, vo.NewErrorResponse(
			"internal_error",
			"Failed to count chat sessions",
			"INTERNAL_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	// 獲取會話列表
	var sessions []models.ChatSession
	if err := h.db.Where("user_id = ? AND is_active = ?", userID, true).
		Order("last_updated_at DESC").
		Offset(offset).Limit(limit).
		Find(&sessions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, vo.NewErrorResponse(
			"internal_error",
			"Failed to get chat sessions",
			"INTERNAL_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	// 轉換為 DTO
	var sessionResponses []dto.ChatSessionResponse
	for _, session := range sessions {
		response := dto.ChatSessionResponse{
			ID:                  session.ID.String(),
			UserID:              session.UserID.String(),
			Title:               session.Title,
			FirstMessageSnippet: session.FirstMessageSnippet,
			LastUpdatedAt:       session.LastUpdatedAt.Format(time.RFC3339),
			MessageCount:        session.MessageCount,
			IsActive:            session.IsActive,
			CreatedAt:           session.CreatedAt.Format(time.RFC3339),
		}
		sessionResponses = append(sessionResponses, response)
	}

	// 構建分頁回應
	totalPages := int((total + int64(limit) - 1) / int64(limit))
	response := dto.ChatSessionListResponse{
		Sessions:   sessionResponses,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
		HasMore:    page < totalPages,
	}

	c.JSON(http.StatusOK, vo.SuccessResponse(response, "Chat sessions retrieved successfully"))
}

// CreateSession 創建聊天會話
// @Summary 創建聊天會話
// @Description 創建新的聊天會話
// @Tags chat
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.ChatSessionRequest true "會話資料"
// @Success 201 {object} vo.Response{data=dto.ChatSessionResponse}
// @Failure 400 {object} vo.ErrorResponse
// @Failure 401 {object} vo.ErrorResponse
// @Router /chat/sessions [post]
func (h *ChatHandler) CreateSession(c *gin.Context) {
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

	var req dto.ChatSessionRequest
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

	// 創建會話
	session := models.ChatSession{
		UserID:        uuid.MustParse(userID),
		Title:         req.Title,
		LastUpdatedAt: time.Now(),
		MessageCount:  0,
		IsActive:      true,
	}

	if err := h.db.Create(&session).Error; err != nil {
		c.JSON(http.StatusInternalServerError, vo.NewErrorResponse(
			"internal_error",
			"Failed to create chat session",
			"INTERNAL_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	// 構建回應
	response := dto.ChatSessionResponse{
		ID:                  session.ID.String(),
		UserID:              session.UserID.String(),
		Title:               session.Title,
		FirstMessageSnippet: session.FirstMessageSnippet,
		LastUpdatedAt:       session.LastUpdatedAt.Format(time.RFC3339),
		MessageCount:        session.MessageCount,
		IsActive:            session.IsActive,
		CreatedAt:           session.CreatedAt.Format(time.RFC3339),
	}

	c.JSON(http.StatusCreated, vo.SuccessResponse(response, "Chat session created successfully"))
}

// GetSessionMessages 獲取會話訊息
// @Summary 獲取會話訊息
// @Description 獲取指定會話的訊息列表
// @Tags chat
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param sessionId path string true "會話ID"
// @Param page query int false "頁碼" default(1)
// @Param limit query int false "每頁數量" default(50)
// @Success 200 {object} vo.Response{data=dto.SessionMessagesResponse}
// @Failure 400 {object} vo.ErrorResponse
// @Failure 401 {object} vo.ErrorResponse
// @Failure 404 {object} vo.ErrorResponse
// @Router /chat/sessions/{sessionId}/messages [get]
func (h *ChatHandler) GetSessionMessages(c *gin.Context) {
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

	sessionID := c.Param("sessionId")
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, vo.NewErrorResponse(
			"bad_request",
			"Session ID is required",
			"VALIDATION_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	// 驗證 UUID 格式
	parsedSessionID, err := uuid.Parse(sessionID)
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

	// 解析查詢參數
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))

	// 參數驗證
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 50
	}

	offset := (page - 1) * limit

	// 獲取總數
	var total int64
	if err := h.db.Model(&models.ChatMessage{}).Where("session_id = ?", parsedSessionID).Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, vo.NewErrorResponse(
			"internal_error",
			"Failed to count messages",
			"INTERNAL_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	// 獲取訊息列表
	var messages []models.ChatMessage
	if err := h.db.Where("session_id = ?", parsedSessionID).
		Order("timestamp ASC"). // 按時間順序排列
		Offset(offset).Limit(limit).
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
	for _, message := range messages {
		sessionIDStr := ""
		if message.SessionID != nil {
			sessionIDStr = message.SessionID.String()
		}

		response := dto.ChatMessageResponse{
			ID:        message.ID.String(),
			UserID:    message.UserID.String(),
			SessionID: &sessionIDStr,
			Role:      message.Role,
			Content:   message.Content,
			Timestamp: message.Timestamp,
			Model:     message.Model,
			Tokens:    message.Tokens,
			CreatedAt: message.CreatedAt.Format(time.RFC3339),
		}
		messageResponses = append(messageResponses, response)
	}

	// 構建分頁回應
	response := dto.SessionMessagesResponse{
		SessionID: sessionID,
		Messages:  messageResponses,
		Total:     total,
		Page:      page,
		Limit:     limit,
		HasMore:   offset+limit < int(total),
	}

	c.JSON(http.StatusOK, vo.SuccessResponse(response, "Session messages retrieved successfully"))
}

// SendSessionMessage 發送會話訊息
// @Summary 發送會話訊息
// @Description 在指定會話中發送訊息給 AI
// @Tags chat
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param sessionId path string true "會話ID"
// @Param request body dto.ChatMessageRequest true "聊天訊息"
// @Success 200 {object} vo.Response{data=dto.ChatMessageResponse}
// @Failure 400 {object} vo.ErrorResponse
// @Failure 401 {object} vo.ErrorResponse
// @Failure 404 {object} vo.ErrorResponse
// @Router /chat/sessions/{sessionId}/messages [post]
func (h *ChatHandler) SendSessionMessage(c *gin.Context) {
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

	sessionID := c.Param("sessionId")
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, vo.NewErrorResponse(
			"bad_request",
			"Session ID is required",
			"VALIDATION_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	// 驗證 UUID 格式
	parsedSessionID, err := uuid.Parse(sessionID)
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
	if err := h.db.Where("id = ? AND user_id = ? AND is_active = ?", parsedSessionID, userID, true).First(&session).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, vo.NewErrorResponse(
				"not_found",
				"Chat session not found or inactive",
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

	// 如果這是第一則訊息，更新 session 的 FirstMessageSnippet
	if session.MessageCount == 0 {
		snippet := req.Content
		if len(snippet) > 100 {
			snippet = snippet[:100]
		}
		h.db.Model(&session).Update("first_message_snippet", snippet)
	}

	// 保存使用者訊息
	userMessage := models.ChatMessage{
		UserID:    uuid.MustParse(userID),
		SessionID: &parsedSessionID,
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

	// 獲取會話歷史以提供更好的上下文
	var recentMessages []models.ChatMessage
	h.db.Where("session_id = ?", parsedSessionID).
		Order("timestamp DESC").
		Limit(10). // 取最近的10條訊息作為上下文
		Find(&recentMessages)

	// 構建對話歷史給 AI
	var conversationHistory []dto.Message
	for i := len(recentMessages) - 1; i >= 0; i-- { // 反向排列以保持時間順序
		message := recentMessages[i]
		role := message.Role
		if role == "bot" {
			role = "assistant"
		}
		conversationHistory = append(conversationHistory, dto.Message{
			Role:    role,
			Content: message.Content,
		})
	}

	// 添加當前使用者訊息
	conversationHistory = append(conversationHistory, dto.Message{
		Role:    "user",
		Content: req.Content,
	})

	// 調用 OpenRouter API 使用完整對話歷史
	aiResponse, err := h.callOpenRouterAPIWithHistory(conversationHistory, req.Model)
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
		SessionID: &parsedSessionID,
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
	h.db.Model(&models.ChatSession{}).Where("id = ?", parsedSessionID).
		Updates(map[string]interface{}{
			"last_updated_at": time.Now(),
			"message_count":   gorm.Expr("message_count + ?", 2), // user + bot message
		})

	// 構建回應
	response := dto.ChatMessageResponse{
		ID:        botMessage.ID.String(),
		UserID:    botMessage.UserID.String(),
		SessionID: &sessionID,
		Role:      botMessage.Role,
		Content:   botMessage.Content,
		Timestamp: botMessage.Timestamp,
		Model:     botMessage.Model,
		Tokens:    botMessage.Tokens,
		CreatedAt: botMessage.CreatedAt.Format(time.RFC3339),
	}

	c.JSON(http.StatusOK, vo.SuccessResponse(response, "Message sent successfully"))
}

// callOpenRouterAPIWithHistory 使用對話歷史調用 OpenRouter API
func (h *ChatHandler) callOpenRouterAPIWithHistory(messages []dto.Message, model string) (*dto.OpenRouterResponse, error) {
	// 構建請求
	request := dto.OpenRouterRequest{
		Model:       model,
		Messages:    messages,
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
