import json
import requests

from src.lib import consts, customer


def query(query: str, k: int = 3, include: bool = False):
    try:
        c = customer.Customer.get()
        if c is None:
            return False

        payload = {"query": query, "k": k, "includeContent": include}

        response = requests.put(
            f"{consts.HOST}/customers/{c.id}/vectorstore/query",
            data=json.dumps(payload),
        )
        if response.status_code != 200:
            raise Exception(f"error sending the request: {response.text}")

        data = response.json()
        return data
    except Exception as e:
        print(f"Unknown error purging the datastore: {e}")
        return None
