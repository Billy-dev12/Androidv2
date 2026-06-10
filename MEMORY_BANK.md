# MEMORY_BANK.md

## Konteks Project
Aplikasi ini adalah wrapper CLI berbasis bahasa pemrograman Go dengan arsitektur MVC (Model-View-Controller) untuk memudahkan interaksi dengan Android Debug Bridge (ADB) serta beberapa tools utility tambahan seperti Firmware Extractor.

## Keputusan Desain Penting
- Menggunakan arsitektur MVC untuk memisahkan logika eksekusi ADB (Model), presentasi TUI (View), dan koordinasi input/kontrol (Controller).
- Menggunakan CLI router untuk parsing argumen masukan langsung maupun mode menu interaktif (TUI).
- Pencegahan celah keamanan seperti Zip Slip dan Directory Traversal pada model ekstraksi firmware (`app/models/firmware.go`).
- Mengimplementasikan pendekatan generik pada Firmware Extractor (Outer Archive Extractor) dengan pemindaian tipe data otomatis (Auto-Detection) setelah ekstraksi untuk fleksibilitas maksimal.
- Mengintegrasikan file `.gitignore` untuk mencegah berkas binary/arsip terkompresi/citra partisi besar terdorong ke repositori remote Git.
- Mengimplementasikan `Samsung Firmware Extractor (Inner)` untuk mengekstrak isi dari file `.tar.md5` hasil ekstraksi luar ke sub-folder masing-masing komponen (`AP`, `BL`, `CP`, `CSC`, `HOME_CSC`) secara teratur tanpa dipaksa mendekompresi berkas `.lz4` internal ke `.img`.

## Masalah yang Sedang Dikerjakan
- Selesai merefaktor Samsung Inner Extractor untuk mengekstrak konten berkas `.tar.md5` lokal ke sub-folder terpisah tanpa dekompresi otomatis `.lz4`.

## Catatan Teknis Penting
- `.gitignore` menyaring berkas `android-tool`, `extracted_*/`, serta file berformat `.zip`, `.tgz`, `.img`, `.lz4`, `.bin`, `.pkg`, `.app` secara rekursif.
- Pencarian file Samsung lokal menggunakan prefix nama berkas (`AP_`, `BL_`, `CP_`, `CSC_`, `HOME_CSC_`) dan ekstensi `.tar` / `.tar.md5`.
- Pemilihan berkas parsial Samsung mendukung parsing string yang dipisahkan koma (case-insensitive, contoh: `ap,bl`).
- Ekstraksi isi berkas `.tar.md5` dilakukan ke sub-folder terpisah (misal: `outputDir/AP/`) menggunakan `ExtractTarRaw`.

## File Kunci Project
- `main.go` — Titik masuk utama aplikasi untuk inisialisasi model, view, controller, dan router.
- `routes/router.go` — Router utama yang mem-parsing argumen CLI dan mengatur alur TUI interaktif.
- `app/models/firmware.go` — Logika dekompresi arsip firmware, penentu tipe brand (`DetectFirmwareType`), dan pencarian komponen `.tar.md5` lokal (`FindSamsungFiles`).
- `app/controllers/firmware_controller.go` — Controller untuk mengarahkan proses ekstraksi dan menyajikan hasil validasi konten ke pengguna.
- `resources/views/console.go` — View pembantu untuk menampilkan tabel device, menu interaktif, dan prompt input.

## TODO Berikutnya
- [ ] Implementasi Fastboot Management Menu (Status: Soon)
- [ ] Implementasi MediaTek Port Monitor (BROM) (Status: Soon)
- [ ] Penanganan Inner Firmware Extractor lainnya (seperti ekstraktor `payload.bin` atau dekripsi `.ozip`)
