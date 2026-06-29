# ShortLink

A URL shortening service written in Go. `POST` a long URL to `/encode` and get a short
code back; `POST` the short URL to `/decode` and get the original URL back. Mappings are
stored in PostgreSQL, so decoding keeps working across restarts.

- **Live demo:** http://34.21.205.23:8080
- **Interactive API (Swagger UI):** http://34.21.205.23:8080/swagger/index.html
- **Repository:** https://github.com/kytruongdev/shortlink

> The detailed run/build instructions live in a separate file: **[RUN.md](RUN.md)**.

---

## Table of Contents

1. [Overview](#1-overview)
2. [API](#2-api)
3. [Architecture & Design](#3-architecture--design)
4. [Running & Testing](#4-running--testing)
5. [Security & Attack Vectors](#5-security--attack-vectors)
6. [Scalability & the Collision Problem](#6-scalability--the-collision-problem)
7. [AI Usage](#7-ai-usage)

---

## 1. Overview

ShortLink exposes two JSON HTTP endpoints:

- `POST /api/v1/encode` — long URL → short URL.
- `POST /api/v1/decode` — short URL (or bare code) → long URL.

**Round-trip guarantee.** Whatever URL goes into `/encode` comes back unchanged from
`/decode`. The same URL submitted twice returns the same short code (deduplication).

**Durable by design.** A short code carries no information about its target — it is a random
7-character key. The mapping `code → url` is therefore *state*, and ShortLink keeps it in
PostgreSQL. Because the state lives in the database and not in process memory, **decoding a
previously encoded URL still works after a restart** — an architectural property, not a
feature bolted on top.

### Request lifecycle

**Encode.** Parse JSON → validate the URL (handler) → normalize it for dedup (controller)
→ look up the normalized URL; if found, return the existing code → otherwise generate a
random code and `INSERT`; on a unique-constraint conflict, regenerate and retry → return
`{short_url, code}`.

**Decode.** Parse JSON → extract the code from the short URL (or accept a bare code) →
look up by code → return the original URL, or `404` if the code is unknown.

---

## 2. API

All business endpoints live under the version prefix `/api/v1`. Every response — success
or error — is JSON.

### Endpoints

| Method | Path | Purpose |
|---|---|---|
| `POST` | `/api/v1/encode` | Shorten a long URL |
| `POST` | `/api/v1/decode` | Resolve a short URL or code back to the original |
| `GET`  | `/healthz` | Liveness probe (process is up) |
| `GET`  | `/readyz` | Readiness probe (database reachable) |
| `GET`  | `/swagger/index.html` | Interactive API documentation |

### `POST /api/v1/encode`

Request:

```json
{ "long_url": "https://codesubmit.io/library/react" }
```

Response `200`:

```json
{ "short_url": "http://34.21.205.23:8080/GeAi9K", "code": "GeAi9K" }
```

Errors: `400 INVALID_BODY` (malformed/oversized JSON or unknown field), `400 INVALID_URL`
(failed validation), `500` (code generation exhausted), `504 TIMEOUT`.

### `POST /api/v1/decode`

Accepts either a full short URL or a bare code.

Request:

```json
{ "short_url": "http://34.21.205.23:8080/GeAi9K" }
```

Response `200`:

```json
{ "long_url": "https://codesubmit.io/library/react" }
```

Errors: `400 INVALID_BODY`, `400 INVALID_SHORT_URL` (empty/invalid), `404 NOT_FOUND`
(unknown code), `504 TIMEOUT`.

> **Resolving a short URL.** A short URL is resolved by `POST /api/v1/decode`, which returns
> the original URL as JSON. There is intentionally **no** `GET /{code}` redirect endpoint —
> adding one would be outside the assignment's two-endpoint scope and would introduce an
> open-redirect surface (see §5).

### Error format

Every error is a JSON object with a stable machine-readable `code` and a human-readable
`message`:

```json
{ "code": "NOT_FOUND", "message": "short url not found" }
```

---

## 3. Architecture & Design

### Layered design

A single Go service in three layers, with PostgreSQL behind the repository:

```
HTTP → handler → controller → repository → PostgreSQL
        (transport)  (business)   (storage)
```

- **handler** — matches routes (chi), decodes/encodes JSON, performs *shape* validation of
  input, maps domain errors to HTTP status codes. No business logic.
- **controller** — the business layer: normalizes and deduplicates URLs, generates short
  codes, orchestrates the repository.
- **repository** — accesses PostgreSQL through `sqlc`-generated, type-safe code. It sits
  behind an interface so the controller does not depend on storage details and can be
  tested with a mock.

Each layer boundary is an **interface**, which is what makes the layers independently
testable (the controller is unit-tested against a mocked repository).

### Project structure

```
shortlink/
├── README.md · RUN.md                      # this file + run instructions
└── api/                                     # Go module: github.com/kytruongdev/shortlink
    ├── Dockerfile · Makefile · sqlc.yaml
    ├── build/
    │   ├── docker-compose.local.yml         # local stack (Postgres + migrate + app)
    │   └── docker-compose.prod.yml          # deployment stack (env-driven, persistent)
    ├── cmd/server/main.go                   # composition root: wires everything, one binary
    ├── data/
    │   ├── migrations/                       # golang-migrate up/down pairs (schema history)
    │   └── queries/                          # SQL source for sqlc
    └── internal/
        ├── config/                           # env config, strict (fail-fast, no defaults)
        ├── model/                            # Link struct + domain constants (code length, max URL)
        ├── controller/link/                  # business: encode/decode, dedup, retry-on-collision
        ├── repository/
        │   ├── sqlc/                          # generated, type-safe DB code (committed)
        │   └── link/                          # wraps sqlc behind a Repository interface
        ├── handler/
        │   ├── router.go                      # registers /api/v1 routes
        │   └── rest/                           # encode/decode handlers (+ Swagger annotations)
        ├── infra/
        │   ├── app/runner.go                  # generic, transport-agnostic graceful-shutdown runner
        │   ├── db/pg/                          # PostgreSQL connection pool
        │   └── httpserver/                     # transport plumbing: middleware, health, JSON, error mapping
        └── pkg/
            ├── apperror/                       # domain error carrying an HTTP status
            └── urlshortener/                   # pure helpers: generate, validate, normalize
```

The layout is **layered + per-resource**: `link` is the single resource, and each layer
package (`controller/link`, `repository/link`) co-locates its interface, implementation,
and generated mock. Helpers with no I/O live in `pkg/urlshortener` so they are trivially
unit-testable. Operational concerns (health probes, JSON helpers, error mapping) live in
`infra/httpserver` and are kept out of the business layers.

### Data model

A single table holds every mapping:

```sql
CREATE TABLE links (
    code           TEXT        PRIMARY KEY,         -- short code (base62, 7 chars)
    original_url   TEXT        NOT NULL,            -- the URL exactly as submitted
    normalized_url TEXT        NOT NULL UNIQUE,     -- canonical form, used for dedup
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

- `code` is the primary key — a fast `code → url` lookup for decode, and the uniqueness
  guarantee that backs collision handling.
- `normalized_url` is `UNIQUE` — it drives deduplication (`url → code`) and guarantees one
  code per canonical URL.
- `original_url` is stored verbatim, so decode returns exactly what was submitted.

**Concurrency & correctness: the database is the source of truth**

Deduplication and code generation are both **check-then-act** sequences, and naive
check-then-act is a classic **TOCTOU** race: under load, two requests can both pass the
"does this already exist?" check before either inserts, producing duplicate or corrupted
mappings.

The fix is to not rely on the application-level check for correctness — the `UNIQUE`
constraints above are the source of truth, and the check is only a fast path. The
authoritative step is the `INSERT`: on a unique violation (Postgres `23505`), the controller
reconciles instead of failing — a duplicate `normalized_url` means a concurrent request won,
so it re-fetches and returns that row; a duplicate `code` is a collision, so it regenerates
and retries. The database, not the application, serializes writes, so the check-to-act window
is harmless.

### Tech stack

| Concern | Choice |
|---|---|
| Language | Go |
| HTTP router | go-chi/chi v5 |
| Database | PostgreSQL 16 |
| DB access | sqlc (generated, type-safe) over pgx v5 |
| Migrations | golang-migrate (`migrate/migrate` container) |
| Mocks | mockery |
| Logging | `log/slog` (structured) |
| API docs | swaggo/swag (Swagger UI) |
| Packaging | Docker (distroless) + docker-compose |
| Tests | `testify`, table-driven |


### Design decisions & trade-offs

**Stateful mapping, not reversible encoding.**
- A reversible code (base64/encrypt) would need no storage — but by the pigeonhole principle
  a short code can't hold an arbitrary-length URL; it would have to be as long as the URL
  (JWT-style), defeating shortening.
- So the mapping is stored — which is also what lets decode survive a restart: state lives in
  the database, not in memory.

**Generating the short code: random base62 + `crypto/rand`, not a counter.**
- A database counter never collides and is fast, but its codes are *sequential* — anyone can
  walk `/aaaa1`, `/aaaa2`, … to enumerate every link — and it needs central coordination.
- Random codes can collide, so generation is backed by a DB uniqueness guarantee + bounded
  retry (§6).
- Payoff: unguessable codes, no coordination; at `62^7 ≈ 3.5×10¹²`, collisions are
  effectively zero at this scale.

**PostgreSQL via sqlc.**
- An embedded store (a JSON file or SQLite) is simpler — one file, no extra container — and
  would be reasonable here.
- Postgres is chosen for being production-shaped, for enforcing the uniqueness invariants the
  collision strategy relies on, and for matching the deployment. Cost: one more container.
- `sqlc` generates type-safe Go from plain SQL — no ORM reflection, no string-built SQL,
  which also removes the SQL-injection surface.

**Deduplication & URL normalization.**
- "Same URL" is defined by *conservative* normalization: lowercase scheme/host, drop the
  default port — so `HTTP://Example.com:80/x` and `http://example.com/x` share one code.
- The query string is left untouched on purpose, because order can carry meaning: `?a=1&b=2`
  and `?b=2&a=1` stay distinct. Over-normalizing would merge URLs that are not equivalent.

**Error handling & configuration.**
- Errors carry a client `code` + HTTP status from where they occur; the stack is captured
  once at the origin and the response is written once at the transport edge — business code
  never touches `net/http` status semantics.
- Config is read once from required env vars (fail-fast); business rules (code length 7,
  max URL 2048) are domain constants, so they can't drift between environments.

---

## 4. Running & Testing

### Quick start

With Docker installed, the whole stack (PostgreSQL + migrations + app) starts with:

```bash
cd api
make up           # builds and runs Postgres, applies migrations, starts the app
curl localhost:8080/healthz
```

For step-by-step setup, build, and API-testing instructions, see **[RUN.md](RUN.md)**.

### Tests & continuous integration

```bash
cd api
make test         # provisions Postgres via compose, then runs the full suite with -race
```

Each layer is tested at the level that fits it:

- **Both endpoints** — handler layer, with `httptest` and a mocked controller.
- **Controller** — dedup and retry-on-collision against a mocked repository.
- **Repository** — against a **real PostgreSQL** (a mock can't verify SQL/`sqlc`), isolated
  per test by a transaction rollback.

Pure helpers have table-driven tests. *Decode-after-restart* needs no test of its own — state
lives in PostgreSQL, so it survives a process restart by construction.

**CI** runs on every PR and push to `master`: `gofmt` → `go vet` → `go build` → `make test`;
a red pipeline blocks the merge.

---

## 5. Security & Attack Vectors

The strongest mitigation here is to **not build the behaviour that creates a vulnerability**.
Where a risk is simply out of scope today, the condition that would bring it back is noted so
a future change re-evaluates it.

### Input validation

`/encode` validates the URL before it touches the database (`pkg/urlshortener/validate.go`):

- **Scheme allowlist** — only `http`/`https`, which keeps `javascript:`, `data:`, and `file:`
  from being stored and later handed back to a client.
- **Userinfo rejected (anti-phishing)** — `https://www.youtube@com.vn/login` is refused.
  Without this, `url.Parse` reads `www.youtube` as *userinfo* and the real host as `com.vn`,
  producing a link that *looks* like YouTube but is not.
- **No IP-literal or `localhost` hosts**, and the host must be dot-separated labels with an
  alphabetic TLD — a public shortener points at domain names, not internal targets.
- **Length cap** — URLs over 2048 characters are rejected.

### Abuse & resource limits

- Request bodies are capped at 1 MiB (`MaxBytesReader`) and unknown JSON fields are rejected;
  every request times out at 5 s (→ `504`); a handler panic is recovered to `500`.
- SQL is parameterized through `sqlc` — there is no string-built SQL, so inputs carry no
  injection surface.
- Concurrent writes cannot duplicate or corrupt mappings — the database `UNIQUE` constraints
  enforce that (see [§3](#concurrency--correctness-the-database-is-the-source-of-truth)).

### Enumeration

Codes are random base62 from `crypto/rand`, not a sequential counter, so an attacker cannot
walk the space (`/abc0000`, `/abc0001`, …) to harvest links or predict the next code. Length
7 (over 5–6) follows the "Gone in Six Characters" research, where a denser space was scanned
to leak private documents.

### Vulnerabilities designed out

- **SSRF — not applicable.** The service never makes an outbound request to a user-supplied
  URL (no preview, no health check); it only stores and returns strings. *If* a preview
  feature were added, it would need to block private IP ranges and the cloud metadata endpoint
  (`169.254.169.254`).
- **Open redirect — not applicable.** `/decode` returns the URL as JSON; it never issues a
  `3xx` to a user-controlled location. *If* a `GET /{code}` redirect were added, it would need
  a domain allowlist or an interstitial warning page.

### Exposure & data handling

- **Request payloads are never logged** — URLs often carry secrets in the query string, so
  logs record only method, path, status, latency, and request id.
- **The database is not internet-exposed** — in the production stack it has no published host
  port and is reachable only on the internal network.
- **Secrets are not committed** — real `.env` files are git-ignored.
- **Transport** — the demo runs over plain HTTP on port 8080; a real deployment would sit
  behind a reverse proxy terminating TLS, with Swagger gated behind auth (it is exposed here
  on purpose, for grading).

---

## 6. Scalability & the Collision Problem

> The assignment is explicit that a scalable service need only be *reasoned about*, not
> *built*. So apart from collision handling, nothing here is implemented — each item names the
> issue and the threshold that would trigger work. Things that serve *scale* (cache, sharding,
> distributed IDs, multiple instances) are documented only; things a serious demo *needs*
> (Docker, tests, env config, migrations, graceful shutdown) are built for real.

### The collision problem

Random codes mean two URLs can draw the same code. This is resolved at the database, not in
application code: `code` is `UNIQUE`, so a colliding `INSERT` fails and the controller
regenerates and retries (bounded at 3 attempts) — correct even under concurrent writes (see
[§3](#concurrency--correctness-the-database-is-the-source-of-truth)). At `62^7 ≈ 3.5×10¹²` the
collision and retry rates are effectively zero.

The retry rate only rises as the table fills (≈ `N/M`). Past hundreds of millions of rows the
options are to lengthen the code, switch to distributed IDs (e.g. Snowflake, base62-encoded
with a scramble to stay unguessable), or hand out pre-allocated codes from a key-generation
service.

### Other bottlenecks, and when they would matter

None of these is built; each would be addressed only at the threshold shown.

| Concern | Today | Becomes a problem when | Approach |
|---|---|---|---|
| Read throughput | single Postgres, primary-key lookups | read QPS / p99 outgrows Postgres | Redis cache-aside for `code → url`, hot links only |
| Availability | single stateless process | HA or more throughput needed | multiple instances behind a load balancer — no code change, already stateless |
| Storage | one table | billions of rows | shard by code prefix; add read replicas |

The read path is the most likely first limit: decodes far outnumber encodes and traffic is
skewed (a viral link can dominate for minutes), which is exactly what a cache absorbs. It is
left out today only because a code lookup is already a sub-millisecond primary-key read.

Further feature-driven evolution (accounts, rate limiting, expiration) would each be an
additive migration rather than a rewrite.

---

## 7. AI Usage

The core of this project — the architecture, design decisions, business logic, data layer,
infrastructure, and deployment — is entirely my own work. AI assistance was limited to a few
supporting tasks: polishing the wording and grammar of this documentation, and speeding up
the writing of repetitive table-driven test cases.
