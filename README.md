# Split Tunnel Manager - Split Tunnel Manager

[🇮🇩 Bahasa Indonesia](#bahasa-indonesia) | [🇬🇧 English](#english)

---

<a name="bahasa-indonesia"></a>
## 🇮🇩 Bahasa Indonesia

**Split Tunnel Manager Split Tunnel Manager** adalah aplikasi desktop modern (dibangun menggunakan Go & Fyne) yang dirancang khusus untuk menambahkan fitur *Split Tunneling* pada layanan VPN komersial yang tidak memiliki fitur tersebut secara bawaan.

Banyak layanan VPN populer seperti **HMA (HideMyAss)** dan **Avira Phantom VPN** akan merutekan seluruh koneksi internet Anda melalui server VPN, yang seringkali memutus akses Anda ke perangkat jaringan lokal (LAN) seperti printer lokal, server NAS, atau komputer lain di jaringan yang sama. 

Split Tunnel Manager memecahkan masalah ini dengan cara memanipulasi *Route Table* Windows secara instan. Anda dapat menentukan IP atau rentang jaringan (CIDR) mana saja yang ingin Anda "Bypass" (kecualikan) dari terowongan VPN.

### Fitur Utama
- 🚀 **Bypass LAN / IP Spesifik:** Tetap bisa mengakses printer atau server lokal meski VPN sedang menyala.
- 🎨 **Modern Fyne UI:** Antarmuka yang intuitif, bersih, dan dilengkapi dengan *Dark Theme*.
- 🛡️ **Auto-Elevation (UAC):** Meminta hak akses Administrator secara otomatis untuk mengubah tabel rute Windows secara transparan.
- ⚡ **Silent Execution:** Modifikasi rute terjadi murni di latar belakang tanpa memunculkan jendela *command prompt* yang mengganggu.

### Cara Penggunaan
1. Hubungkan / Connect VPN Anda (contoh: Avira Phantom atau HMA) melalui aplikasi resminya.
2. Buka **Split Tunnel Manager**.
3. Di tab **Bypass Rules**, tambahkan rentang IP yang tidak ingin masuk ke VPN (misalnya `192.168.1.0/24`).
4. Di tab **Split Tunnel**, pilih Adapter VPN Anda dari *dropdown*.
5. Klik **Start Split Tunnel**. Selesai! Koneksi LAN Anda akan kembali normal.

### Cara Kompilasi (Build)
Pastikan Anda sudah menginstal **Go** dan compiler C (**MSYS2 ucrt64** atau **TDM-GCC**) agar dapat melakukan *build* aplikasi Fyne CGO di Windows.

> **Catatan GCC:** GCC v15/v16 (w64devkit terbaru) tidak kompatibel dengan CGO secara langsung. Script build sudah menyetel `CGO_CFLAGS=-g0` sebagai workaround.

Jalankan perintah berikut di CMD atau PowerShell pada direktori proyek ini:
```bat
.\build.bat
```
Hasil kompilasi (`Split Tunnel Manager.exe`) akan otomatis tersimpan di dalam folder `bin/`.

### Mengubah Icon Aplikasi

Icon aplikasi terdiri dari tiga file yang bekerja bersama:

| File | Keterangan |
|------|------------|
| `Icon.png` | Sumber gambar asli (bisa diubah bebas) |
| `Icon.ico` | Versi multi-resolusi (16px–256px), di-generate dari `Icon.png` |
| `rsrc.syso` | Windows Resource file, di-embed otomatis ke `.exe` saat build |
| `bundled.go` | Asset Go untuk icon di title bar window (di-generate oleh `fyne bundle`) |

Untuk mengganti icon, ikuti langkah berikut:

**1. Siapkan `Icon.png` baru** (ganti file yang ada di root proyek).

**2. Generate `Icon.ico`** menggunakan helper yang disertakan:
```powershell
# Jalankan dari root direktori proyek
go build -o make_ico.exe .\make_ico\; .\make_ico.exe; Remove-Item make_ico.exe
```

**3. Generate `rsrc.syso`** (Windows Resource untuk icon `.exe`):
```powershell
go run github.com/akavel/rsrc@v0.10.2 -ico Icon.ico -o rsrc.syso
```

**4. Update `bundled.go`** (icon untuk title bar window):
```powershell
go run fyne.io/fyne/v2/cmd/fyne@latest bundle -o bundled.go Icon.png
```

**5. Build ulang** aplikasi:
```bat
.\build.bat
```

> Langkah 3 dan 4 juga dijalankan otomatis oleh `build.bat` jika `Icon.ico` sudah ada di direktori.

---

<a name="english"></a>
## 🇬🇧 English

**Split Tunnel Manager** is a modern desktop application (built with Go & Fyne) specifically designed to add *Split Tunneling* capabilities to commercial VPN services that do not support this feature natively.

Many popular VPN services like **HMA (HideMyAss)** and **Avira Phantom VPN** route your entire internet connection through their servers, which often breaks your access to local network (LAN) devices such as local printers, NAS servers, or other computers on the same subnet.

Split Tunnel Manager solves this problem by instantly manipulating the Windows *Route Table*. You can define which specific IPs or network ranges (CIDR) you want to "Bypass" (exclude) from the VPN tunnel.

### Key Features
- 🚀 **LAN / Specific IP Bypass:** Retain access to your local printers or servers even while the VPN is connected.
- 🎨 **Modern Fyne UI:** An intuitive, clean interface featuring a custom *Dark Theme*.
- 🛡️ **Auto-Elevation (UAC):** Automatically requests Administrator privileges required to seamlessly modify Windows route tables.
- ⚡ **Silent Execution:** Route modifications happen completely in the background without spawning disruptive command prompt windows.

### How to Use
1. Connect your VPN (e.g., Avira Phantom or HMA) using its official client.
2. Open **Split Tunnel Manager**.
3. In the **Bypass Rules** tab, add the IP ranges you do not want routed through the VPN (e.g., `192.168.1.0/24`).
4. In the **Split Tunnel** tab, select your VPN's Network Adapter from the dropdown.
5. Click **Start Split Tunnel**. That's it! Your LAN connection will be restored.

### How to Compile (Build)
Ensure you have **Go** and a C compiler (**MSYS2 ucrt64** or **TDM-GCC**) installed to build Fyne CGO applications on Windows.

> **GCC Note:** GCC v15/v16 (latest w64devkit) is not directly compatible with CGO. The build script already sets `CGO_CFLAGS=-g0` as a workaround.

Run the following command in CMD or PowerShell from the root of this project:
```bat
.\build.bat
```
The compiled executable (`Split Tunnel Manager.exe`) will be automatically saved inside the `bin/` folder.

### Changing the App Icon

The application icon consists of three files working together:

| File | Description |
|------|-------------|
| `Icon.png` | Original image source (replace freely) |
| `Icon.ico` | Multi-resolution version (16px–256px), generated from `Icon.png` |
| `rsrc.syso` | Windows Resource file, automatically embedded into `.exe` at build time |
| `bundled.go` | Go asset for the window title bar icon (generated by `fyne bundle`) |

To replace the icon, follow these steps:

**1. Replace `Icon.png`** with your new image (in the project root).

**2. Generate `Icon.ico`** using the included helper:
```powershell
# Run from the project root directory
go build -o make_ico.exe .\make_ico\; .\make_ico.exe; Remove-Item make_ico.exe
```

**3. Generate `rsrc.syso`** (Windows Resource for the `.exe` icon):
```powershell
go run github.com/akavel/rsrc@v0.10.2 -ico Icon.ico -o rsrc.syso
```

**4. Update `bundled.go`** (icon for the window title bar):
```powershell
go run fyne.io/fyne/v2/cmd/fyne@latest bundle -o bundled.go Icon.png
```

**5. Rebuild** the application:
```bat
.\build.bat
```

> Steps 3 and 4 are also run automatically by `build.bat` if `Icon.ico` is present in the directory.
