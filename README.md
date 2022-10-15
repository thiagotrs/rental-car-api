<h1 align="center">🚙 Rental Car API</h1>

## About

Rental Car API for management orders, cars, stations and price.

## Technologies

- Golang
- Modular Monolith
- DDD (principles)
- Clean Architecture (principles)

## Run Project

### Clone Project

```git
git clone https://github.com/thiagotrs/rental-car-api.git
```

### Database

```shell
docker run -d --name db -p 5432:5432 -e POSTGRES_PASSWORD=admin1234 postgres
```

### Migration
```shell
cd rental-car-api
go run cmd/migration/main.go -f ./cmd/migration/sql/logistics_db_up.sql
go run cmd/migration/main.go -f ./cmd/migration/sql/pricing_db_up.sql
go run cmd/migration/main.go -f ./cmd/migration/sql/rental_db_up.sql
```

### API

```shell
cd rental-car-api
go run cmd/api/main.go
```

## Configuration

Create 'config.yaml' file in the root directory with these variables:

```yaml
server:
  host: localhost
  port: 4000

database:
  type: postgres
  user: postgres
  pass: admin1234
  host: localhost
  port: 5432
  name: 
```

### URLs

```
http://localhost:4000/
```

## Author

Thiago Rotondo Sampaio - [GitHub](https://github.com/thiagotrs) / [Linkedin](https://www.linkedin.com/in/thiago-rotondo-sampaio) / [Email](mailto:thiagorot@gmail.com)

## License

This project use MIT license, see the file [LICENSE](./LICENSE.md) for more details

---

<p align="center">Develop by <a href="https://github.com/thiagotrs">Thiago Rotondo Sampaio</a></p>