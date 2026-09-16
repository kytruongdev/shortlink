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

http://<host>:3000        → grafana  ← prometheus ← node_exporter (host)
http://<host>:9090        → prometheus                postgres_exporter (DB)
```

- `frontend` — `web/Dockerfile`: builds the Vite app (`VITE_API_BASE_URL=/api/v1`,
  relative) and serves it with nginx (`web/nginx.conf`). Only service with a
  published port.
- `backend` — `api/Dockerfile`: the Go API. Internal network only.
- `postgres` — data store. Internal only; data in the `pgdata` volume.
- `migrate` — applies `api/data/migrations` once, then exits.
- `node_exporter` / `postgres_exporter` — expose host hardware and PostgreSQL metrics.
- `prometheus` (`:9090`) — scrapes the exporters. `grafana` (`:3000`) — dashboards.

## Components

| File | Purpose |
|------|---------|
| `compose/docker-compose.prod.yml` | The production stack. |
| `.env.prod.example` | Template for `deploy/.env.prod` (secrets, URLs). |
| `../web/Dockerfile`, `../web/nginx.conf` | Frontend image + single-origin proxy. |
| `../api/Dockerfile` | Backend image. |

## Images & CI

`.github/workflows/build-push.yml` builds the backend (`api/`) and frontend
(`web/`) images and, on pushes to `master`, publishes them to GHCR:

- `ghcr.io/kytruongdev/shortlink-api` · `ghcr.io/kytruongdev/shortlink-web`
- tags: `latest` and `sha-<commit>`.

Pull requests build the images only (to validate the Dockerfiles), without
pushing. Locally the compose file builds the same images on the fly with
`--build`; on the VPS they are pulled from GHCR.

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
| `GRAFANA_ADMIN_USER` / `GRAFANA_ADMIN_PASSWORD` | Grafana login. |

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

## Monitoring

Prometheus scrapes two exporters and Grafana visualizes them, covering the two
things the app runs on: **hardware** and the **database**.

- **`node_exporter`** — host CPU, memory, disk, network.
- **`postgres_exporter`** — PostgreSQL up/down, connections, database size, transactions.
- **`prometheus`** — <http://localhost:9090> (targets at `/targets`).
- **`grafana`** — <http://localhost:3000> (login with `GRAFANA_ADMIN_*`). Two
  dashboards are provisioned automatically: **Node Exporter Full** (hardware) and
  **PostgreSQL Database**. The Prometheus datasource is provisioned as default.

Everything is provisioned from `deploy/prometheus/` and `deploy/grafana/`, so the
stack comes up with working dashboards and no manual clicking.

## Provisioning the VPS (Terraform)

`deploy/terraform/` creates the DigitalOcean droplet that runs the stack: a
`s-2vcpu-4gb` droplet (`sgp1`), the deploy SSH key, and a firewall (22/80/443
plus 3000/9090 for Grafana/Prometheus).

```bash
# 1. Create a deploy SSH key once (if you don't have one):
ssh-keygen -t ed25519 -f ~/.ssh/shortlink_deploy -N ""

# 2. Provide the DigitalOcean API token via env (never commit it):
export DIGITALOCEAN_TOKEN=dop_v1_xxx

# 3. Provision:
cd deploy/terraform
terraform init
terraform apply          # prints droplet_ip and an ssh_command output
```

The token is read from `DIGITALOCEAN_TOKEN`, so it stays out of every file.
State (`terraform.tfstate`) and the provider cache (`.terraform/`) are gitignored;
`.terraform.lock.hcl` is committed for reproducible provider versions. Tear the
VPS down with `terraform destroy` when not demoing to conserve credit.

## Deploying with Ansible

`deploy/ansible/` installs Docker on the droplet and deploys the stack.

```bash
cd deploy/ansible
cp vars.yml.example vars.yml              # set real secrets (gitignored)
echo -e "[shortlink]\n<DROPLET_IP>" > inventory.ini   # or from terraform output
ansible-playbook playbook.yml
```

The playbook: installs Docker + the Compose plugin, clones this repo to
`/opt/shortlink`, renders `.env.prod` (with `BASE_URL=http://<droplet_ip>`),
pulls the images from GHCR, starts the stack, and waits for the app to answer.
`inventory.ini` and `vars.yml` are gitignored (host + secrets); commit only the
`.example` files. Images must be public on GHCR, or set `ghcr_pat` (a
`read:packages` token) in `vars.yml` to pull private images.

