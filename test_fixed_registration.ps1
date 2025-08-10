Write-Host "Testing FIXED Registration Flow with Activation Token Generation" -ForegroundColor Cyan
Write-Host "=================================================================" -ForegroundColor Cyan

Write-Host "`n1. Testing User Registration with Activation Token..." -ForegroundColor Yellow
$registerData = @{
    username = "newuser456"
    email = "newuser456@example.com"
    password = "password123"
    first_name = "New"
    last_name = "User"
} | ConvertTo-Json

try {
    $response = Invoke-RestMethod -Uri "http://localhost:8080/api/v1/user/register" -Method POST -Body $registerData -ContentType "application/json"
    Write-Host "Registration Response:" -ForegroundColor Green
    Write-Host ($response | ConvertTo-Json -Depth 10) -ForegroundColor White
    
    $userId = $response.id
    Write-Host "`nUser ID: $userId" -ForegroundColor Cyan
    
    Write-Host "`n2. Check server logs for activation link..." -ForegroundColor Yellow
    Write-Host "Look for activation link in the server console window!" -ForegroundColor Yellow
    
    Write-Host "`n3. Testing with a token (example)..." -ForegroundColor Yellow
    Write-Host "If you see an activation link in server logs, copy the token and test:" -ForegroundColor Yellow
    Write-Host "GET http://localhost:8080/api/v1/users/activate?token=YOUR_TOKEN_HERE" -ForegroundColor Cyan
    
}
catch {
    Write-Host "Registration failed:" -ForegroundColor Red
    Write-Host $_.Exception.Message -ForegroundColor Red
}

Write-Host "`nFixed Registration Flow Test Complete!" -ForegroundColor Green
Write-Host "Check the server console for activation link!" -ForegroundColor Yellow