run:
	docker compose up -d --build

rebuild:
	rm -rf ./data
	$(MAKE) run

connect:
	psql 'postgresql://aicontent:aicontent@localhost:5432/aicontent'