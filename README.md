## Migrate db
```sh
migrate -path ./schema -database 'postgres://postgres:test@0.0.0.0:5436/postgres?sslmode=disable' up
```
