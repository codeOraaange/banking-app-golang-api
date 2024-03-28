# banking-app

## Prerequisites
1. Installed ```golang-migrate-cli```, here is the how to https://github.com/golang-migrate/migrate/tree/master/cmd/migrate

## How to run
1. Run ```go mod tidy``` to install & remote unecessary dependencies
2. Run ```go run main.go``` the server will be start in port :8000 (you can test by call /health-check)

## Configuration .env file
1. Create .env file in root project from copy .env.example and fill the value based on your environment (development, staging, production)

## Migration guide
1. Add migration script
```
migrate create -ext sql -dir db/migrations -seq add_user_table
```

2. To execute migration
```
migrate -database "postgres://{username}:{password}@{host}:{port}/{dbname}?sslmode=disable" -path db/migrations up
```

3. To rollback migration
```
migrate -database "postgres://{username}:{password}@{host}:{port}/{dbname}?sslmode=disable" -path db/migrations down
```

4. Example
```
migrate -database "postgres://postgres:P4ssW0rd@localhost:5434/marketplace_db?sslmode=disable" -path db/migrations up
```

5. If you want to add postgres docker for local development
```
docker pull postgres
docker run --name marketplace-app -e POSTGRES_PASSWORD="P4ssW0rd" -e POSTGRES_DB="marketplace_db" -d -p 5434:5432 postgres
```
default username is 'postgres'

## Unit Test guide
1. To run the test (right now is in specific folder, i don't know how to run all yet XD)
```
go test ./test/controllers -v
go test ./test/services -v
...
```

i'll looking about this later (the most important, coverage)