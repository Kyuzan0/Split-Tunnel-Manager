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

echo.
echo Silakan pilih mode kompilasi:
echo [1] Build Normal     (Sembunyikan Terminal - Mode Produksi)
echo [2] Build Debug      (Tampilkan Terminal - Mode Pengembangan)
echo [3] Build Compressed (Menggunakan UPX - Ukuran File Sangat Kecil)
choice /c 123 /n /m "Masukkan pilihan Anda [1,2,3]: "

if errorlevel 3 goto build_upx
if errorlevel 2 goto build_debug
if errorlevel 1 goto build_normal

:build_upx
echo.
echo Menjalankan go build [Mode UPX COMPRESSED]...
set CGO_ENABLED=1
set CGO_CFLAGS=-g0
set CGO_LDFLAGS=-g0
go build -trimpath -ldflags="-s -w -H=windowsgui" -o "bin\Split Tunnel Manager.exe" .
if %errorlevel% neq 0 goto :error
echo Mengecek UPX...
if not exist "tools\upx.exe" (
    echo Mengunduh UPX dari GitHub...
    if not exist "tools" mkdir tools
    powershell -Command "Invoke-WebRequest -Uri 'https://github.com/upx/upx/releases/download/v4.2.4/upx-4.2.4-win64.zip' -OutFile 'tools\upx.zip'"
    powershell -Command "Expand-Archive -Path 'tools\upx.zip' -DestinationPath 'tools' -Force"
    copy "tools\upx-4.2.4-win64\upx.exe" "tools\upx.exe" >nul
    del "tools\upx.zip"
    rmdir /s /q "tools\upx-4.2.4-win64"
)
echo Mengkompresi menggunakan UPX...
"tools\upx.exe" -9 "bin\Split Tunnel Manager.exe"
goto build_done

:build_debug
echo.
echo Menjalankan go build [Mode DEBUG]...
set CGO_ENABLED=1
set CGO_CFLAGS=-g0
set CGO_LDFLAGS=-g0
go build -trimpath -ldflags="-s -w" -o "bin\Split Tunnel Manager.exe" .
if %errorlevel% neq 0 goto :error
goto build_done

:build_normal
echo.
echo Menjalankan go build [Mode NORMAL]...
set CGO_ENABLED=1
set CGO_CFLAGS=-g0
set CGO_LDFLAGS=-g0
go build -trimpath -ldflags="-s -w -H=windowsgui" -o "bin\Split Tunnel Manager.exe" .
if %errorlevel% neq 0 goto :error
goto build_done

:build_done
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
