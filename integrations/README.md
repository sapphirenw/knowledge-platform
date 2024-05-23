# Integrations

A suite of integration tests written in Python meant to test the api.

## Environment

First, create a new environment with the following commands:

```shell
$ python -m venv ./venv
$ source venv/bin/activate
$ pip install -r requirements.txt
```

## Config

There exists a `config.json` file which stores some metadata about your instance. Set the `name` field as your name, and leave the rest of the fields alone.

Now, you can see all available commands with `python main.py --help`

## Vector Store

You can load documents and folders inside of `./docstore` with `python main.py ingest` to load them into the vector store. when you make changes, the backend will automatically react and account for the changes without duplicating the contents inside the folder.

You can ingest a website running the `python main.py ingest-web --help` command.

You can then vectorize the various content with `python main.py vec-dstore` and `python main.py vec-web`.

You can query the vector store with `python main.py query`.

## Projects

You can create a new project with `python main.py create-project "My Title" "My Topic"`. This will create a project and write the value to your config file.

You can generate ideas for the project with `python main.py generate-ideas`. This will use the id in your config file. This function launches an interactive feedback loop where you can iteratively generate ideas by supplying feedback.