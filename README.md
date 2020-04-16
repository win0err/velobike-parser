# Velobike Stations Parser

**[Velobike Parser](https://github.com/win0err/velobike-parser)** parses data from [velobike.ru](velobike.ru) stations map every minute.\
Parsed data contains information about stations (like address, position, etc.) and their states (total and available places). 

## Installation and running 

### Building a binary

Environment variables:
- `DB_DIALECT`: DBMS — `sqlite3` or `postgres`
- `DB_URI`: 
    - SQLite example: `/data/velobike.db`
    - PostgreSQL example: \
    `host=postgres port=5432 user=velobike password=velobike dbname=velobike sslmode=disable`
- `BACKUP_DIR`: `./data` — path to backup storage directory.

Go 1.14 required. Clone [velobike-parser](https://github.com/win0err/velobike-parser) repository and run:
```bash
$ go build -o main .
$ DB_URI="./data/velobike.db" DB_DIALECT=sqlite3 ./main
$ # or
$ BACKUP_DIR="./data" DB_URI="./data/velobike.db" DB_DIALECT=postgres ./main
```

### Docker
```bash
$ docker build -t velobike-parser .
$ docker run -d \
    --name=velobike-parser \
    -e TZ="Europe/Moscow" \
    -e DB_URI="/data/velobike.db" \
    -v $(pwd)/data:/data \
    --restart=unless-stopped \
    velobike-parser
```
### Docker Compose
Just run `docker-compose up -d` to start parser. \
By default, it uses PostgreSQL for data storage and `./data` directory as backup storage. 

You can connect to PostgreSQL database when containers is up. 

Credentials:
- **Host:** localhost
- **Port:** 54321
- **User:** velobike
- **Password:** velobike
- **Database:** velobike

See `docker-compose.yml` for extra environment variables. 

## Notes
Parser doesn't support SQLite in-memory databases.
SQLite databases doesn't support foreign key state → station.

---
_Developed by [Sergei Kolesnikov](https://github.com/win0err)_