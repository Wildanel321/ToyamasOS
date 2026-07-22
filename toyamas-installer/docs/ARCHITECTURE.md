# ToyamasOS Architecture Overview

## Design Principles

ToyamasOS is designed around four core tenets:

1. **Zero Desktop Overhead**: Standard Debian server installations often pull unwanted printing, wireless management, or hardware discovery services. ToyamasOS disables and masks these background units.
2. **Maximum Memory Density via ZRAM**: On a 1 GB RAM VPS, running Docker containers alongside system services can easily trigger Out-Of-Memory (OOM) kernel kills. ZRAM compresses cold memory pages into RAM using fast algorithms (`zstd` / `lz4`), effectively giving the VPS ~1.5 GB usable RAM headroom without expensive disk swap latency.
3. **Container-Centric Workflow**: The primary deployment mechanism for services on ToyamasOS is Docker and Docker Compose. Docker logs are auto-rotated to protect VPS disk space.
4. **Modular & Audit-Friendly Installer**: Every phase of installation is separated into discrete, numbered Bash scripts with strict POSIX compliance and error handling (`set -euo pipefail`).

---

## Component Diagram

```
                     ┌──────────────────────────────────────┐
                     │          Debian 13 Minimal           │
                     └──────────────────┬───────────────────┘
                                        │
             ┌──────────────────────────┼──────────────────────────┐
             ▼                          ▼                          ▼
  ┌─────────────────────┐    ┌─────────────────────┐    ┌─────────────────────┐
  │   Security Stack    │    │ Memory & Kernel     │    │  Container Engine   │
  ├─────────────────────┤    ├─────────────────────┤    ├─────────────────────┤
  │ UFW Firewall        │    │ ZRAM (50% Mem Cap)  │    │ Docker CE           │
  │ Fail2Ban (SSH Jail) │    │ Sysctl (swappiness) │    │ Docker Compose      │
  └─────────────────────┘    └─────────────────────┘    └─────────────────────┘
                                        │
                                        ▼
                             ┌─────────────────────┐
                             │ Monitoring & Admin  │
                             ├─────────────────────┤
                             │ Netdata Monitoring  │
                             │ htop, curl, git     │
                             └─────────────────────┘
```

---

## Service Lifecycle Management

ToyamasOS relies on systemd for managing background processes:

- **ZRAM Lifecycle**: Initialized on boot via `toyamas-zram.service` executing `/usr/local/bin/toyamas-zram-start`.
- **Disabled Services**: Services including `bluetooth.service`, `cups.service`, `avahi-daemon.service`, and `ModemManager.service` are explicitly stopped, disabled, and masked so they cannot be triggered by dependencies.
- **Docker Lifecycle**: Docker is enabled on boot with `live-restore` enabled in `/etc/docker/daemon.json` so containers remain running during Docker engine updates.
