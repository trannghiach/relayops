DB_URL=postgres://postgres:postgres@localhost:5432/relayops?sslmode=disable

migrate-up:
	go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest -path migrations -database "$(DB_URL)" up

migrate-down:
	go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest -path migrations -database "$(DB_URL)" down 1

migrate-force:
	go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest -path migrations -database "$(DB_URL)" force 1