# combines all schema files together and creates the master schema document
db-schema:
	awk 'FNR==1{print ""}1' ./schema/*.sql > ./schema.sql

run: db-schema
	docker compose up -d --build

rebuild: db-schema
	rm -rf ./data
	docker ps -a --filter "ancestor=database_aicontent-db" -q | xargs -r docker rm -f
	$(MAKE) run

connect:
	psql 'postgresql://aicontent:aicontent@localhost:5432/aicontent'

# view database schema information visually
make schemaspy: run
	docker run -v "$(PWD)/schemaspy/output:/output" -v "$(PWD)/schemaspy/schemaspy.properties:/schemaspy.properties" --network database_postgres-net schemaspy/schemaspy:latest
	open $(PWD)/schemaspy/output/index.html
	docker ps -a --filter "ancestor=schemaspy/schemaspy:latest" -q | xargs -r docker rm