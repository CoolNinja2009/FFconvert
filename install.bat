@echo off
echo ==============================
echo   FFConvert Installer
echo ==============================

echo.
echo [1/5] Building helper...

cd helper
go mod tidy
go build -o ffconvert-helper.exe
if errorlevel 1 (
    echo Failed to build helper.
    pause
    exit /b
)
cd ..

echo.
echo [2/5] Creating install folder...
mkdir C:\ffconvert >nul 2>nul

echo.
echo [3/5] Copying helper...
copy helper\ffconvert-helper.exe C:\ffconvert\ffconvert-helper.exe >nul

echo.
echo [4/5] Generating native host manifest...

echo {> C:\ffconvert\ffconvert-helper.json
echo   "name": "ffconvert_helper",>> C:\ffconvert\ffconvert-helper.json
echo   "description": "FFmpeg conversion helper",>> C:\ffconvert\ffconvert-helper.json
echo   "path": "C:\\ffconvert\\ffconvert-helper.exe",>> C:\ffconvert\ffconvert-helper.json
echo   "type": "stdio",>> C:\ffconvert\ffconvert-helper.json
echo   "allowed_origins": [>> C:\ffconvert\ffconvert-helper.json
echo     "chrome-extension://REPLACE_EXTENSION_ID/" >> C:\ffconvert\ffconvert-helper.json
echo   ]>> C:\ffconvert\ffconvert-helper.json
echo }>> C:\ffconvert\ffconvert-helper.json

echo.
echo [5/5] Writing registry entry...

reg add "HKCU\Software\Google\Chrome\NativeMessagingHosts\ffconvert_helper" ^
 /ve /t REG_SZ /d "C:\ffconvert\ffconvert-helper.json" /f >nul

echo.
echo Installation complete.
echo.
echo Next steps:
echo 1. Load extension from the extension folder.
echo 2. Copy its ID.
echo 3. Replace REPLACE_EXTENSION_ID in:
echo    C:\ffconvert\ffconvert-helper.json
echo.

pause
