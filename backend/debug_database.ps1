# 診斷資料庫狀態的腳本

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

Write-ColorOutput Green "=== 診斷資料庫狀態 ==="

# 1. 檢查健康狀況
Write-ColorOutput Yellow "`n1. 檢查 API 健康狀況"
try {
    $healthResponse = Invoke-RestMethod -Uri "https://mindhelp.onrender.com/health" -Method GET
    Write-ColorOutput Green "✓ API 健康狀況正常"
    Write-Host "狀態: $($healthResponse.status)"
} catch {
    Write-ColorOutput Red "✗ API 健康檢查失敗: $($_.Exception.Message)"
}

# 2. 檢查詳細健康狀況 (包含資料庫連接)
Write-ColorOutput Yellow "`n2. 檢查詳細健康狀況"
try {
    $detailedHealthResponse = Invoke-RestMethod -Uri "https://mindhelp.onrender.com/health/detailed" -Method GET
    Write-ColorOutput Green "✓ 詳細健康檢查成功"
    Write-Host "資料庫狀態: $($detailedHealthResponse.database.status)"
    Write-Host "資料庫類型: $($detailedHealthResponse.database.type)"
} catch {
    Write-ColorOutput Red "✗ 詳細健康檢查失敗: $($_.Exception.Message)"
}

# 3. 檢查諮商師資料
Write-ColorOutput Yellow "`n3. 檢查諮商師資料"
try {
    $counselorsResponse = Invoke-RestMethod -Uri "$baseUrl/counselors" -Method GET
    Write-ColorOutput Green "✓ 諮商師 API 可用"
    Write-Host "諮商師總數: $($counselorsResponse.data.counselors.Count)"
    if ($counselorsResponse.data.counselors.Count -gt 0) {
        $counselorsWithLocation = $counselorsResponse.data.counselors | Where-Object { $_.work_location -ne $null -and $_.work_location -ne "" }
        Write-Host "有工作地點的諮商師: $($counselorsWithLocation.Count)"
        if ($counselorsWithLocation.Count -gt 0) {
            Write-Host "範例工作地點: $($counselorsWithLocation[0].work_location)"
        }
    }
} catch {
    Write-ColorOutput Red "✗ 諮商師 API 失敗: $($_.Exception.Message)"
}

# 4. 檢查諮商所資料
Write-ColorOutput Yellow "`n4. 檢查諮商所資料"
try {
    $centersResponse = Invoke-RestMethod -Uri "$baseUrl/counseling-centers" -Method GET
    Write-ColorOutput Green "✓ 諮商所 API 可用"
    Write-Host "諮商所總數: $($centersResponse.data.counseling_centers.Count)"
    if ($centersResponse.data.counseling_centers.Count -gt 0) {
        $centersWithAddress = $centersResponse.data.counseling_centers | Where-Object { $_.address -ne $null -and $_.address -ne "" }
        Write-Host "有地址的諮商所: $($centersWithAddress.Count)"
        if ($centersWithAddress.Count -gt 0) {
            Write-Host "範例地址: $($centersWithAddress[0].address)"
        }
    }
} catch {
    Write-ColorOutput Red "✗ 諮商所 API 失敗: $($_.Exception.Message)"
}

# 5. 檢查推薦醫師資料
Write-ColorOutput Yellow "`n5. 檢查推薦醫師資料"
try {
    $doctorsResponse = Invoke-RestMethod -Uri "$baseUrl/recommended-doctors" -Method GET
    Write-ColorOutput Green "✓ 推薦醫師 API 可用"
    Write-Host "推薦醫師總數: $($doctorsResponse.data.recommended_doctors.Count)"
    if ($doctorsResponse.data.recommended_doctors.Count -gt 0) {
        $doctorsWithDescription = $doctorsResponse.data.recommended_doctors | Where-Object { $_.description -ne $null -and $_.description -ne "" }
        Write-Host "有描述的推薦醫師: $($doctorsWithDescription.Count)"
        if ($doctorsWithDescription.Count -gt 0) {
            Write-Host "範例描述: $($doctorsWithDescription[0].description.Substring(0, [Math]::Min(50, $doctorsWithDescription[0].description.Length)))..."
        }
    }
} catch {
    Write-ColorOutput Red "✗ 推薦醫師 API 失敗: $($_.Exception.Message)"
}

# 6. 重新檢查地址 API
Write-ColorOutput Yellow "`n6. 重新檢查地址 API"
try {
    $addressesResponse = Invoke-RestMethod -Uri "$baseUrl/maps/addresses" -Method GET
    Write-ColorOutput Green "✓ 地址 API 可用"
    Write-Host "總地址數: $($addressesResponse.data.total)"
    if ($addressesResponse.data.addresses -ne $null) {
        Write-Host "地址詳情:"
        $addressesResponse.data.addresses | ForEach-Object {
            Write-Host "  - $($_.type): $($_.name) ($($_.address))"
        }
    }
} catch {
    Write-ColorOutput Red "✗ 地址 API 失敗: $($_.Exception.Message)"
}

Write-ColorOutput Green "`n=== 診斷完成 ==="
Write-Host "`n建議下一步:"
Write-Host "1. 如果資料庫表為空，需要運行資料種子 (seeder)"
Write-Host "2. 如果有資料但地址欄位為空，需要更新資料"
Write-Host "3. 如果 API 無法連接資料庫，檢查資料庫配置"
