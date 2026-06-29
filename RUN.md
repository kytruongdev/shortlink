# Running ShortLink

Step-by-step instructions to set up, build, and test the service. For what it is and how it
is designed, see [README.md](README.md).

All commands run from the `api/` directory:

```bash
cd api
```
---

## 1. Set up and run (Makefile)

One command builds the app image and starts PostgreSQL, applies the migrations, and starts
the service:

```bash
make up                         # build + start Postgres + migrate + start app
curl localhost:8080/healthz     # → 200 when ready
```

Lifecycle commands:

```bash
make logs                       # follow logs (Ctrl-C to stop following; app keeps running)
make ps                         # container status
make down                       # stop, keep the database volume
make reset                      # stop and wipe the database volume
```

Configuration comes from environment variables; the local stack sets working defaults, so no
setup is required for `make up`:

| Variable | Local default | Description |
|---|---|---|
| `PORT` | `8080` | HTTP port |
| `DATABASE_URL` | `postgres://shortlink:shortlink@localhost:5432/shortlink?sslmode=disable` | PostgreSQL DSN |
| `BASE_URL` | `http://localhost:8080` | Host used to build returned short URLs |
| `LOG_LEVEL` | `info` | `debug` \| `info` \| `warn` \| `error` |

---

## 2. Build

Compile the binary directly:

```bash
make build                      # outputs bin/server
```

---

## 3. Test the API

With the stack running (`make up`), exercise both endpoints.

```bash
# encode: long URL -> short URL
curl -X POST localhost:8080/api/v1/encode \
  -H 'Content-Type: application/json' \
  -d '{"long_url":"https://codesubmit.io/library/react"}'
# → {"short_url":"http://localhost:8080/GeAi9K","code":"GeAi9K"}

# decode: short URL (or bare code) -> long URL
curl -X POST localhost:8080/api/v1/decode \
  -H 'Content-Type: application/json' \
  -d '{"short_url":"http://localhost:8080/GeAi9K"}'
# → {"long_url":"https://codesubmit.io/library/react"}
```

Or explore and run the endpoints interactively in a browser via Swagger UI:

```
http://localhost:8080/swagger/index.html
```

### Run the automated test suite

```bash
make test                       # provisions Postgres via Compose, runs all tests with -race
```

---

## Deploying (production stack)

A second compose file, `build/docker-compose.prod.yml`, runs the same app for deployment:
fully environment-driven, with a persistent database volume and no published database port.
On the server, create `.env.prod` (`POSTGRES_USER`/`POSTGRES_PASSWORD`/`POSTGRES_DB`,
`APP_DATABASE_URL` pointing at the `postgres` service, `BASE_URL` set to the public address,
`LOG_LEVEL`), then start it with the `prod` environment selected:

```bash
make up ENV=prod                # also: make logs ENV=prod | make down ENV=prod
```

The live demo runs exactly this way on a GCP VM.
