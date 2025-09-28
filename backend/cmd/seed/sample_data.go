package main

import (
	"log"

	"mindhelp-backend/internal/database"
	"mindhelp-backend/internal/models"
)

// insertSampleData 插入範例資料
func insertSampleData() error {
	log.Println("插入範例資料...")

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

	// 插入諮商師
	for _, counselor := range counselors {
		var existing models.Counselor
		result := database.GetDB().Where("license_number = ?", counselor.LicenseNumber).First(&existing)
		if result.Error != nil {
			// 不存在，創建新記錄
			if err := database.GetDB().Create(&counselor).Error; err != nil {
				log.Printf("Failed to create counselor %s: %v", counselor.Name, err)
			} else {
				log.Printf("Created counselor: %s", counselor.Name)
			}
		} else {
			log.Printf("Counselor already exists: %s", counselor.Name)
		}
	}

	// 插入諮商所
	for _, center := range centers {
		var existing models.CounselingCenter
		result := database.GetDB().Where("name = ? AND address = ?", center.Name, center.Address).First(&existing)
		if result.Error != nil {
			// 不存在，創建新記錄
			if err := database.GetDB().Create(&center).Error; err != nil {
				log.Printf("Failed to create counseling center %s: %v", center.Name, err)
			} else {
				log.Printf("Created counseling center: %s", center.Name)
			}
		} else {
			log.Printf("Counseling center already exists: %s", center.Name)
		}
	}

	// 插入推薦醫師
	for _, doctor := range doctors {
		var existing models.RecommendedDoctor
		result := database.GetDB().Where("name = ?", doctor.Name).First(&existing)
		if result.Error != nil {
			// 不存在，創建新記錄
			if err := database.GetDB().Create(&doctor).Error; err != nil {
				log.Printf("Failed to create recommended doctor %s: %v", doctor.Name, err)
			} else {
				log.Printf("Created recommended doctor: %s", doctor.Name)
			}
		} else {
			log.Printf("Recommended doctor already exists: %s", doctor.Name)
		}
	}

	log.Printf("範例資料插入完成！")
	log.Printf("- 諮商師: %d 筆", len(counselors))
	log.Printf("- 諮商所: %d 筆", len(centers))
	log.Printf("- 推薦醫師: %d 筆", len(doctors))
	
	return nil
}
