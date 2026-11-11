include .env
export

migrate-up:
	goose up

migrate-down:
	goose down

migrate-status:
	goose status

sqlc:
	sqlc generate