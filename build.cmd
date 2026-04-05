@echo off
echo Tidying go modules...
go mod tidy
if %ERRORLEVEL% neq 0 (
  echo Failed to tidy go modules.
  exit /b %ERRORLEVEL%
)

echo Building Go API server...
go build -v -o api.exe .\cmd\api
if %ERRORLEVEL% neq 0 (
  echo Build failed!
  exit /b %ERRORLEVEL%
)

echo Build complete! You can run it with .\api.exe
