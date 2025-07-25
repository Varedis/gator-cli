test:
	go test ./...

up:
	cd sql/schema && goose postgres "postgres://robscott:@localhost:5432/gator" up

down:
	cd sql/schema && goose postgres "postgres://robscott:@localhost:5432/gator" down

generate:
	sqlc generate