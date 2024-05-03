import requests
import mimetypes
import os
import json
import base64
from datetime import datetime, timezone, timedelta

from src.lib import customer
from . import utils, consts


def ingest(path: str):
    try:
        c = customer.Customer.get()
        if c is None:
            return False

        start = datetime.now()
        if not __process_folder(c=c, parent=0, path=path):
            raise Exception("There was a fatal issue ingesting the docstore")

        end = datetime.now()

        if not __purge(
            c=c,
            timestamp=datetime.now(timezone.utc)
            - (end - start + timedelta(seconds=10)),
        ):
            raise Exception("There was a fatal issue purging the datastore")

        return True
    except Exception as e:
        print(f"Unknown error ingesting the datastore: {e}")
        return False


def vectorize() -> bool:
    try:
        c = customer.Customer.get()
        if c is None:
            return False
        response = requests.put(
            f"{consts.HOST}/customers/{c.id}/vectorizeDocuments",
        )
        if response.status_code != 204:
            raise Exception(f"error sending the vectorization request: {response.text}")

        return True
    except Exception as e:
        print(f"Unknown error purging the datastore: {e}")
        return False


def __process_folder(c: customer.Customer, parent: int, path: str):
    try:
        # create the folder
        folder_name = path.split("/")[-1]
        folder_id = None

        if folder_name != "docstore":
            # upload to the server
            payload = {
                "owner": parent,
                "name": folder_name,
            }
            response = requests.post(
                f"{consts.HOST}/customers/{c.id}/folders",
                data=json.dumps(payload),
            )

            if response.status_code == 409:
                print("This folder already exists")
            elif response.status_code != 200:
                raise Exception(f"error creating the folder: {response.text}")

            data = response.json()
            print(f"Created folder: {data}")
            folder_id = data["id"]

        folders, files = utils.list_files_folders(path=path)

        # process all files in the directory
        for file in files:
            if not __upload_file(c=c, parent=folder_id, path=file):
                raise Exception(f"Error uploading: {file}")

        # process all subfolders
        for folder in folders:
            if not __process_folder(c=c, parent=folder_id, path=folder):
                raise Exception(f"error processing subfolder: {folder}")

        return True
    except Exception as e:
        print(f"There was an error: {e}")
        return None


def __upload_file(c: customer.Customer, parent: int, path: str) -> bool:
    try:
        # open the file
        with open(path, "rb") as f:
            contents = f.read()

        filename = os.path.basename(path)
        mime, _ = mimetypes.guess_type(path)
        if mime is None:
            mime = "content/text"

        input = {
            "parentId": parent,
            "filename": filename,
            "mime": mime,
            "signature": utils.gen_sig(contents),
            "size": len(contents),
        }

        print(f"Sending request ... {input}")

        response = requests.post(
            f"{consts.HOST}/customers/{c.id}/generatePresignedUrl",
            data=json.dumps(input),
        )

        if response.status_code == 409:
            print("This file already exists")
            return True
        elif response.status_code != 200:
            raise Exception(f"failed to generate pre-signed url: {response.text}")

        print("Successfully generated presigned url")
        body = response.json()
        upload_url = base64.b64decode(body["uploadUrl"]).decode("utf-8")
        documentId = body["documentId"]

        print(f"Upload URL: {upload_url}")
        print(f"DocumentId: {documentId}")

        print("Uploading the file ...")

        s3_headers = {
            "Content-Type": mime,
            "Content-Disposition": f'attachment; filename="{filename}"',
        }

        s3_response = requests.put(
            upload_url,
            data=contents,
            headers=s3_headers,
        )

        if s3_response.status_code != 200:
            raise Exception(
                f"There was an issue sending the s3 request: {s3_response.text}"
            )

        # report that the file was successfully uploaded
        notif_response = requests.put(
            f"{consts.HOST}/customers/{c.id}/documents/{documentId}/validate",
        )

        if notif_response.status_code != 204:
            raise Exception(
                f"failed to notifiy of successful upload: {notif_response.text}"
            )

        print("Successfully uploaded file")

        return True
    except IOError as e:
        print(f"Error opening or reading the file: {e}")
        return False
    except Exception as e:
        print(f"Unknown error processing the file: {e}")
        return False


def __purge(c: customer.Customer, timestamp: datetime) -> bool:
    try:
        payload = {"timestamp": timestamp.strftime("%Y-%m-%d %H:%M:%S")}
        response = requests.post(
            f"{consts.HOST}/customers/{c.id}/datastore/purge",
            data=json.dumps(payload),
        )
        if response.status_code != 204:
            raise Exception(f"There was an issue sending the request: {response.text}")

        print("Successfully purged datastore")

        return True
    except Exception as e:
        print(f"Unknown error purging the datastore: {e}")
        return False
