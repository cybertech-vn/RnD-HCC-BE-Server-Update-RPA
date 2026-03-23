package software

import (
	"github.com/SaraNguyen999/setup-server/internal/server"
	"github.com/gofiber/fiber/v2"
)

func init() {
	// Initialize any necessary resources or configurations for the distribute package
	server.GetServer().Register(func(a *fiber.App) {
		a.Post("/api/v1/checkver", server.LicenseMiddleware(), CheckUpdate)
	})
}
