# DOCUMENTATION.md

### 2026-06-10 — Push Update & Integrasi Firmware Extractor
- **File yang diubah:** `main.go`, `routes/router.go`, `dokumentasi.md` (dan menambahkan file baru `app/models/firmware.go`, `app/controllers/firmware_controller.go`)
- **Apa yang dilakukan:** Mengintegrasikan menu Firmware Extractor untuk ekstraksi arsip firmware Xiaomi (.zip, .tgz, .tar.gz, .tar). Melakukan build binary `android-tool` dan melakukan git push ke remote repository.
- **Mengapa:** Menyelesaikan implementasi fitur pengekstraksi firmware Xiaomi yang sebelumnya berstatus "Soon" (Under Development) sesuai perubahan kode lokal user.
- **Status:** ✅ Selesai
