# Dokumentasi Android Debug Bridge (ADB) Day 1 CLI Wrapper

Aplikasi ini adalah wrapper CLI berbasis bahasa pemrograman **Go** dengan arsitektur **MVC (Model-View-Controller)** untuk memudahkan interaksi dengan Android Debug Bridge (ADB). Aplikasi dapat dijalankan dengan parameter argumen langsung maupun melalui menu interaktif (TUI/Terminal User Interface) bertingkat.

---

## 🏗️ Arsitektur Proyek (MVC Pattern)

Kode aplikasi ini distrukturkan dengan pola MVC agar mudah dikembangkan dan dipelihara:

```
Androidv2/
├── app/
│   ├── controllers/
│   │   ├── app_controller.go        # Controller untuk install/uninstall aplikasi
│   │   ├── device_controller.go     # Controller untuk status device (list, reboot, info)
│   │   └── file_controller.go       # Controller untuk transfer berkas (push/pull)
│   └── models/
│       ├── adb.go                   # Executor perintah shell ADB
│       ├── application.go           # Model logika pemasangan & penghapusan paket APK
│       ├── device.go                # Model logika pendeteksian & kontrol perangkat (termasuk info detil)
│       └── file_transfer.go         # Model logika transfer berkas (push/pull)
├── resources/
│   └── views/
│       └── console.go               # View untuk visualisasi output tabel, warna & menu TUI
├── routes/
│   └── router.go                    # Router untuk parsing CLI & state machine Menu Bertingkat
├── main.go                          # Titik masuk (entrypoint) aplikasi
└── go.mod                           # Go Module definition
```

---

## 🚀 Cara Menjalankan

### Persyaratan Sistem
1. Terpasang **Go compiler** (versi 1.16 atau lebih baru recommended).
2. Terpasang **ADB** di PATH sistem Anda (`adb` command dapat diakses lewat terminal).

### 1. Kompilasi & Build
Untuk melakukan build binary aplikasi, jalankan perintah berikut di root folder proyek:
```bash
go build -o android-tool main.go
```

### 2. Penggunaan CLI Langsung
Setelah dicompile, Anda bisa menggunakan perintah-perintah berikut:

* **Melihat Daftar Perangkat Terhubung:**
  ```bash
  ./android-tool devices
  # atau
  ./android-tool list
  ```
* **Melihat Informasi Detil Perangkat:**
  ```bash
  ./android-tool info [device-id]
  ```
* **Melakukan Reboot Device:**
  ```bash
  ./android-tool reboot [device-id]
  ```
* **Menginstal Aplikasi (APK):**
  ```bash
  ./android-tool install <apk-path> [device-id]
  ```
* **Menghapus Aplikasi (Package):**
  ```bash
  ./android-tool uninstall <package-name> [device-id]
  ```
* **Kirim File (Push):**
  ```bash
  ./android-tool push <local-path> <remote-path> [device-id]
  ```
* **Ambil File (Pull):**
  ```bash
  ./android-tool pull <remote-path> <local-path> [device-id]
  ```
* **Melihat Bantuan / Help:**
  ```bash
  ./android-tool help
  ```

---

## 🎮 Mode Menu Interaktif (TUI Bertingkat)

Jika Anda menjalankan aplikasi tanpa argumen tambahan:
```bash
./android-tool
```

Aplikasi akan memuat menu utama:
1. **ADB Management Menu** (Akan masuk ke Sub-Menu ADB)
2. **Fastboot Management Menu** (Soon)
3. **Xiaomi Firmware Extractor** (Soon)
4. **MediaTek Port Monitor (BROM)** (Soon)
5. **Exit**

### Sub-Menu ADB Management:
* **List Devices**: Menampilkan perangkat terhubung dalam format tabel.
* **Device Info**: Membaca brand, model, versi android, sdk version, dan level baterai perangkat terhubung.
* **Reboot Device**: Mengirimkan sinyal reboot ke perangkat target.
* **Install APK**: Meminta input path file APK lokal untuk dipasang ke perangkat.
* **Uninstall Package**: Meminta input nama package aplikasi untuk dihapus dari perangkat.
* **Push File**: Mengunggah berkas lokal ke dalam direktori perangkat Android.
* **Pull File**: Mengunduh berkas dari dalam direktori perangkat Android ke mesin lokal.
* **Back to Main Menu**: Kembali ke halaman menu utama.

### Cara Navigasi:
- **`Arrow Up / Down`**: Berpindah pilihan menu ke atas/bawah.
- **`Angka 1 s.d. 8`**: Memilih menu secara instan.
- **`Enter`**: Menjalankan menu yang sedang disorot.
- **`q`**: Kembali ke menu sebelumnya atau keluar (pada Menu Utama).
