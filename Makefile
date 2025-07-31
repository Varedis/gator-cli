test:
	go test ./...

up:
	cd sql/schema && goose postgres "$(DATABASE_URL)" up

down:
	cd sql/schema && goose postgres "$(DATABASE_URL)" down

generate:
	sqlc generate