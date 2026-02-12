@echo off
echo Removing FFConvert...

reg delete "HKCU\Software\Google\Chrome\NativeMessagingHosts\ffconvert_helper" /f >nul
rmdir /s /q C:\ffconvert

echo Uninstalled.
pause
