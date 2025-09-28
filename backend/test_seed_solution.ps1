# 測試種子資料解決方案

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

Write-ColorOutput Green "=== 測試種子資料解決方案 ==="

# 1. 檢查目前狀態
Write-ColorOutput Yellow "`n步驟 1: 檢查目前資料庫狀態"
try {
    $addressesResponse = Invoke-RestMethod -Uri "$baseUrl/maps/google-addresses" -Method GET
    Write-Host "目前地址總數: $($addressesResponse.data.total)"
    
    if ($addressesResponse.data.total -eq 0) {
        Write-ColorOutput Red "✗ 確認問題：資料庫中沒有地址資料"
    } else {
        Write-ColorOutput Green "✓ 資料庫中已有 $($addressesResponse.data.total) 筆地址資料"
        Write-Host "問題已解決，無需進一步操作！"
        exit 0
    }
} catch {
    Write-ColorOutput Red "✗ 無法檢查地址 API: $($_.Exception.Message)"
}

# 2. 檢查管理員 API (需要認證)
Write-ColorOutput Yellow "`n步驟 2: 檢查管理員 API 可用性"
Write-Host "注意：管理員 API 需要認證，這裡只是檢查端點是否存在"

# 3. 建議的解決方案
Write-ColorOutput Yellow "`n步驟 3: 建議的解決方案"
Write-Host ""
Write-ColorOutput Cyan "方案 A: 使用管理員 API (推薦)"
Write-Host "1. 首先需要登入獲取 JWT token"
Write-Host "2. 使用 token 調用管理員 API:"
Write-Host "   POST $baseUrl/admin/seed-database"
Write-Host ""

Write-ColorOutput Cyan "方案 B: 在伺服器上直接運行 seeder"
Write-Host "1. SSH 到 Render 服務"
Write-Host "2. 執行: go run cmd/seed/main.go"
Write-Host ""

Write-ColorOutput Cyan "方案 C: 修改部署配置"
Write-Host "1. 在 Render 中添加環境變數: RUN_SEEDER=true"
Write-Host "2. 重新部署服務"
Write-Host ""

# 4. 提供測試用的管理員 API 調用範例
Write-ColorOutput Yellow "`n步驟 4: 管理員 API 調用範例"
Write-Host ""
Write-Host "如果您有管理員權限，可以嘗試以下操作："
Write-Host ""

$loginExample = @{
    email = "admin@example.com"
    password = "your-admin-password"
} | ConvertTo-Json

Write-Host "1. 登入獲取 token:"
Write-Host "curl -X POST '$baseUrl/auth/login' \"
Write-Host "  -H 'Content-Type: application/json' \"
Write-Host "  -d '$loginExample'"
Write-Host ""

Write-Host "2. 使用 token 插入種子資料:"
Write-Host "curl -X POST '$baseUrl/admin/seed-database' \"
Write-Host "  -H 'Authorization: Bearer YOUR_JWT_TOKEN' \"
Write-Host "  -H 'Content-Type: application/json'"
Write-Host ""

Write-Host "3. 檢查資料庫統計:"
Write-Host "curl -X GET '$baseUrl/admin/database-stats' \"
Write-Host "  -H 'Authorization: Bearer YOUR_JWT_TOKEN'"
Write-Host ""

# 5. 驗證步驟
Write-ColorOutput Yellow "`n步驟 5: 執行後驗證"
Write-Host "執行任一解決方案後，運行以下命令驗證："
Write-Host ""
Write-Host "# 檢查地址資料"
Write-Host "curl '$baseUrl/maps/google-addresses'"
Write-Host ""
Write-Host "# 檢查個別資料"
Write-Host "curl '$baseUrl/counselors'"
Write-Host "curl '$baseUrl/counseling-centers'"  
Write-Host "curl '$baseUrl/recommended-doctors'"
Write-Host ""

# 6. 預期結果
Write-ColorOutput Yellow "`n步驟 6: 預期結果"
Write-Host "成功執行種子資料後，您應該看到："
Write-Host ""
$expectedResult = @{
    data = @{
        addresses = @(
            @{
                id = "..."
                name = "王心理師"
                address = "台北市信義區信義路五段7號101大樓"
                type = "counselor"
            }
        )
        total = 11
        format = "google_maps_ready"
    }
    success = $true
} | ConvertTo-Json -Depth 3

Write-Host $expectedResult
Write-Host ""

Write-ColorOutput Green "`n=== 總結 ==="
Write-Host "問題：API 回應地址資料為空 (total: 0)"
Write-Host "原因：資料庫中沒有諮商師、諮商所、推薦醫師的地址資料"
Write-Host "解決：運行資料種子程序插入範例資料"
Write-Host "結果：Google Maps API 將有 11 筆地址資料可用"
Write-Host ""
Write-ColorOutput Cyan "建議立即行動："
Write-Host "1. 嘗試管理員 API 方案 (如果有權限)"
Write-Host "2. 或聯繫系統管理員在伺服器上運行 seeder"
Write-Host "3. 運行 .\debug_database.ps1 監控進度"
