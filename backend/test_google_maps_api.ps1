# Google Maps API 測試腳本
# 使用方法: .\test_google_maps_api.ps1

$baseUrl = "http://localhost:8080/api/v1/google-maps"

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

Write-ColorOutput Green "=== Google Maps API 測試開始 ==="

# 1. 測試地理編碼
Write-ColorOutput Yellow "`n1. 測試地理編碼 (Geocoding)"
$geocodeData = @{
    address = "台北市信義區市府路1號"
    language = "zh-TW"
    region = "tw"
} | ConvertTo-Json

try {
    $response = Invoke-RestMethod -Uri "$baseUrl/geocode" -Method POST -Body $geocodeData -ContentType "application/json"
    Write-ColorOutput Green "✓ 地理編碼測試成功"
    Write-Host "結果數量: $($response.results.Count)"
    if ($response.results.Count -gt 0) {
        $location = $response.results[0].geometry.location
        Write-Host "座標: $($location.lat), $($location.lng)"
    }
} catch {
    Write-ColorOutput Red "✗ 地理編碼測試失敗: $($_.Exception.Message)"
}

# 2. 測試反向地理編碼
Write-ColorOutput Yellow "`n2. 測試反向地理編碼 (Reverse Geocoding)"
$reverseGeocodeData = @{
    latitude = 25.0330
    longitude = 121.5654
    language = "zh-TW"
} | ConvertTo-Json

try {
    $response = Invoke-RestMethod -Uri "$baseUrl/reverse-geocode" -Method POST -Body $reverseGeocodeData -ContentType "application/json"
    Write-ColorOutput Green "✓ 反向地理編碼測試成功"
    Write-Host "結果數量: $($response.results.Count)"
    if ($response.results.Count -gt 0) {
        Write-Host "地址: $($response.results[0].formatted_address)"
    }
} catch {
    Write-ColorOutput Red "✗ 反向地理編碼測試失敗: $($_.Exception.Message)"
}

# 3. 測試地點搜尋
Write-ColorOutput Yellow "`n3. 測試地點搜尋 (Places Search)"
$placesSearchData = @{
    query = "心理諮商 台北"
    location = "25.0330,121.5654"
    radius = 5000
    type = "health"
    language = "zh-TW"
    region = "tw"
} | ConvertTo-Json

try {
    $response = Invoke-RestMethod -Uri "$baseUrl/search-places" -Method POST -Body $placesSearchData -ContentType "application/json"
    Write-ColorOutput Green "✓ 地點搜尋測試成功"
    Write-Host "結果數量: $($response.results.Count)"
    if ($response.results.Count -gt 0) {
        Write-Host "第一個結果: $($response.results[0].name)"
    }
} catch {
    Write-ColorOutput Red "✗ 地點搜尋測試失敗: $($_.Exception.Message)"
}

# 4. 測試路線規劃
Write-ColorOutput Yellow "`n4. 測試路線規劃 (Directions)"
$directionsData = @{
    origin = "台北車站"
    destination = "台北101"
    mode = "driving"
    language = "zh-TW"
    region = "tw"
    alternatives = $true
} | ConvertTo-Json

try {
    $response = Invoke-RestMethod -Uri "$baseUrl/directions" -Method POST -Body $directionsData -ContentType "application/json"
    Write-ColorOutput Green "✓ 路線規劃測試成功"
    Write-Host "路線數量: $($response.routes.Count)"
    if ($response.routes.Count -gt 0) {
        $route = $response.routes[0]
        if ($route.legs.Count -gt 0) {
            Write-Host "距離: $($route.legs[0].distance.text)"
            Write-Host "時間: $($route.legs[0].duration.text)"
        }
    }
} catch {
    Write-ColorOutput Red "✗ 路線規劃測試失敗: $($_.Exception.Message)"
}

# 5. 測試附近心理健康服務
Write-ColorOutput Yellow "`n5. 測試附近心理健康服務搜尋"
$nearbyParams = @{
    latitude = 25.0330
    longitude = 121.5654
    radius = 5000
    keyword = "心理諮商"
}

$queryString = ($nearbyParams.GetEnumerator() | ForEach-Object { "$($_.Key)=$($_.Value)" }) -join "&"

try {
    $response = Invoke-RestMethod -Uri "$baseUrl/nearby-mental-health?$queryString" -Method GET
    Write-ColorOutput Green "✓ 附近心理健康服務搜尋測試成功"
    Write-Host "結果數量: $($response.results.Count)"
    if ($response.results.Count -gt 0) {
        Write-Host "第一個結果: $($response.results[0].name)"
    }
} catch {
    Write-ColorOutput Red "✗ 附近心理健康服務搜尋測試失敗: $($_.Exception.Message)"
}

# 6. 測試批次地理編碼
Write-ColorOutput Yellow "`n6. 測試批次地理編碼 (Batch Geocoding)"
$batchGeocodeData = @{
    addresses = @(
        "台北車站",
        "台北101",
        "西門町"
    )
    language = "zh-TW"
    region = "tw"
} | ConvertTo-Json

try {
    $response = Invoke-RestMethod -Uri "$baseUrl/batch-geocode" -Method POST -Body $batchGeocodeData -ContentType "application/json"
    Write-ColorOutput Green "✓ 批次地理編碼測試成功"
    Write-Host "處理地址數量: $($response.total)"
} catch {
    Write-ColorOutput Red "✗ 批次地理編碼測試失敗: $($_.Exception.Message)"
}

# 7. 測試 API 使用統計
Write-ColorOutput Yellow "`n7. 測試 API 使用統計"
try {
    $response = Invoke-RestMethod -Uri "$baseUrl/usage-stats" -Method GET
    Write-ColorOutput Green "✓ API 使用統計測試成功"
    Write-Host "快取條目數量: $($response.data.cache_stats.total_entries)"
    Write-Host "API Key 已配置: $($response.data.api_info.api_key_configured)"
} catch {
    Write-ColorOutput Red "✗ API 使用統計測試失敗: $($_.Exception.Message)"
}

# 8. 測試清除快取
Write-ColorOutput Yellow "`n8. 測試清除快取"
try {
    $response = Invoke-RestMethod -Uri "$baseUrl/clear-cache" -Method POST
    Write-ColorOutput Green "✓ 清除快取測試成功"
    Write-Host "回應: $($response.message)"
} catch {
    Write-ColorOutput Red "✗ 清除快取測試失敗: $($_.Exception.Message)"
}

Write-ColorOutput Green "`n=== Google Maps API 測試完成 ==="
Write-Host "`n注意事項:"
Write-Host "1. 確保已設定 GOOGLE_MAPS_API_KEY 環境變數"
Write-Host "2. 確保後端伺服器正在運行 (預設 localhost:8080)"
Write-Host "3. 某些測試可能需要有效的 Google Maps API 金鑰才能成功"
