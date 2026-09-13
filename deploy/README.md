# ShortLink — Deployment

Infrastructure and configuration to run ShortLink in production on a single VPS
with Docker Compose. This directory holds everything needed to build, ship, and
run the app; nothing here contains application business logic.

> Flow: **Localhost → Docker Desktop → VPS**. Phase 0 (this document's current
> scope) makes the full stack run as one origin on Docker Desktop. Later phases
> add monitoring, CI/CD, and infrastructure automation (see the roadmap below).

## Architecture

One reachable port (80) fronted by the frontend's nginx, which serves the SPA
and reverse-proxies the API — a single origin, so there is no CORS or
cross-site-cookie problem and the frontend bundle carries no hard-coded host.

```
http://<host>/            → frontend (nginx) → static SPA
http://<host>/api/*       → frontend (nginx) → proxy → backend (Go :8080)
http://<host>/swagger/*   → frontend (nginx) → proxy → backend
                             backend → postgres (internal only)
```

- `frontend` — `web/Dockerfile`: builds the Vite app (`VITE_API_BASE_URL=/api/v1`,
  relative) and serves it with nginx (`web/nginx.conf`). Only service with a
  published port.
- `backend` — `api/Dockerfile`: the Go API. Internal network only.
- `postgres` — data store. Internal only; data in the `pgdata` volume.
- `migrate` — applies `api/data/migrations` once, then exits.

## Components

| File | Purpose |
|------|---------|
| `compose/docker-compose.prod.yml` | The production stack. |
| `.env.prod.example` | Template for `deploy/.env.prod` (secrets, URLs). |
| `../web/Dockerfile`, `../web/nginx.conf` | Frontend image + single-origin proxy. |
| `../api/Dockerfile` | Backend image. |

## Configuration

Copy the template and fill in real values:

```bash
cp deploy/.env.prod.example deploy/.env.prod
# edit deploy/.env.prod — set strong POSTGRES_PASSWORD, JWT_SECRET, and BASE_URL
```

| Variable | Notes |
|----------|-------|
| `POSTGRES_USER` / `POSTGRES_PASSWORD` / `POSTGRES_DB` | Database credentials. |
| `DATABASE_URL` | Backend + migrate DSN; host is the service name `postgres`. |
| `BASE_URL` | Public origin used to build short URLs; the address users open. |
| `ALLOWED_ORIGINS` | CORS allow-list; same as `BASE_URL` for single-origin. |
| `LOG_LEVEL` | `debug` / `info` / `warn` / `error`. |
| `JWT_SECRET` | ≥ 32 bytes. |
| `ACCESS_TOKEN_TTL` / `REFRESH_TOKEN_TTL` | e.g. `5m` / `168h`. |
| `COOKIE_SECURE` | `false` over plain HTTP; `true` once HTTPS is enabled. |

## Run locally (Docker Desktop)

```bash
cp deploy/.env.prod.example deploy/.env.prod        # BASE_URL/ALLOWED_ORIGINS = http://localhost
docker compose -p shortlink-prod --env-file deploy/.env.prod \
    -f deploy/compose/docker-compose.prod.yml up -d --build
```

Then open <http://localhost> — shorten a URL (anonymous), register/login, manage
links, and open a short link to see the redirect. Swagger is at
<http://localhost/swagger/index.html>.

Tear down (keep data):
```bash
docker compose -p shortlink-prod -f deploy/compose/docker-compose.prod.yml down
```
Wipe data too: add `-v`.

## Known limitation — anonymous quota behind the proxy

The backend derives the client IP for the anonymous daily quota from the TCP
`RemoteAddr` and deliberately does not trust `X-Forwarded-For` (spoofable). Behind
the nginx proxy that address is the proxy's, so the 10-links/day anonymous quota
becomes effectively global rather than per-visitor. This is acceptable for the
demo (the quota is a soft anti-abuse limit and logged-in users are unaffected).
Making it per-visitor again would mean trusting `X-Forwarded-For` only from the
trusted proxy — a small, deliberate backend change left out of scope here.

## Roadmap (next phases)

- **P1** Monitoring: Prometheus + Grafana + node_exporter (host) + postgres_exporter (DB).
- **P2** CI: GitHub Actions builds and pushes the backend/frontend images to GHCR.
- **P3** Terraform: provision the AWS EC2 VPS.
- **P4** Ansible: install Docker and deploy this stack on the VPS.
- **P5** Jenkins: continuous deployment (pull image + run Ansible).
- **P6** Domain + TLS (bonus).
