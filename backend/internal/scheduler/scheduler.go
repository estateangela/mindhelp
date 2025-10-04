package scheduler

import (
	"context"
	"fmt"
	"log"
	"time"

	"mindhelp-backend/internal/config"
	"mindhelp-backend/internal/database"
	"mindhelp-backend/internal/models"

	"github.com/robfig/cron/v3"
)

// Scheduler 定時任務調度器
type Scheduler struct {
	cron   *cron.Cron
	cfg    *config.Config
	ctx    context.Context
	cancel context.CancelFunc
}

// NotificationMessage 通知訊息結構
type NotificationMessage struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Type      string    `json:"type"` // hourly, weekly
	CreatedAt time.Time `json:"created_at"`
}

// NewScheduler 創建新的調度器
func NewScheduler(cfg *config.Config) *Scheduler {
	// 使用台北時區
	taipeiLocation, err := time.LoadLocation("Asia/Taipei")
	if err != nil {
		log.Printf("Warning: Failed to load Taipei timezone, using UTC: %v", err)
		taipeiLocation = time.UTC
	}

	c := cron.New(
		cron.WithLocation(taipeiLocation),
		cron.WithLogger(cron.VerbosePrintfLogger(log.New(log.Writer(), "Cron: ", log.LstdFlags))),
	)

	ctx, cancel := context.WithCancel(context.Background())

	return &Scheduler{
		cron:   c,
		cfg:    cfg,
		ctx:    ctx,
		cancel: cancel,
	}
}

// Start 啟動定時任務
func (s *Scheduler) Start() error {
	log.Println("Starting notification scheduler...")

	// 每小時發送通知 (在每小時的第0分鐘)
	_, err := s.cron.AddFunc("0 * * * *", s.sendHourlyNotification)
	if err != nil {
		return fmt.Errorf("failed to add hourly cron job: %v", err)
	}

	// 每週一晚上8點發送週報通知
	_, err = s.cron.AddFunc("0 20 * * 1", s.sendWeeklyNotification)
	if err != nil {
		return fmt.Errorf("failed to add weekly cron job: %v", err)
	}

	// 啟動 cron
	s.cron.Start()

	log.Println("Notification scheduler started successfully")
	return nil
}

// Stop 停止定時任務
func (s *Scheduler) Stop() {
	log.Println("Stopping notification scheduler...")
	s.cancel()
	s.cron.Stop()
	log.Println("Notification scheduler stopped")
}

// sendHourlyNotification 發送每小時通知
func (s *Scheduler) sendHourlyNotification() {
	log.Println("Executing hourly notification job...")

	// 隨機選擇一個小時通知訊息
	hourlyMessages := []string{
		"今天心情還好嗎？來和AI 說說話吧 🌿",
		"有些困擾說出口會好一點。來讓 AI 小幫手聽你說說吧 👂",
	}

	message := hourlyMessages[time.Now().Unix()%int64(len(hourlyMessages))]

	// 獲取所有活躍用戶
	var users []models.User
	if err := database.GetDB().Where("deleted_at IS NULL").Find(&users).Error; err != nil {
		log.Printf("Error fetching users for hourly notification: %v", err)
		return
	}

	// 為每個用戶創建通知記錄
	for _, user := range users {
		notification := models.Notification{
			UserID:    user.ID,
			Title:     "心理健康提醒",
			Content:   message,
			Type:      "hourly_reminder",
			IsRead:    false,
			CreatedAt: time.Now(),
		}

		if err := database.GetDB().Create(&notification).Error; err != nil {
			log.Printf("Error creating hourly notification for user %s: %v", user.ID.String(), err)
		}
	}

	log.Printf("Hourly notification sent to %d users", len(users))
}

// sendWeeklyNotification 發送每週通知
func (s *Scheduler) sendWeeklyNotification() {
	log.Println("Executing weekly notification job...")

	// 週報通知訊息
	weeklyMessage := "週一心理健康週報 📊\n\n這週記得照顧好自己的心理健康，如果需要專業協助，記得尋求諮商師或醫療資源的幫助。\n\nMindHelp 團隊關心您 💚"

	// 獲取所有活躍用戶
	var users []models.User
	if err := database.GetDB().Where("deleted_at IS NULL").Find(&users).Error; err != nil {
		log.Printf("Error fetching users for weekly notification: %v", err)
		return
	}

	// 為每個用戶創建通知記錄
	for _, user := range users {
		notification := models.Notification{
			UserID:    user.ID,
			Title:     "週一心理健康週報",
			Content:   weeklyMessage,
			Type:      "weekly_bulletin",
			IsRead:    false,
			CreatedAt: time.Now(),
		}

		if err := database.GetDB().Create(&notification).Error; err != nil {
			log.Printf("Error creating weekly notification for user %s: %v", user.ID.String(), err)
		}
	}

	log.Printf("Weekly notification sent to %d users", len(users))
}

// GetScheduledJobs 獲取已排程的任務資訊
func (s *Scheduler) GetScheduledJobs() []map[string]interface{} {
	entries := s.cron.Entries()
	jobs := make([]map[string]interface{}, len(entries))

	for i, entry := range entries {
		jobs[i] = map[string]interface{}{
			"id":       entry.ID,
			"next_run": entry.Next,
			"schedule": fmt.Sprintf("%v", entry.Schedule),
		}
	}

	return jobs
}

// TriggerHourlyNotification 手動觸發每小時通知
func (s *Scheduler) TriggerHourlyNotification() error {
	log.Println("Manually triggering hourly notification...")
	s.sendHourlyNotification()
	return nil
}

// TriggerWeeklyNotification 手動觸發每週通知
func (s *Scheduler) TriggerWeeklyNotification() error {
	log.Println("Manually triggering weekly notification...")
	s.sendWeeklyNotification()
	return nil
}
