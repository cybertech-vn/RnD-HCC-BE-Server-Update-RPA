package distribute

import (
	"fmt"
	"io"

	"github.com/SaraNguyen999/setup-server/internal/crypto/aes"
	"github.com/SaraNguyen999/setup-server/internal/db/minio"
	"github.com/SaraNguyen999/setup-server/internal/software"
	"github.com/SaraNguyen999/setup-server/pkg/server"
	"github.com/gofiber/fiber/v2"
)

func DownloadVersion(c *fiber.Ctx) error {

	appID := c.Params("app_id")
	version := c.Params("version")

	v, err := software.GetVersion(appID, version)
	if err != nil {
		return server.ResponseReturn(c, true, "version not found", nil, 404)
	}

	objectName := fmt.Sprintf("%s/%s", appID, v.FileName)

	// lấy file từ MinIO
	obj, err := minio.MinioDB.GetObject("", objectName)
	if err != nil {
		return err
	}
	defer obj.Close()

	// đọc toàn bộ file
	data, err := io.ReadAll(obj)
	if err != nil {
		return err
	}

	// tạo AES từ passkey
	crypt := aes.NewAES("securesrc", nil, c.Locals("passkey").(string), 32)

	_, err = crypt.GenKey()
	if err != nil {
		return err
	}

	// encrypt
	ciphertext, nonce, tag, err := crypt.EncryptData(string(data), nil)
	if err != nil {
		return err
	}

	// ghép nonce + ciphertext + tag
	payload := append(nonce, ciphertext...)
	payload = append(payload, tag...)

	// set header download
	c.Set("Content-Disposition", "attachment; filename="+v.FileName+".enc")
	c.Set("Content-Type", "application/octet-stream")

	return c.Send(payload)
}
