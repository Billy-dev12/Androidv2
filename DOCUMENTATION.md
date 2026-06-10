# DOCUMENTATION.md

### 2026-06-10 — Push Update & Integrasi Firmware Extractor
- **File yang diubah:** `main.go`, `routes/router.go`, `dokumentasi.md` (dan menambahkan file baru `app/models/firmware.go`, `app/controllers/firmware_controller.go`)
- **Apa yang dilakukan:** Mengintegrasikan menu Firmware Extractor untuk ekstraksi arsip firmware Xiaomi (.zip, .tgz, .tar.gz, .tar). Melakukan build binary `android-tool` dan melakukan git push ke remote repository.
- **Mengapa:** Menyelesaikan implementasi fitur pengekstraksi firmware Xiaomi yang sebelumnya berstatus "Soon" (Under Development) sesuai perubahan kode lokal user.
- **Status:** ✅ Selesai

### 2026-06-10 — Generalisasi Outer Archive Extractor & Deteksi Konten Otomatis
- **File yang diubah:** `routes/router.go`, `app/controllers/firmware_controller.go`, `app/models/firmware.go`
- **Apa yang dilakukan:** Merefaktor fitur ekstraktor firmware dari yang sebelumnya khusus Xiaomi menjadi generik untuk semua brand HP (ZIP, TAR, TGZ, TAR.MD5). Menambahkan fitur validasi otomatis yang men-scan dan mendeteksi tipe/brand firmware berdasarkan isi berkas di dalamnya (seperti payload.bin untuk Google Pixel/OnePlus, scatter untuk MediaTek, .tar.md5/.lz4 untuk Samsung, dll).
- **Mengapa:** Menyediakan ekstraksi arsip terstandarisasi yang lebih fleksibel dengan kemampuan deteksi tipe berkas secara pintar dibanding melakukan hardcode menu per brand handphone.
- **Status:** ✅ Selesai