## Continuous deployment (Jenkins)

`deploy/jenkins/` runs Jenkins on the droplet and defines the CD pipeline. While
GitHub Actions builds and pushes images (CI), Jenkins pulls them and redeploys (CD).

- **`Dockerfile`** — Jenkins LTS + Docker CLI + Compose plugin; plugins baked in
  (`plugins.txt`); setup wizard skipped.
- **`casc.yaml`** — Configuration-as-Code: admin user (from env) and a
  `shortlink-deploy` pipeline job created automatically (no manual clicking),
  reading `deploy/jenkins/Jenkinsfile` from the repo.
- **`Jenkinsfile`** — the pipeline: refresh repo → `docker compose pull` (latest
  images from GHCR) → `up -d` → smoke test.
- **`docker-compose.jenkins.yml`** — runs Jenkins with the Docker socket and
  `/opt/shortlink` mounted so the pipeline can redeploy the stack; UI on `:8080`.

Deploy it with Ansible (after `playbook.yml`):

```bash
cd deploy/ansible && ansible-playbook jenkins.yml
```

Then open `http://<droplet_ip>:8080`, log in with the Jenkins admin credentials
(`jenkins_admin_*` in `vars.yml`), and run the **shortlink-deploy** job — it pulls
the newest images and redeploys. `DEPLOY_BRANCH` (default `master`) selects which
branch the pipeline uses.

## Domain + HTTPS (Caddy)

`deploy/caddy/` puts Caddy in front of the app for automatic HTTPS. Caddy
terminates TLS (Let's Encrypt) and reverse-proxies to the frontend.

1. Point a DNS **A record** at the droplet IP (e.g., `shortlink.example.xyz → <ip>`).
2. In `deploy/ansible/vars.yml`, set `public_url: https://shortlink.example.xyz`
   and `cookie_secure: "true"`, and set the domain in `deploy/caddy/Caddyfile`.
3. Redeploy: `ansible-playbook playbook.yml`. Caddy obtains the certificate on
   first start (ports 80/443 must be open — they are in the Terraform firewall).

The frontend no longer publishes port 80 directly; Caddy fronts it on 80/443.
`BASE_URL`/`ALLOWED_ORIGINS` become the https URL and `COOKIE_SECURE=true`, so
refresh cookies are sent only over TLS.

## Known limitation — anonymous quota behind the proxy

The backend derives the client IP for the anonymous daily quota from the TCP
`RemoteAddr` and deliberately does not trust `X-Forwarded-For` (spoofable). Behind
the nginx proxy that address is the proxy's, so the 10-links/day anonymous quota
becomes effectively global rather than per-visitor. This is acceptable for the
demo (the quota is a soft anti-abuse limit and logged-in users are unaffected).
Making it per-visitor again would mean trusting `X-Forwarded-For` only from the
trusted proxy — a small, deliberate backend change left out of scope here.

## Roadmap (next phases)

- **P1** Monitoring: Prometheus + Grafana + node_exporter (host) + postgres_exporter (DB). ✅ done
- **P2** CI: GitHub Actions builds and pushes the backend/frontend images to GHCR. ✅ done
- **P3** Terraform: provision the DigitalOcean droplet. ✅ config ready (apply needs the API token)
- **P4** Ansible: install Docker and deploy this stack on the VPS. ✅ done
- **P5** Jenkins: continuous deployment (pull images + redeploy). ✅ done
- **P6** Domain + TLS (Caddy, Let's Encrypt). ✅ done
