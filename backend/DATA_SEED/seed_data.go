package main

// This is a standalone script for seeding data
// Run with: go run seed_data.go

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"mindhelp-backend/internal/config"
	"mindhelp-backend/internal/database"
	"mindhelp-backend/internal/models"
)

func seedData() {
	// 載入配置
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 連接到資料庫
	if err := database.Connect(cfg); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("開始插入資料...")

	// 插入諮商師資料
	//if err := insertCounselors(); err != nil {
	//	log.Fatalf("Failed to insert counselors: %v", err)
	//}

	// 插入諮商所資料
	//if err := insertCounselingCenters(); err != nil {
	//	log.Fatalf("Failed to insert counseling centers: %v", err)
	//}

	// 插入推薦醫師資料
	if err := insertRecommendedDoctors(); err != nil {
		log.Fatalf("Failed to insert recommended doctors: %v", err)
	}

	log.Println("所有資料插入完成！")
}

func insertCounselors() error {
	log.Println("插入諮商師資料...")

	// 讀取 CSV 檔案
	file, err := os.Open("document/專題功能表 (1).xlsx - 諮商師.csv")
	if err != nil {
		return fmt.Errorf("failed to open counselor CSV file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to read counselor CSV file: %v", err)
	}

	// 跳過標題行
	for i, record := range records[1:] {
		if len(record) < 10 {
			continue // 跳過不完整的記錄
		}

		counselor := models.Counselor{
			Name:             strings.TrimSpace(record[0]),
			LicenseNumber:    strings.TrimSpace(record[1]),
			Gender:           strings.TrimSpace(record[2]),
			Specialties:      strings.TrimSpace(record[3]),
			LanguageSkills:   strings.TrimSpace(record[4]),
			WorkLocation:     strings.TrimSpace(record[5]),
			WorkUnit:         strings.TrimSpace(record[6]),
			InstitutionCode:  strings.TrimSpace(record[7]),
			PsychologySchool: strings.TrimSpace(record[8]),
			TreatmentMethods: strings.TrimSpace(record[9]),
		}

		// 檢查是否已存在
		var existing models.Counselor
		result := database.GetDB().Where("license_number = ?", counselor.LicenseNumber).First(&existing)
		if result.Error == nil {
			// 更新現有記錄
			counselor.ID = existing.ID
			if err := database.GetDB().Save(&counselor).Error; err != nil {
				log.Printf("Failed to update counselor %d: %v", i+1, err)
			}
		} else {
			// 創建新記錄
			if err := database.GetDB().Create(&counselor).Error; err != nil {
				log.Printf("Failed to create counselor %d: %v", i+1, err)
			}
		}
	}

	log.Printf("諮商師資料插入完成，共處理 %d 筆記錄", len(records)-1)
	return nil
}

func insertCounselingCenters() error {
	log.Println("插入諮商所資料...")

	// 讀取 CSV 檔案
	file, err := os.Open("document/專題功能表 (1).xlsx - 台北諮商所.csv")
	if err != nil {
		return fmt.Errorf("failed to open counseling center CSV file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to read counseling center CSV file: %v", err)
	}

	// 跳過標題行
	for i, record := range records[2:] { // 跳過前兩行標題
		if len(record) < 4 {
			continue // 跳過不完整的記錄
		}

		// 檢查是否有通訊心理諮商
		onlineCounseling := false
		if len(record) > 3 && strings.TrimSpace(record[3]) == "V" {
			onlineCounseling = true
		}

		center := models.CounselingCenter{
			Name:             strings.TrimSpace(record[0]),
			Address:          strings.TrimSpace(record[1]),
			Phone:            strings.TrimSpace(record[2]),
			OnlineCounseling: onlineCounseling,
		}

		// 檢查是否已存在
		var existing models.CounselingCenter
		result := database.GetDB().Where("name = ? AND address = ?", center.Name, center.Address).First(&existing)
		if result.Error == nil {
			// 更新現有記錄
			center.ID = existing.ID
			if err := database.GetDB().Save(&center).Error; err != nil {
				log.Printf("Failed to update counseling center %d: %v", i+1, err)
			}
		} else {
			// 創建新記錄
			if err := database.GetDB().Create(&center).Error; err != nil {
				log.Printf("Failed to create counseling center %d: %v", i+1, err)
			}
		}
	}

	log.Printf("諮商所資料插入完成，共處理 %d 筆記錄", len(records)-2)
	return nil
}

func insertRecommendedDoctors() error {
	log.Println("插入推薦醫師資料...")

	// 讀取 CSV 檔案
	file, err := os.Open("document/專題功能表 (1).xlsx - 網友推薦醫師＆診所.csv")
	if err != nil {
		return fmt.Errorf("failed to open recommended doctor CSV file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to read recommended doctor CSV file: %v", err)
	}

	// 處理每一行記錄
	for i, record := range records {
		if len(record) == 0 || strings.TrimSpace(record[0]) == "" {
			continue
		}

		description := strings.TrimSpace(record[0])

		// 嘗試從描述中提取經驗次數
		experienceCount := 0
		if strings.Contains(description, "諮商") {
			// 尋找數字
			parts := strings.Split(description, "諮商")
			if len(parts) > 1 {
				// 在諮商後面尋找數字
				afterCounseling := parts[1]
				words := strings.Fields(afterCounseling)
				for _, word := range words {
					if num, err := strconv.Atoi(word); err == nil {
						experienceCount = num
						break
					}
				}
			}
		}

		// 從描述中提取名稱（通常是第一個部分）
		name := description
		if strings.Contains(description, " - ") {
			name = strings.Split(description, " - ")[0]
		} else if strings.Contains(description, "：") {
			name = strings.Split(description, "：")[0]
		}

		doctor := models.RecommendedDoctor{
			Name:            strings.TrimSpace(name),
			Description:     description,
			ExperienceCount: experienceCount,
		}

		// 檢查是否已存在
		var existing models.RecommendedDoctor
		result := database.GetDB().Where("name = ? AND description = ?", doctor.Name, doctor.Description).First(&existing)
		if result.Error == nil {
			// 更新現有記錄
			doctor.ID = existing.ID
			if err := database.GetDB().Save(&doctor).Error; err != nil {
				log.Printf("Failed to update recommended doctor %d: %v", i+1, err)
			}
		} else {
			// 創建新記錄
			if err := database.GetDB().Create(&doctor).Error; err != nil {
				log.Printf("Failed to create recommended doctor %d: %v", i+1, err)
			}
		}
	}

	log.Printf("推薦醫師資料插入完成，共處理 %d 筆記錄", len(records))
	return nil
}

func main() {
	seedData()
}
