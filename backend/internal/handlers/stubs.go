package handlers

import (
	"net/http"

	"mindhelp-backend/internal/config"
	"mindhelp-backend/internal/vo"

	"github.com/gin-gonic/gin"
)

// 這個檔案包含所有缺失的 handler stubs，以滿足路由配置需求

// UserHandler 使用者處理器 stub
type UserHandler struct {
	cfg *config.Config
}

func NewUserHandler(cfg *config.Config) *UserHandler {
	return &UserHandler{cfg: cfg}
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, vo.NewErrorResponse(
		"not_implemented",
		"此功能尚未實現",
		"NOT_IMPLEMENTED",
		nil,
		c.Request.URL.Path,
	))
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, vo.NewErrorResponse(
		"not_implemented",
		"此功能尚未實現",
		"NOT_IMPLEMENTED",
		nil,
		c.Request.URL.Path,
	))
}

func (h *UserHandler) DeleteAccount(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, vo.NewErrorResponse(
		"not_implemented",
		"此功能尚未實現",
		"NOT_IMPLEMENTED",
		nil,
		c.Request.URL.Path,
	))
}

func (h *UserHandler) ChangePassword(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, vo.NewErrorResponse(
		"not_implemented",
		"此功能尚未實現",
		"NOT_IMPLEMENTED",
		nil,
		c.Request.URL.Path,
	))
}

func (h *UserHandler) GetStats(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, vo.NewErrorResponse(
		"not_implemented",
		"此功能尚未實現",
		"NOT_IMPLEMENTED",
		nil,
		c.Request.URL.Path,
	))
}

// QuizHandler 測驗處理器 stub
type QuizHandler struct{}

func NewQuizHandler() *QuizHandler {
	return &QuizHandler{}
}

func (h *QuizHandler) SubmitQuiz(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, vo.NewErrorResponse(
		"not_implemented",
		"此功能尚未實現",
		"NOT_IMPLEMENTED",
		nil,
		c.Request.URL.Path,
	))
}

func (h *QuizHandler) GetQuizHistory(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, vo.NewErrorResponse(
		"not_implemented",
		"此功能尚未實現",
		"NOT_IMPLEMENTED",
		nil,
		c.Request.URL.Path,
	))
}

func (h *QuizHandler) GetQuizzes(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, vo.NewErrorResponse(
		"not_implemented",
		"此功能尚未實現",
		"NOT_IMPLEMENTED",
		nil,
		c.Request.URL.Path,
	))
}

func (h *QuizHandler) GetQuiz(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, vo.NewErrorResponse(
		"not_implemented",
		"此功能尚未實現",
		"NOT_IMPLEMENTED",
		nil,
		c.Request.URL.Path,
	))
}

// BookmarkHandler 收藏處理器 stub
type BookmarkHandler struct{}

func NewBookmarkHandler() *BookmarkHandler {
	return &BookmarkHandler{}
}

func (h *BookmarkHandler) GetArticleBookmarks(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, vo.NewErrorResponse(
		"not_implemented",
		"此功能尚未實現",
		"NOT_IMPLEMENTED",
		nil,
		c.Request.URL.Path,
	))
}

func (h *BookmarkHandler) GetLocationBookmarks(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, vo.NewErrorResponse(
		"not_implemented",
		"此功能尚未實現",
		"NOT_IMPLEMENTED",
		nil,
		c.Request.URL.Path,
	))
}

func (h *BookmarkHandler) BookmarkResource(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, vo.NewErrorResponse(
		"not_implemented",
		"此功能尚未實現",
		"NOT_IMPLEMENTED",
		nil,
		c.Request.URL.Path,
	))
}

func (h *BookmarkHandler) UnbookmarkResource(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, vo.NewErrorResponse(
		"not_implemented",
		"此功能尚未實現",
		"NOT_IMPLEMENTED",
		nil,
		c.Request.URL.Path,
	))
}

// ArticleHandler 文章處理器 stub
type ArticleHandler struct{}

