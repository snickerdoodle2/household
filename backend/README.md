# Backend

## Wymagania

1. `golang 1.22`
1. `docker`
1. `just` - [github](https://github.com/casey/just)
1. `watchexec` - [github](https://github.com/watchexec/watchexec)
1. `migrate cli` - [github](https://github.com/golang-migrate/migrate)

## Uruchomienie

### `.env`
Uzupełnij plik `.env`:
```bash
DATABASE_ADDRESS=localhost:5432
POSTGRES_PASSWORD=<YOUR PASSWORD>
POSTGRES_DB=<JAK NAZYWAMY>
DATABASE_URL=postgresql://postgres:<YOUR PASSWORD>@localhost:5432/<JAK NAZYWAMY>
```

### Pierwsze uruchomienie
1. `docker compose up -d`
1. `just up`
1. `just run`
