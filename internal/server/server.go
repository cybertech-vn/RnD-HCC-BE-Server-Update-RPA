package server

import (
	"strings"

	s "github.com/SaraNguyen999/setup-server/pkg/server"

	"github.com/SaraNguyen999/setup-server/internal/config"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

var (
	f     *s.FiberServer
	sconf config.ServerConfig
)

func init() {
	f = s.NewServer("Setup Server", fiber.Config{
		BodyLimit:         1000 * 1024 * 1024,
		StreamRequestBody: true,
	})
	sconf = config.CONFIG.ServerConfig
	f.SetupLogger(sconf.LogPath)

	f.SetupCors(cors.Config{
		AllowOrigins: strings.Join(sconf.AllowOrigins, ","),
		AllowHeaders: strings.Join(sconf.AllowHeaders, ","),
	})

	f.Register(func(a *fiber.App) {
		a.Options("/*", func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusNoContent)
		})
	})

	f.Register(func(a *fiber.App) {
		a.Get("/heath", func(c *fiber.Ctx) error {
			return s.ResponseReturn(c, true, "ok", nil, 200)
		})
	})

}

func Fallback() {
	f.FallBack(func(c *fiber.Ctx) error {
		return s.ResponseReturn(c, false, "Not Found this api", nil, 404)
	})
}

func GetServer() *s.FiberServer {
	return f
}

func Start() error {
	return f.Start(sconf.ServerAddress)
}

func StartWSN() error {
	f.StartWSN(sconf.ServerAddress)
	return nil
}
