package distribute

import (
	"os"

	"github.com/SaraNguyen999/setup-server/internal/server"
	"github.com/gofiber/fiber/v2"
)

const defaultUploadSigningSecret = "c8a9f1e3b7d4c2a6e5f8b0d1a4c9e7f2d6b4a1c0e8f3d5b7a9c2e6f1d4b8a0"

func uploadSigningSecret() string {
	if s := os.Getenv("SECRET"); s != "" {
		return s
	}
	return defaultUploadSigningSecret
}

func init() {
	// Initialize any necessary resources or configurations for the distribute package
	secret := uploadSigningSecret()
	server.GetServer().Register(func(a *fiber.App) {
		a.Post("/api/v1/upload", server.UploadAuthMiddleware(secret), UploadVersion)
		// Backward-compatible alias for clients still calling /upload directly.
		a.Post("/upload", server.UploadAuthMiddleware(secret), UploadVersion)
		a.Get("/api/v1/download/:app_id/:version", server.LicenseMiddleware(), DownloadVersion)
	})
}
