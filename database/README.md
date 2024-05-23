# Database

A local environment for creating the database running in a docker container.

## Launching the Database

You can launch the database with `make run`. If there are changes to the schema, you can run `make rebuild`. Then, you can connect to the database with `make connect`. If there are errors, you can get logs with `make logs`.