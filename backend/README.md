#### Requirements

1. `golang 1.22`
2. `docker`
3. `just` - [github](https://github.com/casey/just)
4. `watchexec` - [github](https://github.com/watchexec/watchexec)
5. `migrate cli` - [github](https://github.com/golang-migrate/migrate)

#### Setup

##### `.env` Configuration
Fill in the `.env` file:
```bash
DATABASE_ADDRESS=localhost:5432
POSTGRES_PASSWORD=<DB_PASSWORD>
POSTGRES_DB=<DB_NAME>
DATABASE_URL=postgresql://postgres:<DB_PASSWORD>@localhost:5432/<DB_NAME>?sslmode=disable
```

#### First run
1. `docker compose up -d`
2. `just up`
3. `just run`
