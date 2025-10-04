package handlers

import (
	"net/http"
	"strings"

	"mindhelp-backend/internal/database"
	"mindhelp-backend/internal/models"
	"mindhelp-backend/internal/vo"

	"github.com/gin-gonic/gin"
)

// AdminHandler 管理員相關處理器
type AdminHandler struct{}

// NewAdminHandler 創建新的管理員處理器
func NewAdminHandler() *AdminHandler {
	return &AdminHandler{}
}

// SeedDatabase 資料庫種子資料
// @Summary 插入種子資料
// @Description 插入範例的諮商師、諮商所和推薦醫師資料
// @Tags admin
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} vo.ErrorResponse
// @Router /admin/seed-database [post]
func (h *AdminHandler) SeedDatabase(c *gin.Context) {
	// 獲取資料庫連接
	db, err := database.GetDBSafely()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, vo.NewErrorResponse(
			"database_unavailable",
			"Database service is currently unavailable",
			"SERVICE_UNAVAILABLE",
			[]string{err.Error()},
			c.Request.URL.Path,
		))
		return
	}

	// 範例諮商師資料
	counselors := []models.Counselor{
		{
			Name:             "王心理師",
			LicenseNumber:    "PSY001",
			Gender:           "女",
			Specialties:      "焦慮症、憂鬱症、創傷治療",
			LanguageSkills:   "中文、英文",
			WorkLocation:     "台北市信義區信義路五段7號101大樓",
			WorkUnit:         "台北心理診所",
			InstitutionCode:  "TP001",
			PsychologySchool: "台灣大學心理系",
			TreatmentMethods: "認知行為療法、正念治療",
		},
		{
			Name:             "李諮商師",
			LicenseNumber:    "PSY002",
			Gender:           "男",
			Specialties:      "家庭治療、伴侶諮商、青少年輔導",
			LanguageSkills:   "中文、台語",
			WorkLocation:     "台北市大安區復興南路一段390號",
			WorkUnit:         "大安諮商中心",
			InstitutionCode:  "TP002",
			PsychologySchool: "政治大學心理系",
			TreatmentMethods: "系統性家族治療、敘事治療",
		},
		{
			Name:             "陳心理師",
			LicenseNumber:    "PSY003",
			Gender:           "女",
			Specialties:      "兒童心理、學習障礙、注意力不足",
			LanguageSkills:   "中文、英文、日文",
			WorkLocation:     "台北市中山區南京東路二段125號",
			WorkUnit:         "中山兒童心理中心",
			InstitutionCode:  "TP003",
			PsychologySchool: "師範大學特教系",
			TreatmentMethods: "遊戲治療、藝術治療",
		},
	}

	// 範例諮商所資料
	centers := []models.CounselingCenter{
		{
			Name:             "台北心理健康中心",
			Address:          "台北市中正區中山南路1號2樓",
			Phone:            "02-2311-1234",
			OnlineCounseling: true,
		},
		{
			Name:             "信義諮商所",
			Address:          "台北市信義區信義路四段1號8樓",
			Phone:            "02-2722-5678",
			OnlineCounseling: false,
		},
		{
			Name:             "大安心理診所",
			Address:          "台北市大安區敦化南路二段216號3樓",
			Phone:            "02-2733-9012",
			OnlineCounseling: true,
		},
		{
			Name:             "松山諮商中心",
			Address:          "台北市松山區八德路四段138號5樓",
			Phone:            "02-2570-3456",
			OnlineCounseling: true,
		},
	}

	// 範例推薦醫師資料
	doctors := []models.RecommendedDoctor{
		{
			Name:            "張精神科醫師",
			Description:     "張精神科醫師 - 台大醫院精神科主治醫師，專精於憂鬱症、焦慮症治療，具有20年豐富經驗。位於台北市中正區中山南路7號台大醫院。",
			ExperienceCount: 20,
		},
		{
			Name:            "林心理師",
			Description:     "林心理師 - 榮總精神科心理師，專業於創傷後壓力症候群、恐慌症治療。台北市北投區石牌路二段201號榮民總醫院。",
			ExperienceCount: 15,
		},
		{
			Name:            "黃醫師",
			Description:     "黃醫師 - 馬偕醫院精神科，擅長家庭治療、青少年心理問題。台北市中山區中山北路二段92號馬偕紀念醫院。",
			ExperienceCount: 12,
		},
		{
			Name:            "吳心理師",
			Description:     "吳心理師 - 新光醫院心理師，專精於強迫症、恐懼症治療，提供個別諮商與團體治療。台北市士林區文昌路95號新光醫院。",
			ExperienceCount: 8,
		},
	}

	var createdCount int
	var updatedCount int
	var errors []string

	// 插入諮商師
	for _, counselor := range counselors {
		var existing models.Counselor
		result := db.Where("license_number = ?", counselor.LicenseNumber).First(&existing)
		if result.Error != nil {
			// 不存在，創建新記錄
			if err := db.Create(&counselor).Error; err != nil {
				errors = append(errors, "Failed to create counselor "+counselor.Name+": "+err.Error())
			} else {
				createdCount++
			}
		} else {
			// 存在，更新記錄
			counselor.ID = existing.ID
			if err := db.Save(&counselor).Error; err != nil {
				errors = append(errors, "Failed to update counselor "+counselor.Name+": "+err.Error())
			} else {
				updatedCount++
			}
		}
	}

	// 插入諮商所
	for _, center := range centers {
		var existing models.CounselingCenter
		result := db.Where("name = ?", center.Name).First(&existing)
		if result.Error != nil {
			// 不存在，創建新記錄
			if err := db.Create(&center).Error; err != nil {
				errors = append(errors, "Failed to create counseling center "+center.Name+": "+err.Error())
			} else {
				createdCount++
			}
		} else {
			// 存在，更新記錄
			center.ID = existing.ID
			if err := db.Save(&center).Error; err != nil {
				errors = append(errors, "Failed to update counseling center "+center.Name+": "+err.Error())
			} else {
				updatedCount++
			}
		}
	}

	// 插入推薦醫師
	for _, doctor := range doctors {
		var existing models.RecommendedDoctor
		result := db.Where("name = ?", doctor.Name).First(&existing)
		if result.Error != nil {
			// 不存在，創建新記錄
			if err := db.Create(&doctor).Error; err != nil {
				errors = append(errors, "Failed to create recommended doctor "+doctor.Name+": "+err.Error())
			} else {
				createdCount++
			}
		} else {
			// 存在，更新記錄
			doctor.ID = existing.ID
			if err := db.Save(&doctor).Error; err != nil {
				errors = append(errors, "Failed to update recommended doctor "+doctor.Name+": "+err.Error())
			} else {
				updatedCount++
			}
		}
	}

	if len(errors) > 0 {
		c.JSON(http.StatusInternalServerError, vo.ErrorResponse{
			Code:    "SEED_PARTIAL_FAILURE",
			Message: "部分資料插入失敗",
			Error:   "Errors: " + strings.Join(errors, "; "),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"created_count":       createdCount,
			"updated_count":       updatedCount,
			"total_processed":     createdCount + updatedCount,
			"counselors":          len(counselors),
			"counseling_centers":  len(centers),
			"recommended_doctors": len(doctors),
		},
		"message": "資料庫種子資料插入成功",
	})
}

