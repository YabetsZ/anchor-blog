# Test script for High Priority Features Implementation

Write-Host "🚀 Testing High Priority Features Implementation" -ForegroundColor Green
Write-Host "=================================================" -ForegroundColor Green

$baseUrl = "http://localhost:8080"
$accessToken = ""
$userId = ""
$postId = ""

# Test 1: Health Check
Write-Host "`n1. Testing Health Check..." -ForegroundColor Yellow
try {
    $response = Invoke-RestMethod -Uri "$baseUrl/health" -Method GET
    Write-Host "✅ Health Check: $($response.status)" -ForegroundColor Green
} catch {
    Write-Host "❌ Health Check Failed: $($_.Exception.Message)" -ForegroundColor Red
    exit 1
}

# Test 2: Register and Login
Write-Host "`n2. Testing User Registration and Login..." -ForegroundColor Yellow
$registerData = @{
    username = "testuser$(Get-Random)"
    first_name = "Test"
    last_name = "User"
    email = "test$(Get-Random)@example.com"
    password = "testpassword123"
    role = "user"
    profile = @{
        bio = "Test user for high priority features"
        picture_url = "https://example.com/avatar.jpg"
        social_links = @()
    }
} | ConvertTo-Json -Depth 3

try {
    $registerResponse = Invoke-RestMethod -Uri "$baseUrl/api/v1/user/register" -Method POST -Body $registerData -ContentType "application/json"
    Write-Host "✅ User Registration: ID = $($registerResponse.id)" -ForegroundColor Green
    $userId = $registerResponse.id
} catch {
    Write-Host "❌ User Registration Failed: $($_.Exception.Message)" -ForegroundColor Red
    exit 1
}

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
    exit 1
}

# Test 3: Create a blog post
Write-Host "`n3. Testing Blog Post Creation..." -ForegroundColor Yellow
$postData = @{
    title = "High Priority Features Test Post"
    content = "This post is created to test the new CRUD, search, and like/dislike functionality."
    tags = @("test", "crud", "features")
} | ConvertTo-Json

try {
    $postResponse = Invoke-RestMethod -Uri "$baseUrl/api/v1/posts" -Method POST -Body $postData -ContentType "application/json" -Headers $headers
    Write-Host "✅ Blog Post Created: ID = $($postResponse.id)" -ForegroundColor Green
    $postId = $postResponse.id
} catch {
    Write-Host "❌ Blog Post Creation Failed: $($_.Exception.Message)" -ForegroundColor Red
    exit 1
}

# Test 4: Update the blog post
Write-Host "`n4. Testing Blog Post Update..." -ForegroundColor Yellow
$updateData = @{
    title = "Updated High Priority Features Test Post"
    content = "This post has been updated to test the UPDATE functionality."
    tags = @("test", "crud", "features", "updated")
} | ConvertTo-Json

try {
    $updateResponse = Invoke-RestMethod -Uri "$baseUrl/api/v1/posts/$postId" -Method PUT -Body $updateData -ContentType "application/json" -Headers $headers
    Write-Host "✅ Blog Post Updated: Title = $($updateResponse.title)" -ForegroundColor Green
} catch {
    Write-Host "❌ Blog Post Update Failed: $($_.Exception.Message)" -ForegroundColor Red
}

# Test 5: Search functionality
Write-Host "`n5. Testing Blog Search Functionality..." -ForegroundColor Yellow

