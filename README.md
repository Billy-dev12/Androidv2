<div align="center">

# Android V2 Core Engine

**CLI + Interactive TUI — Manajemen Android Device & Firmware Extraction**

![Go](https://img.shields.io/badge/Go-1.26+-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![ADB](https://img.shields.io/badge/ADB-Ready-34A853?style=for-the-badge&logo=android&logoColor=white)
![License](https://img.shields.io/badge/License-MIT-FF6F00?style=for-the-badge)
![Zero Deps](https://img.shields.io/badge/Dependencies-Zero-7B2FF7?style=for-the-badge)

**100% Go Standard Library — Zero External Dependencies**

---

[Fitur](#-fitur) • [Demo](#-demo) • [Instalasi](#-instalasi) • [Cara Pakai](#-cara-pakai) • [Struktur](#-struktur-project) • [Teknologi](#-teknologi)

</div>

---

## 🚀 Fitur

### 📱 Device Management
| Perintah | Fungsi |
|----------|--------|
| `devices` / `list` | Lihat semua perangkat Android yang terhubung |
| `info` | Detail lengkap: brand, model, Android version, root/bootloader status, battery, chipset |
| `reboot` | Reboot device |

### 📦 App Management
| Perintah | Fungsi |
|----------|--------|
| `install <apk>` | Install APK ke device |
| `uninstall <package>` | Hapus package dari device |

### 📁 File Transfer
| Perintah | Fungsi |
|----------|--------|
| `push <local> <remote>` | Kirim file/folder dari PC ke device |
| `pull <remote> <local>` | Ambil file/folder dari device ke PC |
| `screenshot <output>` | Capture layar device ke PNG |

### 🔧 Hardware Diagnostics
`diagnostics` — Laporan lengkap memory, CPU, storage, display, network/signal, dan sensor via `dumpsys`

### 🧩 Firmware Extraction
Extract berbagai format firmware dengan **auto-detection brand**:

| Format | Brand Support |
|--------|---------------|
| `.zip`, `.tgz`, `.tar.gz` | Pixel, Samsung, Huawei, Oppo, Vivo, MediaTek, Motorola |
| `.tar`, `.tar.md5` | Samsung (AP/BL/CP/CSC + LZ4 decompress) |

Fitur tambahan:
- **Partition Scanner** — List `.img`/`.bin` dengan ukuran
- **Build.prop Parser** — Baca properti device dari firmware
- **Anti Zip Slip** — Proteksi directory traversal

### ⚙️ Utility
| Perintah | Fungsi |
|----------|--------|
| `env` | Cek tools lingkungan (adb, fastboot, lz4, java, python, dll) |
| `config show` | Lihat konfigurasi |
| `config set <key> <value>` | Ubah konfigurasi persisten |
| `history` | Riwayat perintah |

---

## 🎮 Demo

```
┌─────────────────────────────────────────┐
│  Android V2 Core Engine v1.0            │
│                                         │
│  [1] List Devices                       │
│  [2] Device Info                        │
│  [3] Reboot Device                      │
│  [4] Install APK                        │
│  [5] Uninstall App                      │
│  [6] Push File                          │
│  [7] Pull File                          │
│  [8] Screenshot                         │
│  [9] Diagnostics                        │
│  [10] Firmware Tools                    │
│  [11] Check Environment                 │
│  [12] Configuration                     │
│  [13] Command History                   │
│  [0] Exit                               │
│                                         │
│  Select: █                              │
└─────────────────────────────────────────┘
```

> **Dual Interface**: Jalankan tanpa argumen untuk **TUI interaktif**, atau langsung kasih perintah untuk **CLI mode**.

---

## 📦 Instalasi

### Prerequisites
- **Go 1.26+** — [Download](https://go.dev/dl/)
- **ADB** — Android Debug Bridge ([cara install](https://developer.android.com/studio/command-line/adb))

### Build dari Source

```bash
# Clone
git clone https://github.com/Billy-dev12/Androidv2.git
cd Androidv2

# Build
go build -o android-tool

# (Opsional) Install ke sistem
sudo mv android-tool /usr/local/bin/
```

### Cek Hasil Build

```bash
./android-tool env
```

---

## 🛠 Cara Pakai

### CLI Mode

```bash
# 🔍 Lihat device terhubung
./android-tool devices

# ℹ️ Info device
./android-tool info

# 📲 Install APK
./android-tool install app.apk

# 🗑 Uninstall
./android-tool uninstall com.example.app

# 📤 Push file ke /sdcard/
./android-tool push file.txt /sdcard/

# 📥 Pull file dari device
./android-tool pull /sdcard/file.txt .

# 🖼 Screenshot
./android-tool screenshot screen.png

# 🔬 Diagnostics hardware
./android-tool diagnostics

# 📦 Extract firmware
./android-tool firmware partitions ./extracted_folder

# 🏷️ Baca build.prop dari firmware
./android-tool firmware buildprop ./extracted_folder

# 🌐 Cek environment
./android-tool env

# ⚙️ Konfigurasi
./android-tool config show
./android-tool config set default_device_id XXXXXXXXXXXX

# 📋 History
./android-tool history

# ❓ Bantuan
./android-tool help
```

### Interactive TUI Mode

```bash
./android-tool
```

| Tombol | Fungsi |
|--------|--------|
| ⬆ ⬇ / W S | Navigasi menu |
| Enter | Pilih menu |
| Angka (1-9, 0) | Shortcut menu |
| Q / Esc | Kembali / Keluar |

---

## 📂 Struktur Project

```
android-tool-mvc/
├── main.go                         # Entry point
├── go.mod                          # Module definition
│
├── app/
│   ├── controllers/                # 6 controllers
│   │   ├── device_controller.go
│   │   ├── app_controller.go
│   │   ├── file_controller.go
│   │   ├── firmware_controller.go
│   │   ├── diagnostics_controller.go
│   │   └── system_controller.go
│   │
│   └── models/                     # 9 models
│       ├── adb.go                  # ADB executor
│       ├── device.go               # Device info
│       ├── application.go          # App management
│       ├── file_transfer.go        # Push/pull
│       ├── firmware.go             # Firmware extractor
│       ├── diagnostics.go          # Hardware diagnostics
│       ├── environment.go          # Tools checker
│       ├── config.go               # Persistent config
│       └── history.go              # Command history
│
├── resources/
│   └── views/
│       └── console.go              # Colored output + TUI
│
└── routes/
    └── router.go                   # CLI parser + router
```

---

## 🧰 Teknologi

| Komponen | Detail |
|----------|--------|
| **Bahasa** | Go 1.26.4 |
| **Module** | `android-tool-mvc` |
| **Dependencies** | ❌ Nol — Standard Library 100% |
| **UI** | ANSI escape codes + `stty` raw mode |
| **Arsitektur** | MVC Pattern |
| **Database** | JSON file-based config & history |
| **Keamanan** | Zip Slip protection |

---

<div align="center">

**Dibangun dengan ❤️ menggunakan Go**

[⬆ Kembali ke atas](#android-v2-core-engine)

</div>
