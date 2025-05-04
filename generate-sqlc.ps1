# Check if Docker Desktop is running
$dockerProcess = Get-Process "Docker Desktop" -ErrorAction SilentlyContinue
if (-not $dockerProcess) {
    Write-Host "Docker Desktop is not running. Please start Docker Desktop and try again." -ForegroundColor Red
    exit 1
}

Write-Host "Running sqlc generate..." -ForegroundColor Cyan
docker run --rm -v "${PWD}:/app" sqlc

if ($LASTEXITCODE -eq 0) {
    Write-Host "Successfully generated SQL code!" -ForegroundColor Green
} else {
    Write-Host "Failed to generate SQL code. Please check the error messages above." -ForegroundColor Red
} 