// GetDatabaseStats 獲取資料庫統計
// @Summary 獲取資料庫統計
// @Description 獲取資料庫中各表的記錄數量
// @Tags admin
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} vo.ErrorResponse
// @Router /admin/database-stats [get]
func (h *AdminHandler) GetDatabaseStats(c *gin.Context) {
	// 獲取資料庫連接
	db, err := database.GetDBSafely()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, vo.NewErrorResponse(
			"database_unavailable",
			"Database service is currently unavailable",
			"SERVICE_UNAVAILABLE",
			[]string{err.Error()},
			c.Request.URL.Path,
		))
		return
	}

	var counselorCount int64
	var centerCount int64
	var doctorCount int64
	var userCount int64

	// 計算各表記錄數
	if err := db.Model(&models.Counselor{}).Count(&counselorCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, vo.ErrorResponse{
			Code:    "DB_ERROR",
			Message: "無法取得諮商師數量",
			Error:   "Counselor count error: " + err.Error(),
		})
		return
	}
	if err := db.Model(&models.CounselingCenter{}).Count(&centerCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, vo.ErrorResponse{
			Code:    "DB_ERROR",
			Message: "無法取得諮商中心數量",
			Error:   "CounselingCenter count error: " + err.Error(),
		})
		return
	}
	if err := db.Model(&models.RecommendedDoctor{}).Count(&doctorCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, vo.ErrorResponse{
			Code:    "DB_ERROR",
			Message: "無法取得推薦醫師數量",
			Error:   "RecommendedDoctor count error: " + err.Error(),
		})
		return
	}
	if err := db.Model(&models.User{}).Count(&userCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, vo.ErrorResponse{
			Code:    "DB_ERROR",
			Message: "無法取得使用者數量",
			Error:   "User count error: " + err.Error(),
		})
		return
	}

	// 計算有地址的記錄數
	var counselorWithLocationCount int64
	var centerWithAddressCount int64
	var doctorWithDescriptionCount int64

	// 使用 db（注意大小寫，確保正確存取 DB 實例）並加上錯誤處理，避免未定義變數與潛在 SQL 錯誤
	if err := db.Model(&models.Counselor{}).
		Where("work_location IS NOT NULL AND work_location != ''").
		Count(&counselorWithLocationCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, vo.ErrorResponse{
			Code:    "DB_ERROR",
			Message: "無法取得有工作地點的諮商師數量",
			Error:   "Counselor with location count error: " + err.Error(),
		})
		return
	}
	if err := db.Model(&models.CounselingCenter{}).
		Where("address IS NOT NULL AND address != ''").
		Count(&centerWithAddressCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, vo.ErrorResponse{
			Code:    "DB_ERROR",
			Message: "無法取得有地址的諮商中心數量",
			Error:   "CounselingCenter with address count error: " + err.Error(),
		})
		return
	}
	if err := db.Model(&models.RecommendedDoctor{}).
		Where("description IS NOT NULL AND description != ''").
		Count(&doctorWithDescriptionCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, vo.ErrorResponse{
			Code:    "DB_ERROR",
			Message: "無法取得有描述的推薦醫師數量",
			Error:   "RecommendedDoctor with description count error: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"total_records": gin.H{
				"counselors":          counselorCount,
				"counseling_centers":  centerCount,
				"recommended_doctors": doctorCount,
				"users":               userCount,
			},
			"records_with_address": gin.H{
				"counselors_with_location": counselorWithLocationCount,
				"centers_with_address":     centerWithAddressCount,
				"doctors_with_description": doctorWithDescriptionCount,
				"total_addressable":        counselorWithLocationCount + centerWithAddressCount + doctorWithDescriptionCount,
			},
		},
		"message": "資料庫統計資訊",
	})
}
