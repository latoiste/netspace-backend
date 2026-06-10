# NetSpace — Backend

Go backend for **NetSpace**, a location‑based, ephemeral social app for cafés. It serves the REST API and the realtime **WebSocket** layer (presence + chat) that the [NetSpace frontend](#-frontend) talks to, backed by **PostgreSQL**.

The core idea: a visitor scans a café's QR, must be physically **inside the venue's GPS geofence** to enter, checks in, and can chat with others at the same café. Sessions are ephemeral — a user's private data is purged on logout.

---

## ✨ Features

- 🔐 **JWT auth** issued at check‑in; protected REST + authenticated WebSocket.
- 🛰️ **Realtime hub** per café — presence (join/leave), public room, DMs (with read receipts), and group chats.
- 📍 **Per‑venue geofence data** (coordinates + radius) served to the client.
- 👥 **Group chat** — create / rename / invite (idempotent) / accept / reject / leave, with persistence.
- 🔔 **Notifications** — DM + group‑invite notifications, persisted and dismissible.
- 🧹 **Housekeeping** — public‑room message retention sweep, stale‑user cleanup on boot, chat‑data purge on logout.
- 🛠️ **Admin API** — login, analytics, active users, force‑logout, location toggle.

---

## 🧱 Tech Stack

| Area | Tech |
|---|---|
| Language | **Go 1.26** (module `github.com/latoiste/netspace`) |
| Router | **chi/v5** |
| Realtime | **gorilla/websocket** |
| Auth | **JWT** (`golang-jwt`) |
| Database | **PostgreSQL** (`lib/pq`) |
| Config | **godotenv** (`.env` for local; real env vars in prod) |

---

## 🏛️ Architecture

```
main.go → app.NewEnv() → handler.NewHandler(...) → handler.StartServer() (:8080 / $PORT)

Manager ── owns ──▶ Hub (one per café/location)
                      │  run() is a single goroutine that owns all hub state
                      │  (Clients map, groups, blocks) → no data races
                      └─▶ Client (one per WebSocket connection)
```

- **One `Hub` per location**, and each hub's `run()` goroutine is the **sole owner** of that hub's in‑memory state (connected clients, groups, blocks). All mutations are funneled through channels into `run()`.
- ⚠️ **Run a single backend instance.** Hub state is in‑memory per process, so horizontally scaling to multiple replicas would split users across hubs. Keep replicas = 1 (state survives in the DB, not across processes).

---

## 📁 Project Structure

```
main.go                # entrypoint
app/env.go             # wiring: DB, auth, manager
handler/               # HTTP handlers + StartServer (routes), WS upgrade
chat/                  # the realtime core
  manager.go  hub.go   # Manager + Hub (run loop, channels)
  client.go            # per-connection read/write pumps + inbound events
  hub_chat.go  hub_group.go  hub_typing.go  hub_block.go  hub_read.go ...
db/                    # repository (SQL) per domain (user, message, group, …)
model/                 # domain types
api/                   # DTOs + WebSocket event types (wire format)
middleware/            # auth, CORS
netspace.sql           # schema + migrations + dev seed
Dockerfile             # multi-stage build
```

---

## 🚀 Getting Started (local)

### Prerequisites
- **Go 1.26+**
- **PostgreSQL 14+**
- `psql` (to run the schema)

### 1. Create the database & schema
```bash
createdb netspace
psql "postgres://USER:PASSWORD@localhost:5432/netspace?sslmode=disable" -f netspace.sql
```
`netspace.sql` creates the tables, runs idempotent migrations, and seeds **3 demo cafés** (`kopiloka`, `koktong`, `kopi-braga`) and their admin accounts.

### 2. Configure environment
Create **`.env`** in the project root (it's git‑ignored):
```bash
DB_CONN_STRING=postgres://USER:PASSWORD@localhost:5432/netspace?sslmode=disable
JWT_SECRET_KEY=replace-with-a-long-random-secret
# PORT=8080   # optional; defaults to 8080
```

### 3. Run
```bash
go run .
# → "Server is listening on port 8080"
```

Build a binary:
```bash
go build -o netspace . && ./netspace
```

---

## ⚙️ Environment Variables

| Variable | Required | Description |
|---|---|---|
| `DB_CONN_STRING` | ✅ | PostgreSQL DSN, e.g. `postgres://user:pass@host:5432/netspace?sslmode=disable`. Add `?sslmode=require` for managed DBs that enforce TLS. |
| `JWT_SECRET_KEY` | ✅ | Secret used to sign/verify JWTs. Use a long random string; keep it out of git. |
| `PORT` | ➖ | Port to listen on. Injected by most hosts (Railway/Render); defaults to `8080`. |

> A missing `.env` is **not** fatal — in production the values come from the platform's environment.

---

## 🗄️ Database & Seed

- The full schema + seed is in **`netspace.sql`** (safe to re‑run; migrations use `ADD COLUMN IF NOT EXISTS`).
- **Seeded venues:** `kopiloka`, `koktong`, `kopi-braga` (with demo coordinates).
- **Seeded admins (DEV / DEMO ONLY):** password is `admin123` for all. Usernames: `kopiloka`, `koktong`, `kopibraga` (log in via the frontend at `/<slug>/admin/login`). **Change these in production.**

---

## 📍 Geofence — Coordinates & Radius (how to change a venue)

Each café lives in the **`Locations`** table with these geofence columns:

| Column | Meaning |
|---|---|
| `latitude` | Venue center latitude |
| `longitude` | Venue center longitude |
| `geofenceRadius` | Allowed radius in **meters** |

The client fetches them via `GET /api/locations/{slug}` and only lets a visitor in when their GPS is within `geofenceRadius` of `(latitude, longitude)`.

**Change a venue's location/radius:**
```sql
UPDATE Locations
SET latitude = -6.200754, longitude = 106.783913, geofenceRadius = 40
WHERE slug = 'kopiloka';
```

**Get a coordinate:** right‑click the spot in Google Maps (first item is `lat, lng`), or from a browser console:
```js
navigator.geolocation.getCurrentPosition(p => console.log(p.coords.latitude, p.coords.longitude));
```

**Radius guidance:** `30–50 m` for a real café (phone GPS accuracy is ~10–30 m). The frontend also has a global testing override — see the frontend README's geofence section.

---

## 🔌 API & WebSocket (overview)

### REST (selected)
| Method | Path | Auth | Purpose |
|---|---|---|---|
| `POST` | `/api/sessions/check-in` | – | Create a session (returns JWT + userId) |
| `GET` | `/api/locations/{slug}` | – | Venue info incl. geofence coords |
| `GET` | `/api/locations/{slug}/users` | ✅ | People currently at the venue |
| `GET` | `/api/locations/{slug}/public-messages` | ✅ | Recent public‑room history |
| `GET` | `/api/chats` · `/api/chats/{userId}/messages` · `/api/groups/{groupId}/messages` | ✅ | Chat list & history |
| `GET` | `/api/notifications` | ✅ | Notifications |
| `POST` | `/api/admin/login` | – | Admin login |
| `GET` | `/api/admin/dashboard/stats` · `/analytics/*` · `/locations/{slug}/users` | ✅ | Admin data |
| `POST` | `/api/admin/users/{userId}/kick` | ✅ | Force‑logout a user |

### WebSocket
Connect to `GET /ws?token=<JWT>&locationSlug=<slug>`. Messages are JSON `{ "event": string, "data": object }`. Examples:
- **Client → server:** `send_message`, `send_public_message`, `create_group`, `invite_to_group`, `rename_group`, `accept_group_invite`, `mark_read`, `block_user`, `dismiss_notification`, …
- **Server → client:** `user_joined`, `user_left`, `new_message`, `messages_read`, `new_public_message`, `group_created`, `member_joined`, `group_renamed`, `new_notification`, `force_logout`, …

The full set of event payloads lives in `api/ws_events.go` (and mirrored in the frontend's `lib/wsTypes.ts`).

---

## 🐳 Docker

A multi‑stage `Dockerfile` is included (static binary on Alpine). It reads `$PORT`.

```bash
docker build -t netspace-backend .
docker run -p 8080:8080 \
  -e DB_CONN_STRING="postgres://user:pass@host:5432/netspace?sslmode=require" \
  -e JWT_SECRET_KEY="your-secret" \
  netspace-backend
```
`.dockerignore` keeps `.env` and git metadata out of the image.

---

## ☁️ Deploy (Railway example)

1. Deploy this repo (Railway detects the `Dockerfile`).
2. Add a **PostgreSQL** plugin; set `DB_CONN_STRING = ${{Postgres.DATABASE_URL}}` (append `?sslmode=disable`/`require` if needed).
3. Set `JWT_SECRET_KEY` to a strong secret. (`PORT` is injected automatically.)
4. Run `netspace.sql` against the managed DB once.
5. Generate a public domain → use its `https://…` URL as the frontend's `NEXT_PUBLIC_API_BASE_URL` (WebSocket upgrades to `wss://` automatically).

---

## 🔐 Production / Security Notes
- **Change the seeded admin passwords** (`admin123`) and use a strong `JWT_SECRET_KEY`.
- **Single instance only** (in‑memory hub state — see Architecture).
- CORS is currently `*` and the WS accepts any origin; tighten to your frontend domain for production.
- No built‑in rate limiting yet — add one (login / check‑in / messages) before opening to the public.

---

## 📄 License

For educational / portfolio use.
