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

### 2026-06-10 — Implementasi .gitignore & Samsung Inner Extractor (Selective Component)
- **File yang diubah:** `.gitignore`, `routes/router.go`, `app/controllers/firmware_controller.go`, `app/models/firmware.go`
- **Apa yang dilakukan:** 
  1. Membuat file `.gitignore` untuk menyaring binary `android-tool`, folder ekstraksi `extracted_*`, dan file firmware besar (`.zip`, `.tgz`, `.img`, `.lz4`, dll) agar tidak terdorong ke Git.
  2. Mengimplementasikan sub-menu `Samsung Firmware Extractor (Inner)` untuk mengekstrak komponen spesifik (AP, BL, CP, CSC, HOME_CSC) dari berkas `.tar.md5` Samsung, baik secara otomatis (All) maupun pilihan spesifik yang dipisahkan koma.
  3. Mengotomatisasikan dekompresi file `.lz4` di dalam folder output menjadi berkas `.img` mentah setelah proses ekstraksi menggunakan utilitas sistem `lz4`.
- **Mengapa:** Memperkuat keamanan repositori Git dari berkas besar yang sensitif, serta memberikan kontrol ekstraksi spesifik pada firmware Samsung yang unik dan berukuran besar.
- **Status:** ✅ Selesai
