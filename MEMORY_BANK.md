# MEMORY_BANK.md

## Konteks Project
Aplikasi ini adalah wrapper CLI berbasis bahasa pemrograman Go dengan arsitektur MVC (Model-View-Controller) untuk memudahkan interaksi dengan Android Debug Bridge (ADB) serta beberapa tools utility tambahan seperti Firmware Extractor.

## Keputusan Desain Penting
- Menggunakan arsitektur MVC untuk memisahkan logika eksekusi ADB (Model), presentasi TUI (View), dan koordinasi input/kontrol (Controller).
- Menggunakan CLI router untuk parsing argumen masukan langsung maupun mode menu interaktif (TUI).
- Pencegahan celah keamanan seperti Zip Slip dan Directory Traversal pada model ekstraksi firmware (`app/models/firmware.go`).
- Mengimplementasikan pendekatan generik pada Firmware Extractor (Outer Archive Extractor) dengan pemindaian tipe data otomatis (Auto-Detection) setelah ekstraksi untuk fleksibilitas maksimal.

## Masalah yang Sedang Dikerjakan
- Selesai mengimplementasikan Outer Archive Extractor (.zip, .tar, .tgz, .tar.md5) dengan validasi isi firmware otomatis.

## Catatan Teknis Penting
- Fungsi ekstraksi mendukung format `.zip`, `.tgz`, `.tar.gz`, `.tar`, dan `.tar.md5`.
- Setelah ekstraksi berhasil, aplikasi memindai struktur file untuk mendeteksi tipe firmware (Google Pixel, OnePlus, Xiaomi, Samsung, Huawei, Oppo, Realme, Vivo, iQOO, Motorola, MediaTek Scatter).
- Menggunakan terminal raw mode (`cbreak`, `min 1`, `-echo`) untuk mendukung navigasi menu TUI.

## File Kunci Project
- `main.go` — Titik masuk utama aplikasi untuk inisialisasi model, view, controller, dan router.
- `routes/router.go` — Router utama yang mem-parsing argumen CLI dan mengatur alur TUI interaktif.
- `app/models/firmware.go` — Logika dekompresi arsip firmware serta penentu tipe brand firmware (`DetectFirmwareType`).
- `app/controllers/firmware_controller.go` — Controller untuk mengarahkan proses ekstraksi dan menyajikan hasil validasi konten ke pengguna.
- `resources/views/console.go` — View pembantu untuk menampilkan tabel device, menu interaktif, dan prompt input.

## TODO Berikutnya
- [ ] Implementasi Fastboot Management Menu (Status: Soon)
- [ ] Implementasi MediaTek Port Monitor (BROM) (Status: Soon)
- [ ] Penanganan Inner Firmware Extractor (seperti dekompresi otomatis `.lz4` Samsung ke `.img`, ekstraktor `payload.bin`, atau dekripsi `.ozip`)