Write-Host "   Testing search by title:" -ForegroundColor Cyan
try {
    $searchResponse = Invoke-RestMethod -Uri "$baseUrl/api/v1/posts/search?q=Updated&type=title" -Method GET
    Write-Host "   ✅ Search by title: Found $($searchResponse.count) posts" -ForegroundColor Green
} catch {
    Write-Host "   ❌ Search by title failed: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host "   Testing search by author:" -ForegroundColor Cyan
try {
    $searchResponse = Invoke-RestMethod -Uri "$baseUrl/api/v1/posts/search?q=$userId&type=author" -Method GET
    Write-Host "   ✅ Search by author: Found $($searchResponse.count) posts" -ForegroundColor Green
} catch {
    Write-Host "   ❌ Search by author failed: $($_.Exception.Message)" -ForegroundColor Red
}

# Test 6: Filter functionality
Write-Host "`n6. Testing Blog Filter Functionality..." -ForegroundColor Yellow

Write-Host "   Testing filter by tags:" -ForegroundColor Cyan
try {
    $filterResponse = Invoke-RestMethod -Uri "$baseUrl/api/v1/posts/filter?tags=test,crud" -Method GET
    Write-Host "   ✅ Filter by tags: Found $($filterResponse.count) posts" -ForegroundColor Green
} catch {
    Write-Host "   ❌ Filter by tags failed: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host "   Testing filter by date:" -ForegroundColor Cyan
$today = Get-Date -Format "yyyy-MM-dd"
try {
    $filterResponse = Invoke-RestMethod -Uri "$baseUrl/api/v1/posts/filter?start_date=$today&end_date=$today" -Method GET
    Write-Host "   ✅ Filter by date: Found $($filterResponse.count) posts" -ForegroundColor Green
} catch {
    Write-Host "   ❌ Filter by date failed: $($_.Exception.Message)" -ForegroundColor Red
}

# Test 7: Like/Dislike functionality
Write-Host "`n7. Testing Like/Dislike Functionality..." -ForegroundColor Yellow

Write-Host "   Testing like post:" -ForegroundColor Cyan
try {
    $likeResponse = Invoke-RestMethod -Uri "$baseUrl/api/v1/posts/$postId/like" -Method POST -Headers $headers
    Write-Host "   ✅ Like post: $($likeResponse.message)" -ForegroundColor Green
} catch {
    Write-Host "   ❌ Like post failed: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host "   Testing get like status:" -ForegroundColor Cyan
try {
    $statusResponse = Invoke-RestMethod -Uri "$baseUrl/api/v1/posts/$postId/like-status" -Method GET -Headers $headers
    Write-Host "   ✅ Like status: Liked = $($statusResponse.liked), Disliked = $($statusResponse.disliked)" -ForegroundColor Green
} catch {
    Write-Host "   ❌ Get like status failed: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host "   Testing dislike post:" -ForegroundColor Cyan
try {
    $dislikeResponse = Invoke-RestMethod -Uri "$baseUrl/api/v1/posts/$postId/dislike" -Method POST -Headers $headers
    Write-Host "   ✅ Dislike post: $($dislikeResponse.message)" -ForegroundColor Green
} catch {
    Write-Host "   ❌ Dislike post failed: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host "   Testing get like status after dislike:" -ForegroundColor Cyan
try {
    $statusResponse = Invoke-RestMethod -Uri "$baseUrl/api/v1/posts/$postId/like-status" -Method GET -Headers $headers
    Write-Host "   ✅ Like status: Liked = $($statusResponse.liked), Disliked = $($statusResponse.disliked)" -ForegroundColor Green
} catch {
    Write-Host "   ❌ Get like status failed: $($_.Exception.Message)" -ForegroundColor Red
}

# Test 8: Logout functionality
Write-Host "`n8. Testing User Logout..." -ForegroundColor Yellow
try {
    $logoutResponse = Invoke-RestMethod -Uri "$baseUrl/api/v1/logout" -Method POST -Headers $headers
    Write-Host "✅ User Logout: $($logoutResponse.message)" -ForegroundColor Green
} catch {
    Write-Host "❌ User Logout Failed: $($_.Exception.Message)" -ForegroundColor Red
}

# Test 9: Verify logout (try to access protected endpoint)
Write-Host "`n9. Testing Logout Verification..." -ForegroundColor Yellow
try {
    $profileResponse = Invoke-RestMethod -Uri "$baseUrl/api/v1/user/profile" -Method GET -Headers $headers
    Write-Host "❌ Logout verification failed: Still able to access protected endpoint" -ForegroundColor Red
} catch {
    Write-Host "✅ Logout verification: Cannot access protected endpoint after logout" -ForegroundColor Green
}

# Test 10: Test delete functionality (need to login again)
Write-Host "`n10. Testing Blog Post Deletion..." -ForegroundColor Yellow
try {
    $loginResponse = Invoke-RestMethod -Uri "$baseUrl/api/v1/user/login" -Method POST -Body $loginData -ContentType "application/json"
    $headers = @{ Authorization = "Bearer $($loginResponse.access_token)" }
    
    $deleteResponse = Invoke-RestMethod -Uri "$baseUrl/api/v1/posts/$postId" -Method DELETE -Headers $headers
    Write-Host "✅ Blog Post Deleted: $($deleteResponse.message)" -ForegroundColor Green
} catch {
    Write-Host "❌ Blog Post Deletion Failed: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host "`n🎉 High Priority Features Test Complete!" -ForegroundColor Green
Write-Host "=================================================" -ForegroundColor Green

Write-Host "`n📊 Features Tested:" -ForegroundColor Cyan
Write-Host "✅ Blog Post CRUD (Create, Read, Update, Delete)" -ForegroundColor Green
Write-Host "✅ Blog Search (by title and author)" -ForegroundColor Green
Write-Host "✅ Blog Filtering (by tags and date)" -ForegroundColor Green
Write-Host "✅ Like/Dislike functionality" -ForegroundColor Green
Write-Host "✅ User Logout with token invalidation" -ForegroundColor Green