<div align="center">

# 🛡️ Gfonseca Security

### Infrastructure intelligence dashboard — DNS · WHOIS · TLS · Port scan

[![Go](https://img.shields.io/badge/Go-1.24-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://go.dev/)
[![Fiber](https://img.shields.io/badge/Fiber-v2-00ACD7?style=for-the-badge&logo=go&logoColor=white)](https://gofiber.io/)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?style=for-the-badge&logo=docker&logoColor=white)](https://www.docker.com/)
[![Nmap](https://img.shields.io/badge/Nmap-Powered-7C3AED?style=for-the-badge)](https://nmap.org/)
[![License](https://img.shields.io/badge/License-MIT-22C55E?style=for-the-badge)](LICENSE)

<br />

<img src="docs/preview.png" alt="Gfonseca Security dashboard preview" width="92%" />

<br />

**Web dashboard for infrastructure reconnaissance** with a modern dark UI and real-time SSE feedback.

[![Live Demo](https://img.shields.io/badge/🌐_Demo-localhost:8080-8B5CF6?style=flat-square)](http://localhost:8080)
[![SSE](https://img.shields.io/badge/⚡_Realtime-SSE-EC4899?style=flat-square)](#-api-reference)
[![UI](https://img.shields.io/badge/🎨_UI-Dark_Glassmorphism-6366F1?style=flat-square)](#-features)

</div>

> [!WARNING]
> **Authorized use only.** Intended for permitted security testing, internal audits, and education. Only scan targets you own or have explicit written permission to analyze.

---

## 📑 Table of contents

| | |
|:---:|:---|
| ✨ | [Features](#-features) |
| 🧱 | [Tech stack](#-tech-stack) |
| 🔄 | [How it works](#-how-it-works) |
| 🚀 | [Quick start](#-quick-start) |
| 📡 | [API reference](#-api-reference) |
| 📁 | [Project structure](#-project-structure) |
| ⚠️ | [Important notices](#️-important-notices) |
| 🗺️ | [Roadmap](#️-roadmap) |
| 🤝 | [Contributing](#-contributing) |

---

## ✨ Features

<table>
<tr>
<td width="50%" valign="top">

### 🌐 DNS Intelligence
Lookup **A**, **MX**, and **NS** records in one card.

```
✓ IPv4 resolution
✓ Mail exchangers
✓ Name servers
```

</td>
<td width="50%" valign="top">

### 🗺️ WHOIS & Geo
Organization, **ASN**, country, city, and region via IP enrichment.

```
✓ AS number
✓ ISP / org name
✓ Geo location
```

</td>
</tr>
<tr>
<td valign="top">

### 🔒 SSL / TLS
Certificate issuer and **days until expiration**.

```
✓ TLS handshake
✓ Issuer CN
✓ Expiry countdown
```

</td>
<td valign="top">

### 🔌 Port Scan (Nmap)
Live **SSE stream** with progress bar and per-port cards.

```
✓ Service detection (-sV)
✓ Real-time events
✓ Open port badges
```

</td>
</tr>
</table>

<br />

| | Module | What you get |
|:---:|:---|:---|
| 🌐 | **DNS** | A, MX, and NS records |
| 🗺️ | **WHOIS / Geo** | Org, ASN, country, city, region |
| 🔒 | **SSL/TLS** | Issuer + expiration |
| 🔌 | **Port scan** | Nmap + SSE live updates |
| 🎨 | **UI** | Responsive cards, skeletons, animations |

---

## 🧱 Tech stack

<div align="center">

| Layer | Technologies |
|:---:|:---|
| ⚙️ **Backend** | ![Go](https://img.shields.io/badge/Go-00ADD8?style=flat-square&logo=go&logoColor=white) ![Fiber](https://img.shields.io/badge/Fiber-00ACD7?style=flat-square&logo=go&logoColor=white) |
| 🖥️ **Frontend** | ![HTML5](https://img.shields.io/badge/HTML-E34F26?style=flat-square&logo=html5&logoColor=white) ![CSS3](https://img.shields.io/badge/CSS-1572B6?style=flat-square&logo=css3&logoColor=white) ![JavaScript](https://img.shields.io/badge/JS-F7DF1E?style=flat-square&logo=javascript&logoColor=black) ![Lucide](https://img.shields.io/badge/Lucide-8B5CF6?style=flat-square) |
| 🔍 **Scanning** | ![Nmap](https://img.shields.io/badge/Nmap_-sV-7C3AED?style=flat-square) predefined common ports |
| 🐳 **Deploy** | ![Docker](https://img.shields.io/badge/Docker-Alpine-2496ED?style=flat-square&logo=docker&logoColor=white) multi-stage build |

</div>

---

## 🔄 How it works

```mermaid
%%{init: {'theme': 'dark', 'themeVariables': { 'primaryColor': '#8b5cf6', 'primaryTextColor': '#fff', 'primaryBorderColor': '#a78bfa', 'lineColor': '#6366f1', 'secondaryColor': '#1e1b4b', 'tertiaryColor': '#0f172a'}}}%%
flowchart TB
    subgraph Client["🖥️ Browser"]
        UI[Gfonseca Security UI]
    end

    subgraph Server["⚙️ Fiber API"]
        API[REST + SSE]
    end

    subgraph Services["🔧 Backends"]
        DNS["🌐 net.Lookup"]
        SSL["🔒 tls.Dial"]
        IPAPI["🗺️ ip-api.com"]
        Nmap["🔌 Nmap"]
    end

    UI -->|"POST /api/dns"| API
    UI -->|"POST /api/whois"| API
    UI -->|"POST /api/ssl"| API
    UI -->|"SSE /api/scan-stream"| API

    API --> DNS
    API --> SSL
    API --> IPAPI
    API --> Nmap
```

| Step | | Action |
|:---:|:---:|:---|
| **1** | 🔎 | Enter a **domain** or **IP** in the search bar |
| **2** | ⚡ | Parallel requests fill DNS, WHOIS, and SSL cards |
| **3** | 📡 | Port scan streams via **SSE** — each open port appears live |

---

## 🚀 Quick start

### 🐳 Docker *(recommended)*

> [!TIP]
> The image ships with **Nmap** and all runtime dependencies — no local install needed.

```bash
git clone https://github.com/guizeira/data_security.git
cd data_security

docker build -t gfonseca-security .
docker run --rm -p 8080:8080 gfonseca-security
```

<div align="center">

[![Open Dashboard](https://img.shields.io/badge/🚀_Open-http://localhost:8080-8B5CF6?style=for-the-badge)](http://localhost:8080)

</div>

---

### 💻 Local development

> [!NOTE]
> **Requirements:** Go **1.24+** and [Nmap](https://nmap.org/) on your `PATH`.

```bash
git clone https://github.com/YOUR_USERNAME/data_security.git
cd data_security

go mod download
go run .
```

---

## 📡 API reference

**POST** endpoints accept JSON:

```json
{
  "target": "example.com"
}
```

| Method | Endpoint | Description |
|:---:|:---|:---|
| `GET` | `/` | 🖥️ Web interface |
| `POST` | `/api/dns` | 🌐 A records, MX, NS |
| `POST` | `/api/whois` | 🗺️ Organization, ASN, location |
| `POST` | `/api/ssl` | 🔒 TLS certificate status |
| `GET` | `/api/scan-stream?target=` | 🔌 Nmap SSE stream |

<details>
<summary><b>📋 Example — DNS lookup</b></summary>

<br />

```bash
curl -s -X POST http://localhost:8080/api/dns \
  -H "Content-Type: application/json" \
  -d '{"target":"example.com"}' | jq
```

</details>

<details>
<summary><b>📋 Example — Live port scan (SSE)</b></summary>

<br />

```bash
curl -N "http://localhost:8080/api/scan-stream?target=example.com"
```

**SSE events:** `progress` · `port` · `result` · `done` · `error`

</details>

---

## 📁 Project structure

```
data_security/
│
├── 📄 main.go                 # Fiber server & routes
├── 📂 internal/handlers/      # DNS · WHOIS · SSL · SSE/Nmap
├── 📂 templates/
│   └── index.html             # Main UI
├── 📂 static/
│   ├── 🎨 css/style.css
│   └── ⚡ js/app.js
├── 📂 docs/
│   └── 🖼️ preview.png         # README screenshot
├── 🐳 Dockerfile
└── 📦 go.mod
```

---

## ⚠️ Important notices

> [!CAUTION]
> **Legal & ethical use** — Unauthorized network scanning may violate ToS or local laws. **You** are responsible for compliance.

| | Topic | Details |
|:---:|:---|:---|
| ⚖️ | **Ethics** | Get permission before scanning third-party targets |
| 📌 | **WHOIS** | IP geolocation enrichment — not classic domain WHOIS |
| 🔌 | **Ports** | Fixed common port list — not a full 1–65535 scan |
| 🌍 | **API** | Org/location via [ip-api.com](http://ip-api.com/) (non-commercial) |

---

## 🗺️ Roadmap

- [ ] 🐳 `docker-compose.yml` for one-command local setup
- [ ] ⚙️ Environment variables (port, timeout, custom port list)
- [ ] 📄 Report export (JSON / PDF)
- [ ] 📜 Scan history & persistence
- [ ] 🌙 Theme toggle (dark / light)

---

## 🤝 Contributing

Contributions are welcome! 🎉

```mermaid
gitGraph
   commit id: "fork"
   branch feat/my-feature
   checkout feat/my-feature
   commit id: "implement"
   commit id: "polish"
   checkout main
   merge feat/my-feature id: "PR merged" tag: "🎉"
```

| Step | Command |
|:---:|:---|
| **1** | Fork this repository |
| **2** | `git checkout -b feat/my-feature` |
| **3** | `git push origin feat/my-feature` |
| **4** | Open a **Pull Request** |

> [!TIP]
> Bug reports, feature ideas, and UI polish PRs are especially appreciated.

---

## 📄 License

This project is open source under the **[MIT License](LICENSE)**.

---

<div align="center">

### Built with 💜 by **Guilherme Fonseca**

[![GitHub](https://img.shields.io/badge/GitHub-@YOUR_USERNAME-181717?style=for-the-badge&logo=github)](https://github.com/YOUR_USERNAME)
[![Stars](https://img.shields.io/github/stars/YOUR_USERNAME/data_security?style=for-the-badge&logo=github&color=8B5CF6)](https://github.com/YOUR_USERNAME/data_security)

<br />

<sub>🛡️ Gfonseca Security · Infrastructure intelligence, done beautifully.</sub>

</div>
