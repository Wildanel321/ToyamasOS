# ToyamasOS Optimizations Specification

This document details the precise kernel, memory, network, and system service optimizations applied by ToyamasOS for **1 GB RAM / 2 vCPU VPS** targets.

---

## 1. ZRAM & Memory Optimization

On low-RAM VPS instances, traditional disk swap leads to severe disk I/O thrashing and high latency. ToyamasOS uses **ZRAM** compressed RAM swap.

- **Compression Algorithm**: `zstd` (Zstandard) or `lz4` fallback.
- **ZRAM Size**: 50% of total physical RAM (approx. 512 MB compressed buffer for 1GB VPS).
- **Swap Priority**: Priority `100` ensuring kernel uses ZRAM before disk swap.
- **Kernel Swappiness (`vm.swappiness = 100`)**: Tells the Linux kernel to aggressively compress anonymous memory pages into ZRAM while maintaining file system page caches in uncompressed RAM.

---

## 2. Kernel Sysctl Parameters (`/etc/sysctl.d/99-toyamas.conf`)

| Parameter | Value | Rationale |
| :--- | :--- | :--- |
| `vm.swappiness` | `100` | Prioritizes memory compression into ZRAM over dropping page cache. |
| `vm.vfs_cache_pressure` | `50` | Retains directory and inode caches longer in RAM to speed up file access. |
| `vm.dirty_background_ratio` | `5` | Begins background writebacks when dirty pages reach 5% of RAM to prevent disk write spikes. |
| `vm.dirty_ratio` | `10` | Caps dirty RAM at 10% before blocking processes during disk flushes. |
| `net.core.somaxconn` | `1024` | Increases maximum socket listen queue size for high HTTP concurrency. |
| `net.ipv4.tcp_max_syn_backlog` | `2048` | Expands SYN connection queue to handle sudden traffic bursts without dropping packets. |
| `net.ipv4.tcp_fin_timeout` | `15` | Reduces socket state linger time from 60s to 15s to reclaim socket memory quickly. |
| `net.ipv4.tcp_keepalive_time` | `300` | Sends keepalive probes every 5 minutes to release dead TCP connections. |
| `net.ipv4.tcp_tw_reuse` | `1` | Allows recycling TIME_WAIT sockets for outgoing connections. |
| `net.ipv4.tcp_syncookies` | `1` | Protects against TCP SYN flood Denial-of-Service attacks. |

---

## 3. Unnecessary Service Bloat Removal

Headless server environments do not require hardware management daemon processes. The installer disables and masks:

1. `bluetooth.service`: Bluetooth protocol stack.
2. `cups.service` & `cups-browsed.service`: Printing subsystem daemon.
3. `avahi-daemon.service` & `socket`: Zero-configuration mDNS networking daemon.
4. `ModemManager.service`: Cellular broadband hardware daemon.

*Saving approx. 40-70 MB RAM and reducing background CPU wakeups.*

---

## 4. Docker Engine Optimizations (`/etc/docker/daemon.json`)

To prevent container logs from filling low-capacity VPS disks:

```json
{
  "log-driver": "json-file",
  "log-opts": {
    "max-size": "10m",
    "max-file": "3"
  },
  "live-restore": true
}
```

- **Log Rotation**: Limits each container log file to 10 MB with a maximum of 3 retained files (30 MB max total per container).
- **Live Restore**: Keeps containers running during Docker engine updates.
