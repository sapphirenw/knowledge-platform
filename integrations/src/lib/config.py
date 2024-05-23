from dataclasses import dataclass
import json

from src.lib import project


@dataclass
class Config:
    name: str
    current_project_id: str

    @classmethod
    def read(cls):
        try:
            with open("./config.json", "r") as f:
                raw = f.read()
                data = json.loads(raw)

                return cls(
                    name=data["name"],
                    current_project_id=data["currentProjectId"],
                )
        except Exception as e:
            print(f"error reading the config: {e}")
            return None

    def write(self):
        with open("./config.json", "w") as f:
            f.write(json.dumps(self.json(), indent=4))
            f.close()

    def json(self):
        return {
            "name": self.name,
            "currentProjectId": self.current_project_id,
        }
