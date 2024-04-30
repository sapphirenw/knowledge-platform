import hashlib
import binascii
import os


def gen_sig(input: str) -> str:
    if isinstance(input, str):
        input = input.encode("utf-8")
    hash_obj = hashlib.sha256(input)
    hash_digest = hash_obj.digest()
    return binascii.hexlify(hash_digest).decode("utf-8")


def list_files_folders(path):
    folders = []
    files = []

    # List all entries in the directory given by "path"
    for entry in os.listdir(path):
        # Join the path to get full file path
        full_path = os.path.join(path, entry)
        # Check if it's a directory or file
        if os.path.isdir(full_path):
            folders.append(full_path)
        elif os.path.isfile(full_path):
            files.append(full_path)

    return folders, files
