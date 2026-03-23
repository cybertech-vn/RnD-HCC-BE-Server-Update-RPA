package distribute

import (
	"github.com/SaraNguyen999/setup-server/internal/server"
	"github.com/gofiber/fiber/v2"
)

func init() {
	// Initialize any necessary resources or configurations for the distribute package
	server.GetServer().Register(func(a *fiber.App) {
		a.Post("/api/v1/upload", server.UploadAuthMiddleware("c8a9f1e3b7d4c2a6e5f8b0d1a4c9e7f2d6b4a1c0e8f3d5b7a9c2e6f1d4b8a0"), UploadVersion)
		a.Get("/api/v1/download/:app_id/:version", server.LicenseMiddleware(), DownloadVersion)
	})
}
