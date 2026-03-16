import json
import os
import secrets
import hashlib
from cryptography.hazmat.primitives.kdf.hkdf import HKDF
from cryptography.hazmat.primitives import hashes
from cryptography.hazmat.primitives.ciphers import Cipher, algorithms, modes
from cryptography.hazmat.backends import default_backend


class AES_Cryptor:
    def __init__(self, name, key=None, key_seed=None, key_length=32):
        self.name = name
        self.key = key
        self.key_seed = key_seed
        self.key_length = key_length

    def __name__(self):
        return self.name

    # ------------------------------------------------------------
    # Tạo seed số ngẫu nhiên (6 số)
    # ------------------------------------------------------------
    def gen_seed(self):
        """
        Sinh 6 số ngẫu nhiên trong [0,1000), dùng để tạo key deterministically.
        """
        key_seed = [secrets.randbelow(1000) for _ in range(6)]
        self.key_seed = "-".join(str(x).zfill(3) for x in key_seed)
        return self.key_seed

    # ------------------------------------------------------------
    # Tạo key AES-256 từ seed số, dùng HKDF-SHA256 (deterministic)
    # ------------------------------------------------------------
    def gen_key(self):
        """
        Generate secret key từ list seed_numbers (deterministic, dùng HKDF-SHA256).
        """
        if not self.key_seed:
            self.gen_seed()

        # Ghép seed thành chuỗi ví dụ: "123,456,789,321,654,987"
        seed_bytes = self.key_seed.encode("utf-8")

        # Dùng SHA256 để nén lại làm input cho HKDF (optional nhưng tốt)
        seed_digest = hashlib.sha256(seed_bytes).digest()

        # Dùng HKDF để derive key (deterministic, bảo mật)
        hkdf = HKDF(
            algorithm=hashes.SHA256(),
            length=self.key_length,
            salt=b"fixed-hkdf-salt-v1",  # có thể đổi tùy app version
            info=b"aes-key-derivation",
        )

        self.key = hkdf.derive(seed_digest)
        return self.key

    # ------------------------------------------------------------
    # AES-GCM Encrypt / Decrypt
    # ------------------------------------------------------------
    def encrypt_data(self, plaintext, nonce=None):
        """
        Mã hóa dữ liệu với AES-256-GCM.
        - Trả về: (ciphertext, nonce, tag)
        """
        if not self.key:
            raise ValueError("Key chưa được tạo. Hãy gọi gen_key() trước.")

        if not nonce:
            nonce = os.urandom(12)  # Nonce unique
        cipher = Cipher(algorithms.AES(self.key), modes.GCM(nonce), backend=default_backend())
        encryptor = cipher.encryptor()
        ciphertext = encryptor.update(plaintext.encode("utf-8")) + encryptor.finalize()
        tag = encryptor.tag
        return ciphertext, nonce, tag

    def decrypt_data(self, key, ciphertext, nonce, tag):
        """
        Giải mã dữ liệu với AES-256-GCM.
        """
        cipher = Cipher(algorithms.AES(key), modes.GCM(nonce, tag), backend=default_backend())
        decryptor = cipher.decryptor()
        plaintext = decryptor.update(ciphertext) + decryptor.finalize()
        return plaintext

# ------------------------------------------------------------
    # SAVE / LOAD KEY
    # ------------------------------------------------------------
    def save_key(self, file_path, save_seed=True):
        """
        Lưu key xuống file JSON.
        - Nếu save_seed=True: chỉ lưu seed (bảo mật hơn, có thể tái sinh key)
        - Nếu save_seed=False: lưu luôn full key bytes (ít bảo mật hơn)
        """
        if save_seed:
            if not self.key_seed:
                raise ValueError("Chưa có seed để lưu. Gọi gen_seed() trước.")
            data = {
                "type": "seed",
                "name": self.name,
                "key_seed": self.key_seed,
                "checksum": hashlib.sha256(",".join(map(str, self.key_seed)).encode()).hexdigest()[:16],
            }
        else:
            if not self.key:
                raise ValueError("Chưa có key để lưu. Gọi gen_key() trước.")
            data = {
                "type": "key",
                "name": self.name,
                "key_hex": self.key.hex(),
                "checksum": hashlib.sha256(self.key).hexdigest()[:16],
            }

        # os.makedirs(os.path.dirname(file_path), exist_ok=True)
        with open(file_path, "w", encoding="utf-8") as f:
            json.dump(data, f, indent=4)
        print(f"[{self.name}] Saved {data['type']} to {file_path}")

    def load_key(self, file_path):
        """
        Load seed hoặc key từ file JSON.
        - Nếu type='seed' → tái tạo key bằng HKDF.
        - Nếu type='key' → dùng key trực tiếp.
        """
        if not os.path.exists(file_path):
            raise FileNotFoundError(f"Không tìm thấy file: {file_path}")

        with open(file_path, "r", encoding="utf-8") as f:
            data = json.load(f)

        if data["type"] == "seed":
            self.key_seed = data["key_seed"]
            self.gen_key()
            print(f"[{self.name}] Loaded seed & regenerated key.")
        elif data["type"] == "key":
            self.key = bytes.fromhex(data["key_hex"])
            print(f"[{self.name}] Loaded raw key from file.")
        else:
            raise ValueError("File không hợp lệ (thiếu type=seed|key).")
        return self.key
    
def NewAES(name, key=None, key_seed=None, key_length=32):
    return AES_Cryptor(name, key, key_seed, key_length)