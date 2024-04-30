import json
import requests

from src.lib import consts, customer


def ingest_site(
    domain: str,
    whitelist: list[str],
    blacklist: list[str],
    insert: bool,
):
    try:
        c = customer.Customer.get()
        if c is None:
            return False

        payload = {
            "domain": domain,
            "whitelist": whitelist,
            "blacklist": blacklist,
            "insert": insert,
        }

        response = requests.post(
            f"{consts.HOST}/customers/{c.id}/websites",
            data=json.dumps(payload),
        )
        if response.status_code != 200:
            raise Exception(f"error sending the request: {response.text}")

        data = response.json()
        return data
    except Exception as e:
        print(f"Unknown error ingesting the website: {e}")
        return None


def vectorize() -> bool:
    try:
        c = customer.Customer.get()
        if c is None:
            return False

        response = requests.put(
            f"{consts.HOST}/customers/{c.id}/vectorizeWebsites",
        )
        if response.status_code != 204:
            raise Exception(f"error sending the vectorization request: {response.text}")

        return True
    except Exception as e:
        print(f"Unknown error purging the datastore: {e}")
        return False
