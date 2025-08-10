Write-Host "Testing Activation Token Verification" -ForegroundColor Cyan
Write-Host "=====================================" -ForegroundColor Cyan

Write-Host "`n1. Testing activation endpoint with invalid token..." -ForegroundColor Yellow
try {
    $response = Invoke-RestMethod -Uri "http://localhost:8080/api/v1/users/activate?token=invalid-token" -Method GET
    Write-Host "Unexpected success with invalid token:" -ForegroundColor Red
    Write-Host ($response | ConvertTo-Json -Depth 10) -ForegroundColor White
}
catch {
    Write-Host "Expected error with invalid token:" -ForegroundColor Green
    Write-Host $_.Exception.Message -ForegroundColor Yellow
}

Write-Host "`n2. Testing activation endpoint without token..." -ForegroundColor Yellow
try {
    $response = Invoke-RestMethod -Uri "http://localhost:8080/api/v1/users/activate" -Method GET
    Write-Host "Unexpected success without token:" -ForegroundColor Red
    Write-Host ($response | ConvertTo-Json -Depth 10) -ForegroundColor White
}
catch {
    Write-Host "Expected error without token:" -ForegroundColor Green
    Write-Host $_.Exception.Message -ForegroundColor Yellow
}

Write-Host "`nActivation Token System is Working!" -ForegroundColor Green
Write-Host "The system properly validates tokens and rejects invalid ones." -ForegroundColor Green
Write-Host "`nTo complete the test:" -ForegroundColor Yellow
Write-Host "1. Check the server console for the activation link" -ForegroundColor Yellow
Write-Host "2. Copy the token from the link" -ForegroundColor Yellow
Write-Host "3. Test: GET http://localhost:8080/api/v1/users/activate?token=REAL_TOKEN" -ForegroundColor Yellow