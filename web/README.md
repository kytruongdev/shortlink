# ShortLink — Web (frontend)

React + TypeScript single-page app for ShortLink: shorten URLs, sign up / log in,
manage your links, and resolve short links. Built with Vite, Tailwind CSS, and
React Router. It talks to the Go API in [`../api`](../api).

## Prerequisites

- **Node.js 20.19+** (or 22.12+) and npm.
- The **backend running** (see [`../api`](../api)) — the SPA calls it over HTTP.

## Setup & run (dev)

```bash
cd web
cp .env.example .env      # then edit if your API is elsewhere
npm install
npm run dev               # http://localhost:5173
```

## Environment

| Variable | Example | Purpose |
|---|---|---|
| `VITE_API_BASE_URL` | `http://localhost:8080/api/v1` | Base URL of the API (includes the `/api/v1` prefix). |

## Running against the backend

The short URLs are served by **this app** (the frontend owns `domain/{code}`), and
the app resolves each code via the API. Two backend settings must match the
frontend origin:

- **`BASE_URL=http://localhost:5173`** — so `short_url` values point at the frontend.
  (Set in `api/build/docker-compose.local.yml` for `make up`, or `api/.env` for `make run`.)
- **`ALLOWED_ORIGINS`** includes `http://localhost:5173` — for CORS with cookies.

Then, from the repo root: start the API (`cd api && make up`), and in another
terminal start the app (`cd web && npm run dev`). Open http://localhost:5173.

> The refresh-token cookie is `SameSite=Strict` and works because `localhost:5173`
> and `localhost:8080` are the same site. Keep `COOKIE_SECURE=false` for local HTTP.

## Build

```bash
npm run build     # type-check + production build to dist/
npm run preview   # serve the built app locally
```

## Structure

```
src/
├── api/          # typed API calls (auth, links) + response types
├── auth/         # AuthContext, useAuth, ProtectedRoute, JWT decode, session bridge
├── components/   # Layout, Navbar, EncodeForm, QRImage, ui/ (Button, Input, Card, …)
├── lib/          # axios client with refresh-on-401 interceptor, helpers
└── pages/        # Home, Login, Register, MyLinks, Redirect, NotFound
```
