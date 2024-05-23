import typer
from typing_extensions import Annotated
import json

from src.lib import datastore, project, sites, vectorstore, customer

app = typer.Typer()


@app.command(help="Gets a customer based on the name in config.json")
def get_customer():
    c = customer.Customer.get()
    if c is None:
        print("there was an error getting the customer")
        exit(1)

    print(json.dumps(c.json(), indent=4))


@app.command(
    help="Ingest all documents with the same folder structure based on the passed folder. If this command is run again on the same folder with changes, these changes will be accounted for and the remote store will be kept in sync"
)
def ingest(
    folder: Annotated[
        str,
        typer.Argument(help="Folder to construct the document store from."),
    ] = "./docstore",
):
    datastore.ingest(path=folder)


@app.command(help="Ingest a website and its pages based on the rules provided.")
def ingest_web(
    domain: Annotated[
        str,
        typer.Argument(help="Domain to parse. Must include the protocol (http/https)."),
    ],
    insert: Annotated[
        bool,
        typer.Option(
            help="Insert the website into the database or just return the pages based on the rules."
        ),
    ] = False,
    whitelist: Annotated[
        str,
        typer.Option(
            help="A comma-separated list of regexp that page routes MUST match."
        ),
    ] = "",
    blacklist: Annotated[
        str,
        typer.Option(
            help="A comma-separated list of regexp that page routes can NOT match."
        ),
    ] = "",
):
    # parse the whitelist / blacklist
    wlist = []
    blist = []
    if whitelist != "":
        wlist = whitelist.split(",")
    if blacklist != "":
        blist = blacklist.split(",")

    response = sites.ingest_site(
        domain=domain,
        whitelist=wlist,
        blacklist=blist,
        insert=insert,
    )
    if response is None:
        print("There was an error running this command")
        exit(1)

    print("---\nSite:")
    print(json.dumps(response["site"], indent=4))

    print("---\nPages:")
    for i in response["pages"]:
        print(i["url"])


@app.command(help="Vectorize all ingested documents.")
def vec_dstore():
    if not datastore.vectorize():
        print("There was an issue vectorizing the datastore")
        exit(1)


@app.command(help="Vectorize all ingested websites")
def vec_web():
    if not sites.vectorize():
        print("There was an issue vectorizing the websites")
        exit(1)


@app.command(
    help="Queries the vectorstore, returning all documents and website page matches"
)
def query(
    query: str,
    k: Annotated[
        int,
        typer.Option(help="How many of each type to return from the vector query"),
    ] = 3,
    include: Annotated[
        bool,
        typer.Option(
            help="Whether to include the entire content from the response or not"
        ),
    ] = False,
):
    response = vectorstore.query(
        query=query,
        k=k,
        include=include,
    )
    if response is None:
        print("there was an issue sending the request")
        exit(1)

    print(json.dumps(response, indent=4))


@app.command(
    help="Creates a project with your customer and writes the id to the config file"
)
def create_project(
    title: str,
    topic: str,
):
    c = customer.Customer.get()
    response = c.create_project(title=title, topic=topic)
    if response is None:
        print("issue creating the project")
        exit(1)

    print(json.dumps(response, indent=4))


@app.command(
    help="Gets the project with the `currentProjectId` variable in the config.json"
)
def get_project():
    response = project.Project.get()
    if response is None:
        print("issue getting the project")
        exit(1)

    print(json.dumps(response.json(), indent=4))


@app.command()
def generate_ideas(
    k: Annotated[
        int,
        typer.Option(help="How many project ideas to generate"),
    ] = 3,
):
    p = project.Project.get()
    if p is None:
        print("issue getting the project")
        exit(1)

    p.generate_ideas(k=k)


if __name__ == "__main__":
    app()
