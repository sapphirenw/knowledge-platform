from dataclasses import dataclass
from datetime import datetime
import json

import requests

from src.lib import config, consts


@dataclass
class Project:
    id: str
    customer_id: str
    title: str
    topic: str
    idea_generation_model_id: str
    created_at: datetime
    updated_at: datetime

    @classmethod
    def get(cls):
        try:
            c = config.Config.read()
            if c.current_project_id == "":
                raise Exception(
                    "there is no project id in the config file. first create a project with `python main.py create-project --help`"
                )

            response = requests.get(
                f"{consts.HOST}/projects/{c.current_project_id}",
            )
            if response.status_code != 200:
                raise Exception(f"There was an issue with the request: {response.text}")

            p = response.json()

            return cls(
                id=p["id"],
                customer_id=p["customerId"],
                title=p["title"],
                topic=p["topic"],
                idea_generation_model_id=p["ideaGenerationModelId"],
                created_at=datetime.fromisoformat(p["createdAt"].rstrip("Z")),
                updated_at=datetime.fromisoformat(p["updatedAt"].rstrip("Z")),
            )
        except Exception as e:
            print(f"error getting the project: {e}")
            return None

    def generate_ideas(cls, k: int):
        try:
            c = config.Config.read()
            if c.current_project_id == "":
                raise Exception(
                    "there is no project id in the config file. first create a project with `python main.py create-project --help`"
                )
            print(f"generating project ideas for project: {c.current_project_id} ...")

            # data for the idea generation
            data = {"k": k}

            # run inside a loop for as much feedback as needed
            while True:
                response = requests.post(
                    f"{consts.HOST}/projects/{c.current_project_id}/generateIdeas",
                    data=json.dumps(data),
                )
                if response.status_code != 200:
                    raise Exception(
                        f"There was an issue with the request: {response.text}"
                    )

                content = response.json()

                print("\nIdeas:")
                for i in content["ideas"]:
                    print(f"- {i['title']}")

                print("")
                feedback = input("Feedback (empty to exit): ")
                if feedback == "":
                    break

                data["feedback"] = feedback
                data["conversationId"] = content["conversationId"]

                print("Re-generating with feedback ...")

        except Exception as e:
            print(f"error generating ideas: {e}")
            return None

    def json(self):
        return {
            "id": self.id,
            "customer_id": self.customer_id,
            "title": self.title,
            "topic": self.topic,
            "ideaGenerationModelId": self.idea_generation_model_id,
            "createdAt": self.created_at.isoformat(),
            "updatedAt": self.updated_at.isoformat(),
        }
