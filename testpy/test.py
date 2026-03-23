import hashlib
import hmac
from time import time
from urllib.parse import urlparse

import requests
from aes import NewAES
from fp import get_fingerprint

baseUrl = "http://localhost:7878/api/v1"
id=get_fingerprint()

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
    url = f"{baseUrl}/upload"

    app_id = "license_agent"
    secret = "c8a9f1e3b7d4c2a6e5f8b0d1a4c9e7f2d6b4a1c0e8f3d5b7a9c2e6f1d4b8a0"

    parsed = urlparse(url)
    path = parsed.path

    headers = sign_request(app_id, secret, "POST", path)

    with open('test.txt', 'rb') as f:
        files = {'file': f}
        data = {
            "version": "1.0.3"
        }

        response = requests.post(
            url,
            headers=headers,
            files=files,
            data=data
        )

    print(response.status_code)
    print(response.text)
    
def get_pass():
    res=requests.post(f"http://localhost:7776/api/v1/psk/{id}",verify=False)
    print(res.text)
    return res.json()["data"]["passkey"]


def test_download():
    url = f"{baseUrl}/download/license_agent/latest"
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
    url = f"{baseUrl}/checkver"
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
    test_upload()
    # test_download()
    # test_checkversion()