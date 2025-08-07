# Test script for Redis-based view tracking API

Write-Host "🧪 Testing Anchor Blog API with Redis View Tracking" -ForegroundColor Green
Write-Host "=================================================" -ForegroundColor Green

$baseUrl = "http://localhost:8080"

# Test 1: Health Check
Write-Host "`n1. Testing Health Check..." -ForegroundColor Yellow
try {
    $response = Invoke-RestMethod -Uri "$baseUrl/health" -Method GET
    Write-Host "✅ Health Check: $($response.status)" -ForegroundColor Green
} catch {
    Write-Host "❌ Health Check Failed: $($_.Exception.Message)" -ForegroundColor Red
}

# Test 2: Register a new user
Write-Host "`n2. Testing User Registration..." -ForegroundColor Yellow
$registerData = @{
    username = "testuser$(Get-Random)"
    first_name = "Test"
    last_name = "User"
    email = "test$(Get-Random)@example.com"
    password = "testpassword123"
    role = "user"
    profile = @{
        bio = "Test user for Redis view tracking"
        picture_url = "https://example.com/avatar.jpg"
        social_links = @(
            @{
                platform = "github"
                url = "https://github.com/testuser"
            }
        )
    }
} | ConvertTo-Json -Depth 3

try {
    $registerResponse = Invoke-RestMethod -Uri "$baseUrl/api/v1/user/register" -Method POST -Body $registerData -ContentType "application/json"
    Write-Host "✅ User Registration: ID = $($registerResponse.id)" -ForegroundColor Green
    $userId = $registerResponse.id
} catch {
    Write-Host "❌ User Registration Failed: $($_.Exception.Message)" -ForegroundColor Red
}

# Test 3: Login
Write-Host "`n3. Testing User Login..." -ForegroundColor Yellow
$loginData = @{
    username = ($registerData | ConvertFrom-Json).username
    password = "testpassword123"
} | ConvertTo-Json

try {
    $loginResponse = Invoke-RestMethod -Uri "$baseUrl/api/v1/user/login" -Method POST -Body $loginData -ContentType "application/json"
    Write-Host "✅ User Login: Access token received" -ForegroundColor Green
    $accessToken = $loginResponse.access_token
    $headers = @{ Authorization = "Bearer $accessToken" }
} catch {
    Write-Host "❌ User Login Failed: $($_.Exception.Message)" -ForegroundColor Red
    return
}

# Test 4: Create a blog post
Write-Host "`n4. Testing Blog Post Creation..." -ForegroundColor Yellow
$postData = @{
    title = "Test Post for Redis View Tracking"
    content = "This is a test post to demonstrate Redis-based view tracking with IP throttling. Each unique IP can only increment the view count once per 24 hours."
    tags = @("redis", "view-tracking", "test")
} | ConvertTo-Json

try {
    $postResponse = Invoke-RestMethod -Uri "$baseUrl/api/v1/posts" -Method POST -Body $postData -ContentType "application/json" -Headers $headers
    Write-Host "✅ Blog Post Created: ID = $($postResponse.id)" -ForegroundColor Green
    $postId = $postResponse.id
} catch {
    Write-Host "❌ Blog Post Creation Failed: $($_.Exception.Message)" -ForegroundColor Red
    return
}

# Test 5: View the post multiple times (should only increment once due to IP throttling)
Write-Host "`n5. Testing Redis View Tracking (IP Throttling)..." -ForegroundColor Yellow

Write-Host "   First view (should increment view count):" -ForegroundColor Cyan
try {
    $viewResponse1 = Invoke-RestMethod -Uri "$baseUrl/api/v1/posts/$postId" -Method GET
    Write-Host "   ✅ Post viewed: View count = $($viewResponse1.view_count)" -ForegroundColor Green
} catch {
    Write-Host "   ❌ First view failed: $($_.Exception.Message)" -ForegroundColor Red
}

Start-Sleep -Seconds 1

Write-Host "   Second view (should NOT increment due to IP throttling):" -ForegroundColor Cyan
try {
    $viewResponse2 = Invoke-RestMethod -Uri "$baseUrl/api/v1/posts/$postId" -Method GET
    Write-Host "   ✅ Post viewed again: View count = $($viewResponse2.view_count)" -ForegroundColor Green
    
    if ($viewResponse1.view_count -eq $viewResponse2.view_count) {
        Write-Host "   🎉 IP Throttling Working! View count didn't increase on second view" -ForegroundColor Green
    } else {
        Write-Host "   ⚠️  IP Throttling may not be working - view count increased" -ForegroundColor Yellow
    }
} catch {
    Write-Host "   ❌ Second view failed: $($_.Exception.Message)" -ForegroundColor Red
}

# Test 6: Get post view count
Write-Host "`n6. Testing Get Post View Count..." -ForegroundColor Yellow
try {
    $viewCountResponse = Invoke-RestMethod -Uri "$baseUrl/api/v1/posts/$postId/views" -Method GET
    Write-Host "✅ Post View Count: $($viewCountResponse.view_count) views" -ForegroundColor Green
} catch {
    Write-Host "❌ Get View Count Failed: $($_.Exception.Message)" -ForegroundColor Red
}

# Test 7: Get total view statistics
Write-Host "`n7. Testing Total View Statistics..." -ForegroundColor Yellow
try {
    $statsResponse = Invoke-RestMethod -Uri "$baseUrl/api/v1/stats/views" -Method GET
    Write-Host "✅ Total Views Across All Posts: $($statsResponse.total_views)" -ForegroundColor Green
} catch {
    Write-Host "❌ Get Total Views Failed: $($_.Exception.Message)" -ForegroundColor Red
}

# Test 8: Get popular posts
Write-Host "`n8. Testing Popular Posts..." -ForegroundColor Yellow
try {
    $popularResponse = Invoke-RestMethod -Uri "$baseUrl/api/v1/posts/popular?limit=5" -Method GET
    Write-Host "✅ Popular Posts Retrieved: $($popularResponse.count) posts" -ForegroundColor Green
    if ($popularResponse.posts.Count -gt 0) {
        Write-Host "   Most popular post: '$($popularResponse.posts[0].title)' with $($popularResponse.posts[0].view_count) views" -ForegroundColor Cyan
    }
} catch {
    Write-Host "❌ Get Popular Posts Failed: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host "`n🎉 Redis View Tracking Test Complete!" -ForegroundColor Green
Write-Host "=================================================" -ForegroundColor Green