func NewArticleHandler() *ArticleHandler {
	return &ArticleHandler{}
}

func (h *ArticleHandler) BookmarkArticle(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, vo.NewErrorResponse(
		"not_implemented",
		"此功能尚未實現",
		"NOT_IMPLEMENTED",
		nil,
		c.Request.URL.Path,
	))
}

func (h *ArticleHandler) UnbookmarkArticle(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, vo.NewErrorResponse(
		"not_implemented",
		"此功能尚未實現",
		"NOT_IMPLEMENTED",
		nil,
		c.Request.URL.Path,
	))
}

func (h *ArticleHandler) GetArticles(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, vo.NewErrorResponse(
		"not_implemented",
		"此功能尚未實現",
		"NOT_IMPLEMENTED",
		nil,
		c.Request.URL.Path,
	))
}

func (h *ArticleHandler) GetArticle(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, vo.NewErrorResponse(
		"not_implemented",
		"此功能尚未實現",
		"NOT_IMPLEMENTED",
		nil,
		c.Request.URL.Path,
	))
}

// ReviewHandler 評論處理器 stub
type ReviewHandler struct{}

func NewReviewHandler() *ReviewHandler {
	return &ReviewHandler{}
}

func (h *ReviewHandler) UpdateReview(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, vo.NewErrorResponse(
		"not_implemented",
		"此功能尚未實現",
		"NOT_IMPLEMENTED",
		nil,
		c.Request.URL.Path,
	))
}

func (h *ReviewHandler) DeleteReview(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, vo.NewErrorResponse(
		"not_implemented",
		"此功能尚未實現",
		"NOT_IMPLEMENTED",
		nil,
		c.Request.URL.Path,
	))
}

func (h *ReviewHandler) CreateReview(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, vo.NewErrorResponse(
		"not_implemented",
		"此功能尚未實現",
		"NOT_IMPLEMENTED",
		nil,
		c.Request.URL.Path,
	))
}

func (h *ReviewHandler) ReportContent(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, vo.NewErrorResponse(
		"not_implemented",
		"此功能尚未實現",
		"NOT_IMPLEMENTED",
		nil,
		c.Request.URL.Path,
	))
}

func (h *ReviewHandler) GetResourceReviews(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, vo.NewErrorResponse(
		"not_implemented",
		"此功能尚未實現",
		"NOT_IMPLEMENTED",
		nil,
		c.Request.URL.Path,
	))
}

// NotificationHandler 通知處理器 stub
type NotificationHandler struct{}

func NewNotificationHandler() *NotificationHandler {
	return &NotificationHandler{}
}

func (h *NotificationHandler) GetNotifications(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, vo.NewErrorResponse(
		"not_implemented",
		"此功能尚未實現",
		"NOT_IMPLEMENTED",
		nil,
		c.Request.URL.Path,
	))
}

func (h *NotificationHandler) MarkAsRead(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, vo.NewErrorResponse(
		"not_implemented",
		"此功能尚未實現",
		"NOT_IMPLEMENTED",
		nil,
		c.Request.URL.Path,
	))
}

func (h *NotificationHandler) GetNotificationSettings(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, vo.NewErrorResponse(
		"not_implemented",
		"此功能尚未實現",
		"NOT_IMPLEMENTED",
		nil,
		c.Request.URL.Path,
	))
}

func (h *NotificationHandler) UpdateNotificationSettings(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, vo.NewErrorResponse(
		"not_implemented",
		"此功能尚未實現",
		"NOT_IMPLEMENTED",
		nil,
		c.Request.URL.Path,
	))
}

func (h *NotificationHandler) UpdatePushToken(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, vo.NewErrorResponse(
		"not_implemented",
		"此功能尚未實現",
		"NOT_IMPLEMENTED",
		nil,
		c.Request.URL.Path,
	))
}

// ShareHandler 分享處理器 stub
type ShareHandler struct {
	cfg *config.Config
}

