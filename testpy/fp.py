import ctypes
import os

# Load thư viện theo đường dẫn tuyệt đối tại thư mục hiện tại của file.
_HERE = os.path.dirname(os.path.abspath(__file__))
_CANDIDATES = ["libfp.dll", "libfp.so"] if os.name == "nt" else ["libfp.so", "libfp.dll"]

lib = None
load_err = None
for name in _CANDIDATES:
    lib_path = os.path.join(_HERE, name)
    if not os.path.exists(lib_path):
        continue
    try:
        lib = ctypes.CDLL(lib_path)
        break
    except OSError as err:
        load_err = err

if lib is None:
    searched = ", ".join(os.path.join(_HERE, c) for c in _CANDIDATES)
    if load_err is not None:
        raise OSError(f"Cannot load fingerprint library. Tried: {searched}. Last error: {load_err}")
    raise FileNotFoundError(f"Fingerprint library not found. Expected one of: {searched}")

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
