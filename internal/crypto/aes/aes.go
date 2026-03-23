package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/crypto/hkdf"
)

type AESCryptor struct {
	Name      string
	Key       []byte
	KeySeed   string
	KeyLength int
}

func NewAES(name string, key []byte, keySeed string, keyLength int) *AESCryptor {
	if keyLength == 0 {
		keyLength = 32
	}
	return &AESCryptor{Name: name, Key: key, KeySeed: keySeed, KeyLength: keyLength}
}

func (a *AESCryptor) GetName() string {
	return a.Name
}

// GenSeed generates a random seed string consisting of 6 numbers (0-999) joined by "-".
func (a *AESCryptor) GenSeed() string {
	var seeds []string
	for range 6 {
		r := randBelow(1000)
		seeds = append(seeds, fmt.Sprintf("%03d", r))
	}
	a.KeySeed = strings.Join(seeds, "-")
	return a.KeySeed
}

// GenKey generates the AES key from the seed using HKDF-SHA256.
func (a *AESCryptor) GenKey() ([]byte, error) {
	if a.KeySeed == "" {
		a.GenSeed()
	}
	seedBytes := []byte(a.KeySeed)
	seedDigest := sha256.Sum256(seedBytes)
	h := hkdf.New(sha256.New, seedDigest[:], []byte("fixed-hkdf-salt-v1"), []byte("aes-key-derivation"))
	key := make([]byte, a.KeyLength)
	_, err := io.ReadFull(h, key)
	if err != nil {
		return nil, err // To match Python's raising behavior
	}
	a.Key = key
	return key, nil
}

// EncryptData encrypts the plaintext using AES-256-GCM.
// Returns ciphertext, nonce, tag separately.
func (a *AESCryptor) EncryptData(plaintext string, nonce []byte) ([]byte, []byte, []byte, error) {
	return a.EncryptBytes([]byte(plaintext), nonce)
}

func (a *AESCryptor) EncryptBytes(data []byte, nonce []byte) ([]byte, []byte, []byte, error) {

	if len(a.Key) == 0 {
		return nil, nil, nil, fmt.Errorf("key not generated")
	}

	if nonce == nil {
		nonce = make([]byte, 12)
		rand.Read(nonce)
	}

	block, err := aes.NewCipher(a.Key)
	if err != nil {
		return nil, nil, nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, nil, err
	}

	out := gcm.Seal(nil, nonce, data, nil)

	ciphertext := out[:len(out)-16]
	tag := out[len(out)-16:]

	return ciphertext, nonce, tag, nil
}

// DecryptData decrypts the ciphertext using AES-256-GCM.
// This is compatible with the Python version and requires the key separately.
func (a *AESCryptor) DecryptData(key, ciphertext, nonce, tag []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	ctWithTag := append(ciphertext, tag...)
	plain, err := gcm.Open(nil, nonce, ctWithTag, nil)
	if err != nil {
		return "", err
	}
	return string(plain), nil
}

// SaveKey saves the key or seed to a JSON file.
func (a *AESCryptor) SaveKey(filePath string) {
	var data map[string]string

	if a.KeySeed == "" {
		panic("No seed to save. Call GenSeed() first.")
	}
	// Mimic the Python code exactly, including the character-wise join (potential original bug, but for compatibility)
	joined := strings.Join(strings.Split(a.KeySeed, ""), ",")
	sum := sha256.Sum256([]byte(joined))
	checksum := hex.EncodeToString(sum[:])[:16]
	data = map[string]string{
		"type":     "seed",
		"name":     a.Name,
		"key_seed": a.KeySeed,
		"checksum": checksum,
	}

	jsonData, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		panic(err)
	}
	err = os.MkdirAll(filepath.Dir(filePath), 0755)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(filePath, jsonData, 0644)
	if err != nil {
		panic(err)
	}
	fmt.Printf("[%s] Saved %s to %s\n", a.Name, data["type"], filePath)
}

// LoadKey loads the key or seed from a JSON file and sets it on the instance.
func (a *AESCryptor) LoadKey(filePath string) []byte {
	buf, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	var data map[string]string
	err = json.Unmarshal(buf, &data)
	if err != nil {
		panic(err)
	}
	typ, ok := data["type"]
	if !ok {
		panic("File invalid (missing type=seed|key).")
	}
	switch typ {
	case "seed":
		seed, ok := data["key_seed"]
		if !ok {
			panic("Missing key_seed in file.")
		}
		a.KeySeed = seed
		a.GenKey()
		fmt.Printf("[%s] Loaded seed & regenerated key.\n", a.Name)
	case "key":
		kh, ok := data["key_hex"]
		if !ok {
			panic("Missing key_hex in file.")
		}
		a.Key, err = hex.DecodeString(kh)
		if err != nil {
			panic(err)
		}
		fmt.Printf("[%s] Loaded raw key from file.\n", a.Name)
	default:
		panic("File invalid (type must be seed|key).")
	}
	return a.Key
}

func randBelow(n int) int {
	bi, err := rand.Int(rand.Reader, big.NewInt(int64(n)))
	if err != nil {
		panic(err)
	}
	return int(bi.Int64())
}
