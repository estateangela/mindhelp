package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Quiz 心�?測�?資�?模�?
type Quiz struct {
	ID          uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Title       string         `json:"title" gorm:"size:200;not null"`
	Description string         `json:"description" gorm:"type:text"`
	Category    string         `json:"category" gorm:"size:50"` // anxiety, depression, etc.
	IsActive    bool           `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// 關聯
	Questions   []QuizQuestion   `json:"questions,omitempty" gorm:"foreignKey:QuizID"`
	Submissions []QuizSubmission `json:"submissions,omitempty" gorm:"foreignKey:QuizID"`
}

// QuizQuestion 測�?題目資�?模�?
type QuizQuestion struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	QuizID    uuid.UUID      `json:"quiz_id" gorm:"type:uuid;not null;index"`
	Question  string         `json:"question" gorm:"type:text;not null"`
	Options   string         `json:"options" gorm:"type:text"` // JSON array stored as text
	OrderNum  int            `json:"order_num" gorm:"not null"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// 關聯
	Quiz Quiz `json:"quiz,omitempty" gorm:"foreignKey:QuizID"`
}

// QuizSubmission 測�??�交資�?模�?
type QuizSubmission struct {
	ID          uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID      uuid.UUID      `json:"user_id" gorm:"type:uuid;not null;index"`
	QuizID      uuid.UUID      `json:"quiz_id" gorm:"type:uuid;not null;index"`
	Answers     string         `json:"answers" gorm:"type:text;not null"` // JSON object stored as text
	Score       int            `json:"score" gorm:"not null"`
	Result      string         `json:"result" gorm:"type:text"`
	CompletedAt time.Time      `json:"completed_at" gorm:"not null"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// 關聯
	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Quiz Quiz `json:"quiz,omitempty" gorm:"foreignKey:QuizID"`
}

// TableName 指定資料表名稱�?�?
func (Quiz) TableName() string {
	return "quizzes"
}

func (QuizQuestion) TableName() string {
	return "quiz_questions"
}

func (QuizSubmission) TableName() string {
	return "quiz_submissions"
}

// BeforeCreate hooks
func (q *Quiz) BeforeCreate(tx *gorm.DB) error {
	if q.ID == uuid.Nil {
		q.ID = uuid.New()
	}
	return nil
}

func (qq *QuizQuestion) BeforeCreate(tx *gorm.DB) error {
	if qq.ID == uuid.Nil {
		qq.ID = uuid.New()
	}
	return nil
}

func (qs *QuizSubmission) BeforeCreate(tx *gorm.DB) error {
	if qs.ID == uuid.Nil {
		qs.ID = uuid.New()
	}
	return nil
}


