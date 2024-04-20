include $(shell sed 's/export //g' .env > .env.make && echo .env.make)
export

sql:
	awk 'FNR==1{print ""}1' ./schema/queries/* > ./schema/queries.sql
	sqlc generate
	rm ./schema/queries.sql

run:
	go run ./src/*.go

# tests for various packages
test-database:
	go test -v ./src/database/*.go
test-docstore:
	go test -v ./src/docstore/*.go
test-document:
	go test -v ./src/document/*.go
test-utils:
	go test -v ./src/utils/*.go
test-webscrape:
	go test -v ./src/webscrape/*.go
test-customer:
	go test -v ./src/webscrape/*.go

# run all tests
test-small: test-database test-docstore test-document test-utils test-webscrape
	# run a subset of customer tests
	go test -v ./src/customer/*.go -run TestCustomerFolderStructure
	go test -v ./src/customer/*.go -run TestCustomerUploadDocument
	go test -v ./src/customer/*.go -run TestCustomerPurgeDatastore

test:
	go test -v ./...