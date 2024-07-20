include $(shell sed 's/export //g' .env > .env.make && echo .env.make)
export

mig-dev:
	goose postgres "$$POSTGRES_URL_DEV" up

mig-test:
	goose postgres $$POSTGRES_URL_TEST up

mig-prod:
	goose postgres "$$POSTGRES_URL_PROD" up