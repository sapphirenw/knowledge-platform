include $(shell sed 's/export //g' .env > .env.make && echo .env.make)
export

sql:
	awk 'FNR==1{print ""}1' ./schema/queries/* > ./schema/queries.sql
	sqlc generate
	rm ./schema/queries.sql

run:
	go run ./src/*.go