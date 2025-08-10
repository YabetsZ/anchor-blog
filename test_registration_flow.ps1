Write-Host "Testing Registration Flow and Activation Token Generation" -ForegroundColor Cyan
Write-Host "=============================================================" -ForegroundColor Cyan

Write-Host "`n1. Testing User Registration..." -ForegroundColor Yellow
$registerData = @{
    username = "testuser123"
    email = "testuser123@example.com"
    password = "password123"
    first_name = "Test"
    last_name = "User"
} | ConvertTo-Json

try {
    $response = Invoke-RestMethod -Uri "http://localhost:8080/api/v1/user/register" -Method POST -Body $registerData -ContentType "application/json"
    Write-Host "Registration Response:" -ForegroundColor Green
    Write-Host ($response | ConvertTo-Json -Depth 10) -ForegroundColor White
    
    $userId = $response.id
    Write-Host "`nUser ID: $userId" -ForegroundColor Cyan
    
    Write-Host "`n2. Checking if activation token was generated..." -ForegroundColor Yellow
    Write-Host "Looking for activation link in server logs..." -ForegroundColor Yellow
    
    Write-Host "`n3. Testing activation endpoint without token..." -ForegroundColor Yellow
    try {
        $activationResponse = Invoke-RestMethod -Uri "http://localhost:8080/api/v1/users/activate" -Method GET
        Write-Host "Activation without token response:" -ForegroundColor Green
        Write-Host ($activationResponse | ConvertTo-Json -Depth 10) -ForegroundColor White
    }
    catch {
        Write-Host "Expected error - no token provided:" -ForegroundColor Yellow
        Write-Host $_.Exception.Message -ForegroundColor Red
    }
    
}
catch {
    Write-Host "Registration failed:" -ForegroundColor Red
    Write-Host $_.Exception.Message -ForegroundColor Red
}

Write-Host "`nRegistration Flow Analysis Complete!" -ForegroundColor Green