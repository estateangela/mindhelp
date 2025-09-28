package dto

import "time"

// CounselorRequest 諮商師請求結構
type CounselorRequest struct {
	Name             string `json:"name" validate:"required"`
	LicenseNumber    string `json:"license_number" validate:"required"`
	Gender           string `json:"gender"`
	Specialties      string `json:"specialties"`
	LanguageSkills   string `json:"language_skills"`
	WorkLocation     string `json:"work_location"`
	WorkUnit         string `json:"work_unit"`
	InstitutionCode  string `json:"institution_code"`
	PsychologySchool string `json:"psychology_school"`
	TreatmentMethods string `json:"treatment_methods"`
}

// CounselorResponse 諮商師回應結構
type CounselorResponse struct {
	ID               string    `json:"id"`
	Name             string    `json:"name"`
	LicenseNumber    string    `json:"license_number"`
	Gender           string    `json:"gender"`
	Specialties      string    `json:"specialties"`
	LanguageSkills   string    `json:"language_skills"`
	WorkLocation     string    `json:"work_location"`
	WorkUnit         string    `json:"work_unit"`
	InstitutionCode  string    `json:"institution_code"`
	PsychologySchool string    `json:"psychology_school"`
	TreatmentMethods string    `json:"treatment_methods"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// CounselorListResponse 諮商師列表回應結構
type CounselorListResponse struct {
	Counselors []CounselorResponse `json:"counselors"`
	Total      int64               `json:"total"`
	Page       int                 `json:"page"`
	PageSize   int                 `json:"page_size"`
}
