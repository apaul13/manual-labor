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

# Go get
- go get github.com/joho/godotenv
- go get github.com/jackc/pgx/v5
- go get github.com/gin-gonic/gin
- go get github.com/gin-contrib/cors
