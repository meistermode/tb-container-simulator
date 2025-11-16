@echo off
echo building...
cd /d "%~dp0"
go build -o container-simulator.exe
if %errorlevel% equ 0 (
    echo.
    echo build successful!.
) else (
    echo.
    echo build failed!
)
pause
