# 創建範例資料的腳本

$baseUrl = "https://mindhelp.onrender.com/api/v1"

# 顏色函數
function Write-ColorOutput($ForegroundColor) {
    $fc = $host.UI.RawUI.ForegroundColor
    $host.UI.RawUI.ForegroundColor = $ForegroundColor
    if ($args) {
        Write-Output $args
    }
    else {
        $input | Write-Output
    }
    $host.UI.RawUI.ForegroundColor = $fc
}

Write-ColorOutput Green "=== 創建範例資料 ==="

# 注意：這個腳本需要管理員權限來創建資料
# 實際生產環境中應該通過管理 API 或直接在伺服器上運行 seeder

Write-ColorOutput Yellow "`n1. 檢查目前資料狀態"
try {
    # 檢查諮商師
    $counselorsResponse = Invoke-RestMethod -Uri "$baseUrl/counselors" -Method GET
    Write-Host "目前諮商師數量: $($counselorsResponse.data.counselors.Count)"
    
    # 檢查諮商所
    $centersResponse = Invoke-RestMethod -Uri "$baseUrl/counseling-centers" -Method GET  
    Write-Host "目前諮商所數量: $($centersResponse.data.counseling_centers.Count)"
    
    # 檢查推薦醫師
    $doctorsResponse = Invoke-RestMethod -Uri "$baseUrl/recommended-doctors" -Method GET
    Write-Host "目前推薦醫師數量: $($doctorsResponse.data.recommended_doctors.Count)"
    
} catch {
    Write-ColorOutput Red "✗ 無法獲取目前資料狀態: $($_.Exception.Message)"
}

Write-ColorOutput Yellow "`n2. 建議的解決方案"
Write-Host "由於 API 返回空地址資料，可能的原因和解決方案："
Write-Host ""
Write-Host "原因 1: 資料庫表為空"
Write-ColorOutput Cyan "解決方案 1A: 在伺服器上運行 seeder"
Write-Host "  cd /app && go run cmd/seed/main.go"
Write-Host ""
Write-ColorOutput Cyan "解決方案 1B: 手動創建範例資料 (需要管理員 API)"
Write-Host "  需要實現管理員 API 端點來創建資料"
Write-Host ""

Write-Host "原因 2: 資料存在但地址欄位為空"
Write-ColorOutput Cyan "解決方案 2: 更新現有資料的地址欄位"
Write-Host "  UPDATE counselors SET work_location = '台北市信義區' WHERE work_location IS NULL;"
Write-Host "  UPDATE counseling_centers SET address = '台北市大安區' WHERE address IS NULL;"
Write-Host ""

Write-Host "原因 3: 查詢條件太嚴格"
Write-ColorOutput Cyan "解決方案 3: 檢查查詢邏輯"
Write-Host "  檢查 maps_handler.go 中的 WHERE 條件"
Write-Host ""

Write-ColorOutput Yellow "`n3. 測試用範例資料 (JSON 格式)"
Write-Host "如果需要手動創建資料，可以使用以下範例："

$sampleCounselors = @(
    @{
        name = "王心理師"
        license_number = "001"
        gender = "女"
        specialties = "焦慮症、憂鬱症"
        language_skills = "中文、英文"
        work_location = "台北市信義區信義路五段7號"
        work_unit = "台北心理診所"
        institution_code = "TP001"
        psychology_school = "台灣大學心理系"
        treatment_methods = "認知行為療法"
    },
    @{
        name = "李諮商師"
        license_number = "002"  
        gender = "男"
        specialties = "家庭治療、伴侶諮商"
        language_skills = "中文"
        work_location = "台北市大安區復興南路一段390號"
        work_unit = "大安諮商中心"
        institution_code = "TP002"
        psychology_school = "政治大學心理系"
        treatment_methods = "系統性家族治療"
    }
)

$sampleCenters = @(
    @{
        name = "台北心理健康中心"
        address = "台北市中正區中山南路1號"
        phone = "02-2311-1234"
        online_counseling = $true
    },
    @{
        name = "信義諮商所"
        address = "台北市信義區信義路四段1號"
        phone = "02-2722-5678"
        online_counseling = $false
    }
)

Write-Host "範例諮商師資料:"
$sampleCounselors | ConvertTo-Json -Depth 2
Write-Host ""
Write-Host "範例諮商所資料:"
$sampleCenters | ConvertTo-Json -Depth 2

Write-ColorOutput Yellow "`n4. 下一步建議"
Write-Host "1. 運行診斷腳本: .\debug_database.ps1"
Write-Host "2. 在伺服器上執行: go run cmd/seed/main.go"
Write-Host "3. 或者聯繫系統管理員添加範例資料"
Write-Host "4. 檢查資料庫連接和權限設定"

Write-ColorOutput Green "`n=== 完成 ==="
