Write-Host "Running Unit Tests for Registration Activation System" -ForegroundColor Cyan
Write-Host "=======================================================" -ForegroundColor Cyan

Write-Host "`n🧪 Running UserHandler Registration Tests..." -ForegroundColor Yellow
go test -v ./api/handler/user -run TestUserHandler_Register

Write-Host "`n🧪 Running ActivationService Tests..." -ForegroundColor Yellow  
go test -v ./internal/service/user -run TestActivationService

Write-Host "`n🧪 Running UserServices Tests..." -ForegroundColor Yellow
go test -v ./internal/service/user -run TestUserServices

Write-Host "`n🧪 Running Integration Tests..." -ForegroundColor Yellow
go test -v ./test/integration -run TestRegistration

Write-Host "`n📊 Running All Tests with Coverage..." -ForegroundColor Yellow
go test -v -cover ./api/handler/user ./internal/service/user ./test/integration

Write-Host "`n🎉 Test Run Complete!" -ForegroundColor Green
Write-Host "Check the output above for any failures or issues." -ForegroundColor Green