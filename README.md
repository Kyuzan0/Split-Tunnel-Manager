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
Pastikan Anda sudah menginstal **Go** dan compiler C (seperti **MinGW/GCC**) agar dapat melakukan *build* aplikasi Fyne CGO di Windows.

Jalankan perintah berikut di CMD atau PowerShell pada direktori proyek ini:
```bat
.\build.bat
```
Hasil kompilasi (`Split Tunnel Manager.exe`) akan otomatis tersimpan di dalam folder `bin/`.

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
Ensure you have **Go** and a C compiler (like **MinGW/GCC**) installed to build Fyne CGO applications on Windows.

Run the following command in CMD or PowerShell from the root of this project:
```bat
.\build.bat
```
The compiled executable (`Split Tunnel Manager.exe`) will be automatically saved inside the `bin/` folder.
