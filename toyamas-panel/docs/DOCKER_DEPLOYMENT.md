# Toyamas Panel - Docker Deployment Guide

This guide details how to deploy **Toyamas Panel** using Docker and Docker Compose on a Debian 13 minimal server or VPS.

---

## 🚀 Quick Deployment with Docker Compose

### 1. Prerequisites
Ensure Docker and Docker Compose are installed (e.g. using `toyamas-installer`).

### 2. Clone Repository & Navigate
```bash
git clone https://github.com/Wildanel321/ToyamasOS.git
cd ToyamasOS/toyamas-panel
```

### 3. Launch Panel Container
```bash
docker compose up -d --build
```

Verify container status:
```bash
docker compose ps
```

Access the dashboard at `http://<YOUR_VPS_IP>:8080`.

**Default Login Credentials**:
- **Username**: `admin`
- **Password**: `toyamas123`

---

## ⚙️ Environment Variables

Customize environment variables in `docker-compose.yml` or via a `.env` file:

| Variable | Default Value | Description |
| :--- | :--- | :--- |
| `PORT` | `8080` | Internal and external HTTP server port |
| `DB_PATH` | `/app/data/toyamas.db` | SQLite database file location |
| `ADMIN_USER` | `admin` | Initial administrator username |
| `ADMIN_PASS` | `toyamas123` | Initial administrator password |

---

## 🔒 Security Best Practices

### Customizing Admin Credentials
Before deploying in production, update default credentials in `docker-compose.yml`:
```yaml
environment:
  - ADMIN_USER=my_secure_admin
  - ADMIN_PASS=StrongCustomPassword99!
```

### Mount Points Rationale
1. `/var/run/docker.sock:/var/run/docker.sock`: Allows Toyamas Panel to list, start, stop, and restart Docker containers.
2. `./data:/app/data`: Persists SQLite user accounts, session tokens, and audit logs across container restarts.

---

## 🌐 Nginx Reverse Proxy with SSL (Optional)

To serve Toyamas Panel over HTTPS with Let's Encrypt SSL on port 443:

```nginx
server {
    listen 80;
    server_name panel.yourdomain.com;
    return 301 https://$host$request_uri;
}

server {
    listen 443 ssl http2;
    server_name panel.yourdomain.com;

    ssl_certificate /etc/letsencrypt/live/panel.yourdomain.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/panel.yourdomain.com/privkey.pem;

    location / {
        proxy_pass http://127.0.0.1:8080;
        proxy_http_version 1.1;

        # WebSocket support
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "Upgrade";

        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```
