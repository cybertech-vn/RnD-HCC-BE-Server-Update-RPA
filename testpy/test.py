import hashlib
import hmac
import os
from time import time

import requests
from aes import NewAES

try:
    from fp import get_fingerprint
except Exception:
    get_fingerprint = None

BASE_URL = os.getenv("BASE_URL", "http://localhost:7878")
API_PREFIX = os.getenv("API_PREFIX", "/api/v1")


def load_env_file(path):
    if not os.path.exists(path):
        return
    with open(path, "r", encoding="utf-8") as f:
        for raw in f:
            line = raw.strip()
            if not line or line.startswith("#") or "=" not in line:
                continue
            key, value = line.split("=", 1)
            os.environ[key.strip()] = value.strip()


_HERE = os.path.dirname(os.path.abspath(__file__))
load_env_file(os.path.join(_HERE, ".env"))

BASE_URL = os.getenv("BASE_URL", BASE_URL)
API_PREFIX = os.getenv("API_PREFIX", API_PREFIX)


def get_env(name, default=""):
    v = os.getenv(name, "").strip()
    return v if v else default


def resolve_client_id():
    # Prefer explicit env var for local tests.
    env_id = get_env("CLIENT_ID")
    if env_id:
        return env_id
    if get_fingerprint is None:
        return "license_agent"
    return get_fingerprint()


id = resolve_client_id()

def sign_request(app_id, secret, method, path):
    timestamp = str(int(time()))

    data = app_id + timestamp + method + path

    signature = hmac.new(
        secret.encode(),
        data.encode(),
        hashlib.sha256
    ).hexdigest()

    return {
        "X-App-Id": app_id,
        "X-Timestamp": timestamp,
        "X-Signature": signature
    }


def test_upload():
    upload_path = f"{API_PREFIX}/upload"
    url = f"{BASE_URL}{upload_path}"

    app_id = get_env("APP_ID", "license_agent")
    secret = get_env("SECRET")
    version = get_env("VERSION", "1.0.4")

    headers = sign_request(app_id, secret, "POST", upload_path)

    test_file = os.path.join(os.path.dirname(__file__), "test.txt")
    with open(test_file, 'rb') as f:
        files = {'file': f}
        data = {
            "version": version
        }

        response = requests.post(
            url,
            headers=headers,
            files=files,
            data=data,
            timeout=60,
        )

    print(response.status_code)
    print(response.text)
    
def get_pass():
    res=requests.post(f"http://localhost:7776/api/v1/psk/{id}",verify=False)
    print(res.text)
    return res.json()["data"]["passkey"]


def test_download():
    url = f"{BASE_URL}{API_PREFIX}/download/license_agent/latest"
    passkey = get_pass()

    headers = {
        "X-Client-ID": id,
        "X-PassKey": passkey
    }

    r = requests.get(url, headers=headers)
    r.raise_for_status()

    data = r.content

    nonce = data[:12]
    tag = data[-16:]
    ciphertext = data[12:-16]

    crypt = NewAES("securesrc", key_seed=passkey)
    key = crypt.gen_key()

    plaintext = crypt.decrypt_data(key, ciphertext, nonce, tag)

    # tính SHA256 checksum
    sha256 = hashlib.sha256()
    sha256.update(plaintext)
    checksum = sha256.hexdigest()

    # ghi file
    with open("agent.exe", "wb") as f:
        f.write(plaintext)

    print("Download + decrypt OK:", "agent.exe")
    print("SHA256:", checksum)
    print("Download test passed!")    
    
def test_checkversion():
    url = f"{BASE_URL}{API_PREFIX}/checkver"
    passkey = get_pass()
    header = {
        "X-Client-ID": id,
        "X-PassKey": passkey
    }
    
    
    data = {
        "app_id": "license_agent",
        "version": "latest"
    }
    response = requests.post(url, json=data, headers=header)
    
    print("data return:",response.status_code, response.text)
    print("Check version test passed!")


if __name__ == "__main__":
    # test_upload()
    test_download()
    # test_checkversion()