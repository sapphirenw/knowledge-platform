from dataclasses import dataclass
import json


@dataclass
class Config:
    name: str

    @classmethod
    def read(cls):
        try:
            with open("./config.json", "r") as f:
                raw = f.read()
                data = json.loads(raw)

                return cls(name=data["name"])
        except Exception as e:
            print(f"error reading the config: {e}")
            return None
