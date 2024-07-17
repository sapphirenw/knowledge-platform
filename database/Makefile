# combines all schema files together and creates the master schema document
db-schema:
	awk 'FNR==1{print ""}1' ./schema/*.sql > ./schema.sql

sql: db-schema
	awk 'FNR==1{print ""}1' ./queries/* > ./queries.sql
	awk 'FNR==1{print ""}1' ./migrations/* >> ./queries.sql
	sqlc generate
	rm ./queries.sql

mig: sql
	goose postgres "postgresql://aicontent:aicontent@localhost:5432/aicontent" up

clean:
	rm -rf ./data
	docker ps -a --filter "ancestor=database_aicontent-db-testing" -q | xargs -r docker rm -f
	@IMAGES=$$(docker images -f "dangling=true" -q); \
	if [ -n "$$IMAGES" ]; then \
		docker rmi $$IMAGES; \
	fi

run: db-schema
	docker compose up -d --remove-orphans

rebuild: db-schema
	$(MAKE) clean
	$(MAKE) run

test: db-schema
	$(MAKE) clean
	docker compose up --build
	$(MAKE) clean

logs:
	docker logs aicontent-db-testing

connect:
	psql 'postgresql://aicontent:aicontent@localhost:5432/aicontent'

# view database schema information visually
schemaspy: run
	docker run -v "$(PWD)/schemaspy/output:/output" -v "$(PWD)/schemaspy/schemaspy.properties:/schemaspy.properties" --network database_postgres-net schemaspy/schemaspy:latest
	open $(PWD)/schemaspy/output/index.html
	docker ps -a --filter "ancestor=schemaspy/schemaspy:latest" -q | xargs -r docker rm