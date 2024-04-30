include $(shell sed 's/export //g' .env > .env.make && echo .env.make)
export

sql:
	awk 'FNR==1{print ""}1' ../database/queries/* > ../database/queries/queries.sql
	sqlc generate
	rm ../database/queries/queries.sql

run:
	go run ./src/*.go

# tests for various packages
test-database:
	go test -v -coverprofile coverage.out ./src/database/*.go
test-docstore:
	go test -v -coverprofile coverage.out ./src/docstore/*.go
test-utils:
	go test -v -coverprofile coverage.out ./src/utils/*.go
test-webscrape:
	go test -v -coverprofile coverage.out ./src/webscrape/*.go
test-customer:
	go test -v -coverprofile coverage.out ./src/customer/*.go

test:
	go test -v ./...

coverage:
	go tool cover -html=coverage.out

test-coverage:
	go test -v -coverprofile coverage.out ./...
	go tool cover -html=coverage.out