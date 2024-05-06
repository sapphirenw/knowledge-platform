include $(shell sed 's/export //g' .env > .env.make && echo .env.make)
export

sql:
	$(MAKE) -C ../database -f ../database/Makefile db-schema
	awk 'FNR==1{print ""}1' ../database/queries/* > ../database/queries.sql
	sqlc generate
	rm ../database/queries.sql

run: sql
	go run ./src/*.go

# tests for various packages
test-database: sql
	go test -v -coverprofile coverage.out ./src/database/*.go
test-docstore: sql
	go test -v -coverprofile coverage.out ./src/docstore/*.go
test-utils: sql
	go test -v -coverprofile coverage.out ./src/utils/*.go
test-webscrape: sql
	go test -v -coverprofile coverage.out ./src/webscrape/*.go
test-embeddings: sql
	go test -v -coverprofile coverage.out ./src/embeddings/*.go
test-customer: sql
	go test -v -coverprofile coverage.out ./src/customer/*.go
test-llm: sql
	go test -v -coverprofile coverage.out ./src/llm/*.go
test-project: sql
	go test -v -coverprofile coverage.out ./src/project/*.go
test-conv: sql
	go test -v -coverprofile coverage.out ./src/conversation/*.go

test: sql
	go test -v ./...

coverage:
	go tool cover -html=coverage.out

test-coverage: sql
	go test -v -coverprofile coverage.out ./...
	go tool cover -html=coverage.out