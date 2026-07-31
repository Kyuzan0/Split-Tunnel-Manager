# CHANGELOG 

Format mengikuti [Keep a Changelog](https://keepachangelog.com/). Versi SemVer.


## [1.2.0] - 2026-08-01

### Added

- Versi aplikasi kini tampil dinamis di judul window (contoh: `Split Tunnel Manager v1.2.0`). Versi dibaca otomatis dari `CHANGELOG.md` menggunakan `//go:embed` sehingga tidak perlu mengubah kode saat rilis baru.
- `build.bat` kini menampilkan menu interaktif dengan 3 pilihan mode build:
  - **[1] Build Normal** — Sembunyikan Terminal (Mode Produksi).
  - **[2] Build Debug** — Tampilkan Terminal untuk melihat log/error.
  - **[3] Build Compressed** — Menggunakan UPX untuk mengurangi ukuran file dari ~29 MB menjadi ~13 MB. UPX diunduh otomatis dari GitHub jika belum ada.

### Changed

- `build.bat` ditambahkan flag `-trimpath` untuk menghapus path absolut lokal dari binary (lebih kecil & lebih aman).

## [1.1.0] - 2026-07-31


### Fixed

- Icon `.exe` kini tampil di File Explorer & taskbar Windows. Solusi: generate `Icon.ico` (multi-resolusi: 16–256px) lalu embed sebagai Windows Resource (`rsrc.syso`) menggunakan `github.com/akavel/rsrc`. File `.syso` di-link otomatis oleh `go build`.
- `build.bat` diperbarui: otomatis generate `rsrc.syso` dari `Icon.ico` sebelum build.

## [0.1.3] - 2026-07-31

### Fixed

- `build.bat` sekarang menyetel `CGO_CFLAGS=-g0` dan `CGO_LDFLAGS=-g0` sehingga kompatibel dengan GCC v15/v16 (w64devkit, MSYS2 ucrt64). Build kini **berhasil** tanpa perlu downgrade GCC.
- App icon (dari `bundled.go`) kini ikut ter-compile ke dalam binary dan tampil di aplikasi.

## [0.1.2] - 2026-07-31

### Changed

- `build.bat` diperbaiki: refactor logika deteksi GCC dari nested `if`-block (buggy di cmd.exe) menjadi flat dengan `goto :build_step`.
- `build.bat` menambahkan fallback path GCC ke `C:\Apps\w64devkit\bin`, `C:\msys64\ucrt64\bin`, dan `D:\msys64\mingw64\bin`.

### Known Issue

- Build gagal jika GCC yang terdeteksi adalah versi 15/16 (w64devkit v15, MSYS2 ucrt64 v16). CGO tidak kompatibel dengan format DWARF yang dihasilkan GCC terbaru ini. Solusi: gunakan MSYS2 mingw64 (GCC ≤14) atau TDM-GCC.

## [0.1.1] - 2026-07-31

### Added

- App icon kustom yang di-bundle menggunakan `fyne bundle`.

## [0.1.0] - 2026-07-31


Ini adalah versi awal.