func NewShareHandler(cfg *config.Config) *ShareHandler {
	return &ShareHandler{cfg: cfg}
}

func (h *ShareHandler) CreateShare(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, vo.NewErrorResponse(
		"not_implemented",
		"此功能尚未實現",
		"NOT_IMPLEMENTED",
		nil,
		c.Request.URL.Path,
	))
}

func (h *ShareHandler) GetUserShares(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, vo.NewErrorResponse(
		"not_implemented",
		"此功能尚未實現",
		"NOT_IMPLEMENTED",
		nil,
		c.Request.URL.Path,
	))
}

func (h *ShareHandler) GetShare(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, vo.NewErrorResponse(
		"not_implemented",
		"此功能尚未實現",
		"NOT_IMPLEMENTED",
		nil,
		c.Request.URL.Path,
	))
}

func (h *ShareHandler) GetShareStats(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, vo.NewErrorResponse(
		"not_implemented",
		"此功能尚未實現",
		"NOT_IMPLEMENTED",
		nil,
		c.Request.URL.Path,
	))
}

// ConfigHandler 配置處理器 stub
type ConfigHandler struct{}

func NewConfigHandler() *ConfigHandler {
	return &ConfigHandler{}
}

func (h *ConfigHandler) GetConfig(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, vo.NewErrorResponse(
		"not_implemented",
		"此功能尚未實現",
		"NOT_IMPLEMENTED",
		nil,
		c.Request.URL.Path,
	))
}

// GoogleMapsHandler Google Maps 處理器 stub
type GoogleMapsHandler struct {
	cfg *config.Config
}

func NewGoogleMapsHandler(cfg *config.Config) *GoogleMapsHandler {
	return &GoogleMapsHandler{cfg: cfg}
}

func (h *GoogleMapsHandler) Geocode(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, vo.NewErrorResponse(
		"not_implemented",
		"此功能尚未實現",
		"NOT_IMPLEMENTED",
		nil,
		c.Request.URL.Path,
	))
}

func (h *GoogleMapsHandler) ReverseGeocode(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, vo.NewErrorResponse(
		"not_implemented",
		"此功能尚未實現",
		"NOT_IMPLEMENTED",
		nil,
		c.Request.URL.Path,
	))
}

func (h *GoogleMapsHandler) SearchPlaces(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, vo.NewErrorResponse(
		"not_implemented",
		"此功能尚未實現",
		"NOT_IMPLEMENTED",
		nil,
		c.Request.URL.Path,
	))
}

func (h *GoogleMapsHandler) GetDirections(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, vo.NewErrorResponse(
		"not_implemented",
		"此功能尚未實現",
		"NOT_IMPLEMENTED",
		nil,
		c.Request.URL.Path,
	))
}

func (h *GoogleMapsHandler) GetDistanceMatrix(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, vo.NewErrorResponse(
		"not_implemented",
		"此功能尚未實現",
		"NOT_IMPLEMENTED",
		nil,
		c.Request.URL.Path,
	))
}

func (h *GoogleMapsHandler) GetNearbyMentalHealthServices(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, vo.NewErrorResponse(
		"not_implemented",
		"此功能尚未實現",
		"NOT_IMPLEMENTED",
		nil,
		c.Request.URL.Path,
	))
}

func (h *GoogleMapsHandler) BatchGeocode(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, vo.NewErrorResponse(
		"not_implemented",
		"此功能尚未實現",
		"NOT_IMPLEMENTED",
		nil,
		c.Request.URL.Path,
	))
}

func (h *GoogleMapsHandler) GetAPIUsageStats(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, vo.NewErrorResponse(
		"not_implemented",
		"此功能尚未實現",
		"NOT_IMPLEMENTED",
		nil,
		c.Request.URL.Path,
	))
}

func (h *GoogleMapsHandler) ClearCache(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, vo.NewErrorResponse(
		"not_implemented",
		"此功能尚未實現",
		"NOT_IMPLEMENTED",
		nil,
		c.Request.URL.Path,
	))
}
