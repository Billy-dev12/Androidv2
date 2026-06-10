# MEMORY_BANK.md

## Konteks Project
Aplikasi ini adalah wrapper CLI berbasis bahasa pemrograman Go dengan arsitektur MVC (Model-View-Controller) untuk memudahkan interaksi dengan Android Debug Bridge (ADB) serta beberapa tools utility tambahan seperti Firmware Extractor.

## Keputusan Desain Penting
- Menggunakan arsitektur MVC untuk memisahkan logika eksekusi ADB (Model), presentasi TUI (View), dan koordinasi input/kontrol (Controller).
- Menggunakan CLI router untuk parsing argumen masukan langsung maupun mode menu interaktif (TUI).
- Pencegahan celah keamanan seperti Zip Slip dan Directory Traversal pada model ekstraksi firmware (`app/models/firmware.go`).

## Masalah yang Sedang Dikerjakan
- Melakukan push update kode lokal yang berisi fitur Firmware Extractor baru ke repository Git terkonfigurasi.

## Catatan Teknis Penting
- Fungsi ekstraksi mendukung format `.zip`, `.tgz`, `.tar.gz`, dan `.tar`.
- Menggunakan terminal raw mode (`cbreak`, `min 1`, `-echo`) untuk mendukung navigasi menu interaktif dengan tombol panah (arrow keys) dan input angka secara langsung.

## File Kunci Project
- `main.go` — Titik masuk utama aplikasi untuk inisialisasi model, view, controller, dan router.
- `routes/router.go` — Router utama yang mem-parsing argumen CLI dan mengatur alur TUI interaktif.
- `app/models/firmware.go` — Logika dekompresi arsip firmware Xiaomi.
- `app/controllers/firmware_controller.go` — Controller untuk validasi input path firmware dan mengarahkan proses ekstraksi.
- `resources/views/console.go` — View pembantu untuk menampilkan tabel device, menu interaktif, dan prompt input.

## TODO Berikutnya
- [ ] Implementasi Fastboot Management Menu (Status: Soon)
- [ ] Implementasi MediaTek Port Monitor (BROM) (Status: Soon)
