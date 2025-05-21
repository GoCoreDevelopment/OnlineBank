run:
	PORT="8080" DATABASE_URL="postgres://postgres:1234@localhost:5432/postgres?sslmode=disable" SECRET_KEY_JWT="sdfsdfweffdx12ds" go run cmd/main.go
migrate-up:
	migrate -path ./migrations -database "postgres://postgres:1234@localhost:5432/postgres?sslmode=disable" up
migrate-down:
	migrate -path ./migrations -database "postgres://postgres:1234@localhost:5432/postgres?sslmode=disable" down
