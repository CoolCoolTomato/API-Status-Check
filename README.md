# API Status Check

A self-hosted API availability monitoring system. Periodically sends test requests to configured API endpoints and records availability, latency, and response previews.

## Features

- Monitor multiple API endpoints every 10 minutes
- Uses OpenAI-compatible Responses API (`POST /responses`)
- Stores full check history and maintains a rolling recent-100 log
- REST API for managing monitored endpoints
- JWT-based admin authentication
- React dashboard with per-channel availability bars

## Quick Start

### Local Development

**Backend**

```bash
# Create .env in project root
echo "ADMIN_KEY=your-secret-key" > .env

go run ./cmd/server/
# Server starts on http://localhost:8080
```

**Frontend**

```bash
cd web
pnpm install
pnpm dev
# Dev server starts on http://localhost:5173
```

### Docker

```bash
docker compose -f docker/docker-compose.yml up -d
# App available at http://localhost:8080
```

Data is persisted to `docker/data/` on the host.

## Configuration

| Variable    | Source          | Description              |
|-------------|-----------------|--------------------------|
| `ADMIN_KEY` | `.env` or env   | Admin login key for JWT  |

`.env` takes priority over environment variables.

## API Endpoints

### Auth
| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/auth/login` | Get JWT token |

### Admin (requires JWT)
| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/admin/apis` | List all API configs |
| POST | `/api/admin/apis` | Create API config |
| GET | `/api/admin/apis/:id` | Get single config |
| PUT | `/api/admin/apis/:id` | Update config |
| PATCH | `/api/admin/apis/:id` | Toggle enabled |
| DELETE | `/api/admin/apis/:id` | Delete config |

### Checks (public)
| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/checks/recent` | Recent 100 checks |
| GET | `/api/checks/history` | Full history |
| POST | `/api/checks/run` | Trigger manual check |

## Data Files

| File | Description |
|------|-------------|
| `data/apis.json` | API configurations |
| `data/check_history.json` | All check records |
| `data/recent_100_checks.json` | Latest 100 records |
