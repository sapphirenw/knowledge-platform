from datetime import datetime
from typing import Optional
import requests
from dataclasses import dataclass

from src.lib import consts, config


@dataclass
class Customer:
    id: int
    name: str
    datastore: str
    created_at: datetime
    updated_at: datetime

    @classmethod
    def get(cls):
        try:
            # read the config
            conf = config.Config.read()
            if conf is None:
                raise Exception("error reading the config")

            response = requests.get(
                f"{consts.HOST}/tests/customers/get?name={conf.name}",
            )
            if response.status_code != 200:
                raise Exception(
                    f"There was an issue getting the customer: {response.text}"
                )

            customer = response.json()

            return cls(
                id=customer["id"],
                name=customer["name"],
                datastore=customer["datastore"],
                created_at=datetime.fromisoformat(customer["createdAt"].rstrip("Z")),
                updated_at=datetime.fromisoformat(customer["updatedAt"].rstrip("Z")),
            )

        except Exception as e:
            print(f"error getting the customer: {e}")
            return None

    def json(self):
        return {
            "id": self.id,
            "name": self.name,
            "datastore": self.datastore,
            "createdAt": self.created_at.isoformat(),
            "updatedAt": self.updated_at.isoformat(),
        }
