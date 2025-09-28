package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"mindhelp-backend/internal/database"
	"mindhelp-backend/internal/dto"
	"mindhelp-backend/internal/middleware"
	"mindhelp-backend/internal/models"
	"mindhelp-backend/internal/vo"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// QuizHandler 測驗處理器
type QuizHandler struct {
	db *gorm.DB
}

// NewQuizHandler 創建新的測驗處理器
func NewQuizHandler() *QuizHandler {
	return &QuizHandler{
		db: database.GetDB(),
	}
}

// GetQuizzes 獲取測驗列表
// @Summary 獲取測驗列表
// @Description 獲取可用的心理測驗列表
// @Tags quiz
// @Accept json
// @Produce json
// @Param category query string false "測驗類別"
// @Param page query int false "頁碼" default(1)
// @Param limit query int false "每頁數量" default(10)
// @Success 200 {object} vo.Response{data=dto.QuizListResponse}
// @Failure 400 {object} vo.ErrorResponse
// @Router /quizzes [get]
func (h *QuizHandler) GetQuizzes(c *gin.Context) {
	// 解析查詢參數
	category := strings.TrimSpace(c.Query("category"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	// 參數驗證
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 50 {
		limit = 10
	}

	offset := (page - 1) * limit

	// 構建查詢
	dbQuery := h.db.Model(&models.Quiz{}).Where("is_active = ?", true)

	// 添加類別篩選
	if category != "" {
		dbQuery = dbQuery.Where("category = ?", category)
	}

	// 獲取總數
	var total int64
	if err := dbQuery.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, vo.NewErrorResponse(
			"internal_error",
			"Failed to count quizzes",
			"INTERNAL_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	// 獲取測驗列表
	var quizzes []models.Quiz
	if err := dbQuery.Order("created_at DESC").Offset(offset).Limit(limit).Find(&quizzes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, vo.NewErrorResponse(
			"internal_error",
			"Failed to get quizzes",
			"INTERNAL_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	// 轉換為 DTO
	var quizResponses []dto.QuizResponse
	for _, quiz := range quizzes {
		response := dto.QuizResponse{
			ID:          quiz.ID.String(),
			Title:       quiz.Title,
			Description: quiz.Description,
			Category:    quiz.Category,
			IsActive:    quiz.IsActive,
			CreatedAt:   quiz.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
		quizResponses = append(quizResponses, response)
	}

	// 構建分頁回應
	totalPages := int((total + int64(limit) - 1) / int64(limit))
	response := dto.QuizListResponse{
		Quizzes:    quizResponses,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
		HasMore:    page < totalPages,
	}

	c.JSON(http.StatusOK, vo.SuccessResponse(response, "Quizzes retrieved successfully"))
}

// GetQuiz 獲取測驗詳情（包含題目）
// @Summary 獲取測驗詳情
// @Description 獲取測驗的完整內容和題目
// @Tags quiz
// @Accept json
// @Produce json
// @Param id path string true "測驗ID"
// @Success 200 {object} vo.Response{data=dto.QuizResponse}
// @Failure 400 {object} vo.ErrorResponse
// @Failure 404 {object} vo.ErrorResponse
// @Router /quizzes/{id} [get]
func (h *QuizHandler) GetQuiz(c *gin.Context) {
	quizID := c.Param("id")
	if quizID == "" {
		c.JSON(http.StatusBadRequest, vo.NewErrorResponse(
			"bad_request",
			"Quiz ID is required",
			"VALIDATION_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	// 驗證 UUID 格式
	parsedID, err := uuid.Parse(quizID)
	if err != nil {
		c.JSON(http.StatusBadRequest, vo.NewErrorResponse(
			"bad_request",
			"Invalid quiz ID format",
			"VALIDATION_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	// 查找測驗
	var quiz models.Quiz
	if err := h.db.Where("id = ? AND is_active = ?", parsedID, true).First(&quiz).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, vo.NewErrorResponse(
				"not_found",
				"Quiz not found",
				"NOT_FOUND",
				nil,
				c.Request.URL.Path,
			))
			return
		}
		c.JSON(http.StatusInternalServerError, vo.NewErrorResponse(
			"internal_error",
			"Failed to get quiz",
			"INTERNAL_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	// 獲取測驗題目
	var questions []models.QuizQuestion
	if err := h.db.Where("quiz_id = ?", parsedID).Order("order_num ASC").Find(&questions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, vo.NewErrorResponse(
			"internal_error",
			"Failed to get quiz questions",
			"INTERNAL_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	// 轉換題目為 DTO
	var questionResponses []dto.QuizQuestionResponse
	for _, question := range questions {
		// 解析選項
		var options []string
		if question.Options != "" {
			json.Unmarshal([]byte(question.Options), &options)
		}

		response := dto.QuizQuestionResponse{
			ID:       question.ID.String(),
			Question: question.Question,
			Options:  options,
			OrderNum: question.OrderNum,
		}
		questionResponses = append(questionResponses, response)
	}

	// 構建回應
	response := dto.QuizResponse{
		ID:          quiz.ID.String(),
		Title:       quiz.Title,
		Description: quiz.Description,
		Category:    quiz.Category,
		Questions:   questionResponses,
		IsActive:    quiz.IsActive,
		CreatedAt:   quiz.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	c.JSON(http.StatusOK, vo.SuccessResponse(response, "Quiz retrieved successfully"))
}

// SubmitQuiz 提交測驗答案
// @Summary 提交測驗答案
// @Description 提交測驗答案並獲取結果
// @Tags quiz
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "測驗ID"
// @Param request body dto.QuizSubmissionRequest true "測驗答案"
// @Success 201 {object} vo.Response{data=dto.QuizSubmissionResponse}
// @Failure 400 {object} vo.ErrorResponse
// @Failure 401 {object} vo.ErrorResponse
// @Failure 404 {object} vo.ErrorResponse
// @Router /quizzes/{id}/submit [post]
func (h *QuizHandler) SubmitQuiz(c *gin.Context) {
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

	quizID := c.Param("id")
	parsedID, err := uuid.Parse(quizID)
	if err != nil {
		c.JSON(http.StatusBadRequest, vo.NewErrorResponse(
			"bad_request",
			"Invalid quiz ID format",
			"VALIDATION_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	var req dto.QuizSubmissionRequest
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

	// 檢查測驗是否存在且啟用
	var quiz models.Quiz
	if err := h.db.Where("id = ? AND is_active = ?", parsedID, true).First(&quiz).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, vo.NewErrorResponse(
				"not_found",
				"Quiz not found",
				"NOT_FOUND",
				nil,
				c.Request.URL.Path,
			))
			return
		}
		c.JSON(http.StatusInternalServerError, vo.NewErrorResponse(
			"internal_error",
			"Failed to get quiz",
			"INTERNAL_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	// 計算分數（簡化版本，實際可能需要更複雜的評分邏輯）
	score := h.calculateScore(req.Answers, parsedID)
	result := h.generateResult(quiz.Category, score)

	// 序列化答案
	answersJSON, _ := json.Marshal(req.Answers)

	// 創建提交記錄
	submission := models.QuizSubmission{
		UserID:      uuid.MustParse(userID),
		QuizID:      parsedID,
		Answers:     string(answersJSON),
		Score:       score,
		Result:      result,
		CompletedAt: time.Now(),
	}

	if err := h.db.Create(&submission).Error; err != nil {
		c.JSON(http.StatusInternalServerError, vo.NewErrorResponse(
			"internal_error",
			"Failed to save quiz submission",
			"INTERNAL_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	// 構建回應
	response := dto.QuizSubmissionResponse{
		ID:          submission.ID.String(),
		QuizTitle:   quiz.Title,
		Score:       submission.Score,
		Result:      submission.Result,
		CompletedAt: submission.CompletedAt.Format("2006-01-02T15:04:05Z07:00"),
		CreatedAt:   submission.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	c.JSON(http.StatusCreated, vo.SuccessResponse(response, "Quiz submitted successfully"))
}

// GetQuizHistory 獲取使用者測驗歷史
// @Summary 獲取測驗歷史
// @Description 獲取使用者的測驗完成歷史
// @Tags quiz
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "頁碼" default(1)
// @Param limit query int false "每頁數量" default(10)
// @Success 200 {object} vo.Response{data=dto.QuizHistoryResponse}
// @Failure 401 {object} vo.ErrorResponse
// @Router /users/me/quiz_history [get]
func (h *QuizHandler) GetQuizHistory(c *gin.Context) {
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
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	// 參數驗證
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 50 {
		limit = 10
	}

	offset := (page - 1) * limit

	// 獲取總數
	var total int64
	if err := h.db.Model(&models.QuizSubmission{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, vo.NewErrorResponse(
			"internal_error",
			"Failed to count quiz submissions",
			"INTERNAL_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	// 獲取提交記錄
	var submissions []models.QuizSubmission
	if err := h.db.Preload("Quiz").Where("user_id = ?", userID).
		Order("completed_at DESC").
		Offset(offset).Limit(limit).
		Find(&submissions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, vo.NewErrorResponse(
			"internal_error",
			"Failed to get quiz submissions",
			"INTERNAL_ERROR",
			nil,
			c.Request.URL.Path,
		))
		return
	}

	// 轉換為 DTO
	var submissionResponses []dto.QuizSubmissionResponse
	for _, submission := range submissions {
		response := dto.QuizSubmissionResponse{
			ID:          submission.ID.String(),
			QuizTitle:   submission.Quiz.Title,
			Score:       submission.Score,
			Result:      submission.Result,
			CompletedAt: submission.CompletedAt.Format("2006-01-02T15:04:05Z07:00"),
			CreatedAt:   submission.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
		submissionResponses = append(submissionResponses, response)
	}

	// 構建分頁回應
	totalPages := int((total + int64(limit) - 1) / int64(limit))
	response := dto.QuizHistoryResponse{
		Submissions: submissionResponses,
		Total:       total,
		Page:        page,
		Limit:       limit,
		TotalPages:  totalPages,
		HasMore:     page < totalPages,
	}

	c.JSON(http.StatusOK, vo.SuccessResponse(response, "Quiz history retrieved successfully"))
}

// calculateScore 計算測驗分數
func (h *QuizHandler) calculateScore(answers map[string]int, quizID uuid.UUID) int {
	// 這是簡化版本的計分邏輯
	// 實際應用中可能需要根據不同測驗類型實現不同的計分規則
	
	// 獲取題目數量
	var questionCount int64
	h.db.Model(&models.QuizQuestion{}).Where("quiz_id = ?", quizID).Count(&questionCount)
	
	if questionCount == 0 {
		return 0
	}

	// 簡單計分：每題等權重
	totalAnswers := len(answers)
	if totalAnswers == 0 {
		return 0
	}

	// 根據答案計分（這裡簡化為答案值的總和）
	score := 0
	for _, answer := range answers {
		score += answer
	}

	// 標準化分數到合理範圍 (0-27 for GAD-7, 0-27 for PHQ-9, etc.)
	maxPossibleScore := int(questionCount) * 3 // 假設最高選項值為3
	if maxPossibleScore > 0 {
		score = (score * 27) / maxPossibleScore
	}

	return score
}

// generateResult 根據分數和測驗類型生成結果解釋
func (h *QuizHandler) generateResult(category string, score int) string {
	switch category {
	case "anxiety":
		if score <= 4 {
			return "您的焦慮程度在正常範圍內。"
		} else if score <= 9 {
			return "您可能有輕度焦慮症狀。建議留意情緒變化。"
		} else if score <= 14 {
			return "您可能有中度焦慮症狀。建議尋求專業諮詢。"
		} else {
			return "您可能有嚴重焦慮症狀。強烈建議尋求專業醫療協助。"
		}
	case "depression":
		if score <= 4 {
			return "您的情緒狀態在正常範圍內。"
		} else if score <= 9 {
			return "您可能有輕度憂鬱症狀。建議多關注自己的情緒健康。"
		} else if score <= 14 {
			return "您可能有中度憂鬱症狀。建議尋求專業諮詢。"
		} else {
			return "您可能有嚴重憂鬱症狀。強烈建議尋求專業醫療協助。"
		}
	case "stress":
		if score <= 9 {
			return "您的壓力水平較低，心理狀態良好。"
		} else if score <= 18 {
			return "您有中等程度的壓力，建議學習壓力管理技巧。"
		} else {
			return "您的壓力水平較高，建議尋求專業協助學習應對策略。"
		}
	default:
		return "感謝您完成測驗。建議將結果與專業人員討論以獲得更詳細的解釋。"
	}
}
