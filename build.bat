@echo off
setlocal
cd /d "%~dp0"

echo ========================================================
echo Membangun Split Tunnel Manager...
echo ========================================================

if not exist "bin" mkdir bin

echo Mengecek ketersediaan GCC (CGO)...
where gcc >nul 2>nul
if %errorlevel% equ 0 goto build_step

echo GCC tidak terdeteksi di PATH global! Mencoba mencari di direktori umum...
if exist "C:\Apps\w64devkit\bin\gcc.exe" set "PATH=C:\Apps\w64devkit\bin;%PATH%"
if exist "C:\w64devkit\bin\gcc.exe" set "PATH=C:\w64devkit\bin;%PATH%"
if exist "D:\w64devkit\bin\gcc.exe" set "PATH=D:\w64devkit\bin;%PATH%"
if exist "C:\TDM-GCC-64\bin\gcc.exe" set "PATH=C:\TDM-GCC-64\bin;%PATH%"
if exist "D:\msys64\mingw64\bin\gcc.exe" set "PATH=D:\msys64\mingw64\bin;%PATH%"
if exist "C:\msys64\mingw64\bin\gcc.exe" set "PATH=C:\msys64\mingw64\bin;%PATH%"
if exist "C:\msys64\ucrt64\bin\gcc.exe" set "PATH=C:\msys64\ucrt64\bin;%PATH%"

where gcc >nul 2>nul
if %errorlevel% neq 0 (
    echo [ERROR] GCC tetap tidak ditemukan! Silakan buka Terminal di VS Code dan jalankan:
    echo go build -ldflags="-s -w -H=windowsgui" -o "bin\Split Tunnel Manager.exe" .
    goto :error
)
echo GCC berhasil ditemukan dan ditambahkan ke sesi ini!

:build_step

if exist "Icon.ico" (
    echo Meng-embed icon ke Windows resource...
    go run github.com/akavel/rsrc@v0.10.2 -ico Icon.ico -o rsrc.syso >nul 2>nul
)

echo Menjalankan go build...
set CGO_ENABLED=1
set CGO_CFLAGS=-g0
set CGO_LDFLAGS=-g0
go build -trimpath -ldflags="-s -w -H=windowsgui" -o "bin\Split Tunnel Manager.exe" .

if %errorlevel% neq 0 goto :error

echo.
echo ========================================================
echo [SUKSES] Executable tersimpan di bin\Split Tunnel Manager.exe!
echo Anda wajib menjalankan "Split Tunnel Manager.exe" tersebut dengan klik kanan -^> "Run as Administrator"
echo ========================================================
pause
endlocal
exit /b 0

:error
echo.
echo ========================================================
echo [ERROR] Gagal membangun aplikasi! 
echo Pastikan environment CGO (MinGW/GCC) Anda aktif.
echo ========================================================
pause
endlocal
exit /b 1
