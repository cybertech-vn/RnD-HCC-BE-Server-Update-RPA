package server

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"

	"github.com/SaraNguyen999/setup-server/internal/config"
	"github.com/SaraNguyen999/setup-server/pkg/requests"
	"github.com/SaraNguyen999/setup-server/pkg/server"
	"github.com/gofiber/fiber/v2"
)

var SERVER *requests.Client

func init() {
	SERVER = requests.NewClient(config.CONFIG.CloudServerConfig.URL, nil, true)
}

func MiddlewareExample() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Next()
	}
}

func LicenseMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {

		clientID := c.Get("X-Client-ID")
		passKey := c.Get("X-PassKey")

		if clientID == "" || passKey == "" {
			return server.ResponseReturn(c, true, "missing license headers", nil, 403)
		}
		headers := map[string]string{
			"X-Client-ID": clientID,
			"X-PassKey":   passKey,
		}

		resp, err := SERVER.Post("/api/v1/license/license-check", nil, headers)
		if err != nil {
			fmt.Printf("Error contacting license service: %v\n", err)
			return server.ResponseReturn(c, true, "license service unreachable", nil, 502)
		}

		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			return server.ResponseReturn(c, true, "invalid license", nil, 403)
		}
		// lưu vào context
		c.Locals("client_id", clientID)
		c.Locals("passkey", passKey)

		return c.Next()
	}
}

func UploadAuthMiddleware(secret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		appID := c.Get("X-App-Id")
		timestamp := c.Get("X-Timestamp")
		signature := c.Get("X-Signature")

		if appID == "" || timestamp == "" || signature == "" {
			return fiber.ErrUnauthorized
		}

		// ===== 1. CHECK TIMESTAMP (chống replay) =====
		ts, err := strconv.ParseInt(timestamp, 10, 64)
		if err != nil {
			return fiber.ErrUnauthorized
		}

		// chỉ cho phép lệch 5 phút
		if time.Now().Unix()-ts > 300 {
			return fiber.ErrUnauthorized
		}

		// ===== 2. TẠO SIGNATURE SERVER =====
		data := fmt.Sprintf("%s%s%s%s",
			appID,
			timestamp,
			c.Method(),
			c.Path(),
		)

		h := hmac.New(sha256.New, []byte(secret))
		h.Write([]byte(data))
		expected := hex.EncodeToString(h.Sum(nil))

		// ===== 3. SO SÁNH =====
		if !hmac.Equal([]byte(expected), []byte(signature)) {
			return fiber.ErrUnauthorized
		}

		// ===== 4. inject context =====
		c.Locals("app_id", appID)

		return c.Next()
	}
}
