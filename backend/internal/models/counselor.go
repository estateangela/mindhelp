package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Counselor 諮商師模型
type Counselor struct {
	ID                uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name              string         `json:"name" gorm:"size:255;not null"`                    // 名字
	LicenseNumber     string         `json:"license_number" gorm:"size:50;uniqueIndex;not null"` // 編號
	Gender            string         `json:"gender" gorm:"size:10"`                            // 性別
	Specialties       string         `json:"specialties" gorm:"type:text"`                     // 專長
	LanguageSkills    string         `json:"language_skills" gorm:"type:text"`                 // 語言專長
	WorkLocation      string         `json:"work_location" gorm:"size:255"`                    // 工作地點
	WorkUnit          string         `json:"work_unit" gorm:"size:255"`                        // 工作單位
	InstitutionCode   string         `json:"institution_code" gorm:"size:50"`                  // 機構代碼
	PsychologySchool  string         `json:"psychology_school" gorm:"size:255"`                // 心理學派
	TreatmentMethods  string         `json:"treatment_methods" gorm:"type:text"`               // 治療方式
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName 指定表名
func (Counselor) TableName() string {
	return "counselors"
}
