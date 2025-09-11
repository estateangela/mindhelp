# MindHelp Backend API 測試腳本
param(
    [string]$BaseUrl = "http://localhost:8080",
    [switch]$Verbose
)

Write-Host "🧪 MindHelp Backend API 測試開始..." -ForegroundColor Green

# 等待伺服器啟動
Start-Sleep -Seconds 2

function Test-ApiEndpoint {
    param(
        [string]$Method,
        [string]$Url,
        [string]$Description,
        [hashtable]$Headers = @{},
        [string]$Body = $null
    )
    
    Write-Host "📋 測試: $Description" -ForegroundColor Yellow
    Write-Host "   方法: $Method $Url" -ForegroundColor Cyan
    
    try {
        $params = @{
            Uri = $Url
            Method = $Method
            Headers = $Headers
            ContentType = "application/json"
            TimeoutSec = 10
        }
        
        if ($Body) {
            $params.Body = $Body
        }
        
        $response = Invoke-WebRequest @params
        $statusCode = $response.StatusCode
        $responseBody = $response.Content
        
        if ($statusCode -ge 200 -and $statusCode -lt 300) {
            Write-Host "   ✅ 成功 ($statusCode)" -ForegroundColor Green
            if ($Verbose -and $responseBody) {
                Write-Host "   回應: $($responseBody.Substring(0, [Math]::Min(100, $responseBody.Length)))..." -ForegroundColor Gray
            }
            return @{ Success = $true; StatusCode = $statusCode; Body = $responseBody }
        } else {
            Write-Host "   ❌ 失敗 ($statusCode)" -ForegroundColor Red
            return @{ Success = $false; StatusCode = $statusCode; Body = $responseBody }
        }
    }
    catch {
        Write-Host "   ❌ 錯誤: $($_.Exception.Message)" -ForegroundColor Red
        return @{ Success = $false; Error = $_.Exception.Message }
    }
}

# 測試結果統計
$testResults = @()

# 1. 健康檢查
$result = Test-ApiEndpoint -Method "GET" -Url "$BaseUrl/health" -Description "健康檢查端點"
$testResults += $result

# 2. Swagger 文檔
$result = Test-ApiEndpoint -Method "GET" -Url "$BaseUrl/swagger/index.html" -Description "Swagger 文檔"
$testResults += $result

# 3. 用戶註冊（應該會失敗，因為需要請求體）
$registerBody = @{
    email = "test@mindhelp.com"
    password = "testpassword123"
    username = "testuser"
    full_name = "Test User"
} | ConvertTo-Json

$result = Test-ApiEndpoint -Method "POST" -Url "$BaseUrl/api/v1/auth/register" -Description "用戶註冊端點" -Body $registerBody
$testResults += $result

# 4. 測試不需要認證的位置搜索
$result = Test-ApiEndpoint -Method "GET" -Url "$BaseUrl/api/v1/locations/search?q=hospital" -Description "位置搜索端點"
$testResults += $result

# 5. 測試需要認證的端點（應該返回 401）
$result = Test-ApiEndpoint -Method "GET" -Url "$BaseUrl/api/v1/chat/history" -Description "聊天歷史端點 (未認證)"
$testResults += $result

# 測試摘要
Write-Host "`n📊 測試摘要:" -ForegroundColor Magenta
$successCount = ($testResults | Where-Object { $_.Success }).Count
$totalCount = $testResults.Count

Write-Host "   成功: $successCount/$totalCount" -ForegroundColor Green
Write-Host "   失敗: $($totalCount - $successCount)/$totalCount" -ForegroundColor Red

if ($successCount -eq $totalCount) {
    Write-Host "🎉 所有API端點測試通過！" -ForegroundColor Green
} else {
    Write-Host "⚠️  部分API端點需要檢查" -ForegroundColor Yellow
}

return $testResults
