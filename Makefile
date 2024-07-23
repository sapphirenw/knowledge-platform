
run:
	docker compose up --build

build-api:
	docker buildx build --platform linux/amd64,linux/arm64 -t ghcr.io/sapphirenw/aicontent-api:latest --push ./api

build-db:
	$(MAKE) -C "$(shell pwd)/database" -f "$(shell pwd)/database/Makefile" sql
	docker buildx build --platform linux/amd64,linux/arm64 -t ghcr.io/sapphirenw/aicontent-db:latest --push ./database

build-webapp:
	docker buildx build --platform linux/amd64,linux/arm64 -t ghcr.io/sapphirenw/aicontent-webapp:latest --push ./webapp

build: build-api build-webapp