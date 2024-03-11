sql:
	awk 'FNR==1{print ""}1' ./schema/queries/* > ./schema/queries.sql
	sqlc generate
	rm ./schema/queries.sql

dev:
	go run ./src/*.go