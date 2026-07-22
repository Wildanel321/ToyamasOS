# ToyamasOS Custom ISO Builder (`ToyamasOS-1.0.iso`) 💿

[![OS Target](https://img.shields.io/badge/Debian-13%20Trixie%20Minimal-red.svg)](https://www.debian.org)
[![ISO Type](https://img.shields.io/badge/Boot-UEFI%20%2B%20Legacy%20BIOS-blue.svg)](#iso-boot-specifications)
[![Tools](https://img.shields.io/badge/Tools-live--build%20%7C%20debootstrap%20%7C%20squashfs-green.svg)](#build-prerequisites)

This folder contains the complete automated ISO build pipeline for generating **`ToyamasOS-1.0.iso`** — a custom bootable Linux distribution installer based on **Debian 13 (Trixie) Minimal** pre-configured for low-resource VPS (1GB RAM & 2 vCPU).

---

## 🏗️ ISO Build Pipeline Flow

```
Debian 13 Minimal Base (debootstrap)
  └──> Toyamas Installer (ZRAM, Bloat Removal, Sysctl Tuning)
        └──> Toyamas Panel (Web Dashboard & 1-Click App Store)
              └──> Toyamas CLI (toyamas)
                    └──> Toyamas AI Assistant (Ollama & Qwen3)
                          └──> ToyamasOS-1.0.iso
```

---

## 🧰 Build Prerequisites

Building the ISO requires a Debian/Ubuntu Linux host or Virtual Machine with root privileges:

```bash
sudo apt-get update
sudo apt-get install -y live-build debootstrap squashfs-tools xorriso grub-pc-bin grub-efi-amd64-bin
```

---

## 🚀 Building `ToyamasOS-1.0.iso`

Execute the single-command automated builder script:

```bash
cd toyamas-iso
chmod +x build-iso.sh
sudo ./build-iso.sh
```

Upon completion, the pipeline produces:
- `ToyamasOS-1.0.iso`: Hybrid bootable ISO image (~650 MB).
- `ToyamasOS-1.0.iso.sha256`: SHA256 integrity checksum file.

---

## 💾 Installation & Booting Instructions

### 1. Flash to USB Drive (Bare Metal / VPS)
```bash
sudo dd if=ToyamasOS-1.0.iso of=/dev/sdX bs=4M status=progress conv=fsync
```
*(Replace `/dev/sdX` with your actual USB flash drive block device).*

### 2. Boot in Virtual Machines (Proxmox / QEMU / VirtualBox)
Attach `ToyamasOS-1.0.iso` as a virtual CD-ROM drive and select standard boot. Supports both **UEFI** and **Legacy BIOS** boot modes.
