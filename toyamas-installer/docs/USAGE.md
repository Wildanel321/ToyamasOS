# ToyamasOS Usage & Administration Guide

## 1. Quickstart Installation

Run the main installer script with root privileges:

```bash
sudo ./install.sh
```

### Installation Flag Options

```bash
# Skip system apt update and upgrade
sudo ./install.sh --skip-update

# Run only Docker installation (Step 3)
sudo ./install.sh --step=3

# View version
sudo ./install.sh --version
```

---

## 2. Managing Core Services

### Docker & Containers

Check Docker daemon status and version:
```bash
systemctl status docker
docker info
docker compose version
```

### UFW Firewall Management

Check active firewall rules:
```bash
sudo ufw status verbose
```

Allow a new port (e.g. port 8080 for webapp):
```bash
sudo ufw allow 8080/tcp comment "My Web Application"
```

### Fail2Ban SSH Protection

Check banned IP addresses:
```bash
sudo fail2ban-client status sshd
```

Unban an IP address:
```bash
sudo fail2ban-client set sshd unbanip <IP_ADDRESS>
```

### ZRAM Memory Status

Verify ZRAM compressed memory swap status:
```bash
swapon --show
zramctl
```

Restart ZRAM service:
```bash
sudo systemctl restart toyamas-zram.service
```

---

## 3. Logs & Troubleshooting

The installer writes full execution logs to:
```
/var/log/toyamas-installer.log
```

To tail installation logs in real time:
```bash
tail -f /var/log/toyamas-installer.log
```

If a step fails during installation, the script will output the exact script file, line number, and error status code.
