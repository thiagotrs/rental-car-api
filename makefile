migrate-db:
	go run cmd/migration/main.go -f ./cmd/migration/sql

gen-proto-files:
	protoc --proto_path=proto --go_out=proto --go_opt=paths=source_relative *.proto

start-webserver:
	reflex -r '\.go$' -s -- sh -c "go run cmd/api/main.go"

.PHONY: migrate-db gen-proto-files start-webserver