# Postgres Docker Setup

## Run Database
```bash
docker compose up -d
```

Stop:
```bash
docker compose down
```

Reset (clears data + re-runs init scripts):
```bash
docker compose down -v
docker compose up -d
```

## Connection
- Host: `localhost`
- Port: `5432`
- User: `{env.DB_USER}`
- Password: `{env.DB_PASSWORD}`
- DB: `manual_labor`

Connection string:
```
postgres://{env.DB_USER}:{env.DB_PASSWORD}@localhost:5432/manual_labor
```

## Initialization
SQL files in `db/init/` run automatically on first startup.

## Project Structure
```
docker-compose.yml
db/
  init/
    *.sql
```
## Run backend (development)

Option A — Run with Docker + Air (recommended for development)
1. Ensure your project `.env` has DB_USER and DB_PASSWORD (the compose service uses these to build `DB_URL`).
2. Start only the backend (build dev image that includes Air):
```
docker compose up -d --build backend
```
3. Stream logs and watch Air build/restart output:
```
docker compose logs -f backend
```
4. Edit files under `./backend` on your host — Air inside the container will detect changes, rebuild, and restart the server automatically.

Quick checks inside the running container
- Exec into the container:
```
docker compose exec backend sh
```
- Verify Air is installed and config is present:
```
which air
ls -la /app/.air.toml
```

Common troubleshooting
- If Air doesn't see file changes:
  - Verify `./backend:/app` is mounted (compose `volumes` section).
  - Ensure `.air.toml` exists at `/app/.air.toml` inside the container.
  - If permission errors occur because the container runs as a non-root user, either:
    - Temporarily run as root by adding `user: root` to the `backend` service in `docker-compose.yml` (development only), or
    - Use an entrypoint that fixes ownership on startup (I can add this if you want).
- If module downloads are slow during image build, the compose file uses `/go/pkg/mod` and `/go/bin` as cached volumes to speed iterations.

Option B — Run locally with Go (alternate)
1. Install Go on your machine.
2. From project root:
```
cd backend
go run main.go
```

Notes
- The Docker dev setup expects `backend/dockerfile.dev` (contains Air) and `backend/.air.toml` to be present. The compose file should mount `./backend:/app` so live edits are visible inside the container.
- To rebuild/recreate the backend container when making compose changes:
```
docker compose up -d --build --force-recreate backend
```


# Go get
- go get github.com/joho/godotenv
- go get github.com/jackc/pgx/v5
- go get github.com/jackc/pgx/stdlib
- go get github.com/jackc/pgx/v5/pgxpool
- go get github.com/gin-gonic/gin
- go get github.com/gin-contrib/cors
- go get github.com/stephenafamo/bob
