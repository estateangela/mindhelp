package dto

// AppConfigResponse 應用程式配置回應
type AppConfigResponse struct {
	Features      FeatureConfig              `json:"features"`
	Filters       FilterConfig               `json:"filters"`
	APIEndpoints  map[string]string          `json:"api_endpoints,omitempty"`
	SupportInfo   SupportInfo                `json:"support_info,omitempty"`
}

// FeatureConfig 功能配置
type FeatureConfig struct {
	EnableReviews           bool `json:"enableReviews"`
	EnableTherapistProfiles bool `json:"enableTherapistProfiles"`
	EnableGroupChat         bool `json:"enableGroupChat"`
	EnableVideoConsult      bool `json:"enableVideoConsult"`
	EnableSharing           bool `json:"enableSharing"`
}

// FilterConfig 篩選配置
type FilterConfig struct {
	ResourceTypes []FilterOption `json:"resourceTypes"`
	Specialties   []FilterOption `json:"specialties"`
	QuizCategories []FilterOption `json:"quizCategories,omitempty"`
}

// FilterOption 篩選選項
type FilterOption struct {
	Key         string `json:"key"`
	DisplayName string `json:"displayName"`
	Description string `json:"description,omitempty"`
}

// SupportInfo 支援資訊
type SupportInfo struct {
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	Website     string `json:"website"`
	WorkingHours string `json:"working_hours"`
}

// ConfigItem 配置項目 (內部使用)
type ConfigItem struct {
	Key         string `json:"key"`
	Value       string `json:"value"`
	Description string `json:"description"`
	IsActive    bool   `json:"is_active"`
}
