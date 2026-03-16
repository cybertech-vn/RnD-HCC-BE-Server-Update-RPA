import ctypes
import os

# Load thư viện (thay đường dẫn nếu cần)
if os.name == "nt":  # Windows
    lib = ctypes.CDLL("./libfp.dll")
else:
    lib = ctypes.CDLL("./libfp.so")

# Khai báo kiểu (rất quan trọng)
lib.CGenerateFingerprint.restype = ctypes.c_char_p
lib.CGetHostname.restype = ctypes.c_char_p
lib.FreeCStr.argtypes = [ctypes.c_char_p]

# === Sử dụng ===
def get_fingerprint():
     fingerprint = lib.CGenerateFingerprint()
     return fingerprint.decode("utf-8")

def get_hostname():
    hostname = lib.CGetHostname()
    return hostname.decode("utf-8")
