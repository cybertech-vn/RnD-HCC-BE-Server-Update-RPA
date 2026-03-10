package server

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	logs "github.com/SaraNguyen999/setup-server/pkg/log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/sirupsen/logrus"
)

type FiberServer struct {
	Name   string
	logger *logrus.Logger
	app    *fiber.App
}

func (f FiberServer) GetApp() *fiber.App {
	return f.app
}

func (f FiberServer) Register(registerFunc func(*fiber.App)) {
	registerFunc(f.app)
}

func (f *FiberServer) SetupLogger(path string) {
	// Tạo thư mục log nếu chưa có
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		f.logger.Error("cannot create log directory: " + err.Error())
		return
	}

	// Nếu chưa có logger thì khởi tạo
	if f.logger == nil {
		f.logger = logs.SetupLogger(
			logs.WithFileOutput(path),
			logs.WithTextFormat(), // hoặc JSON
			logs.WithLogLevel(logrus.InfoLevel),
		)
	}

	// Middleware log request
	f.app.Use(func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()

		// Nếu có error thì ghi log lỗi riêng
		if err != nil {
			f.logger.WithFields(logrus.Fields{
				"method": c.Method(),
				"path":   c.Path(),
				"ip":     c.IP(),
				"err":    err.Error(),
			}).Error("Request error")
			return err
		}

		// Ghi log request thành công
		f.logger.WithFields(logrus.Fields{
			"method":      c.Method(),
			"path":        c.Path(),
			"ip":          c.IP(),
			"status":      c.Response().StatusCode(),
			"duration_ms": time.Since(start).Milliseconds(),
		}).Info("request")

		return nil
	})
}

func (f FiberServer) SetupCors(c cors.Config) {
	f.app.Use(cors.New(c))
}

func (f *FiberServer) Start(addr string) error {
	if err := f.app.Listen(addr); err != nil {
		f.logger.Error(err)
		return err
		// f.logger.WithError(err).Error("Server stopped")
	} else {
		f.logger.Info("Server start")
	}

	return nil
}

func (f FiberServer) StartWSN(p string) {
	// Tạo kênh để nhận tín hiệu hệ thống
	sigChannel := make(chan os.Signal, 1)
	signal.Notify(sigChannel, syscall.SIGINT, syscall.SIGTERM)

	// Chạy server trong goroutine
	go func() {
		if err := f.app.Listen(p); err != nil {
			f.logger.WithError(err).Fatal("Failed to start server")
		}
	}()

	// Đợi tín hiệu và thực hiện shutdown
	sig := <-sigChannel
	f.logger.Infof("Received signal: %v", sig)
	f.Shutdown()
}

func (f *FiberServer) Shutdown() error {
	if err := f.app.Shutdown(); err != nil {
		f.logger.WithError(err).Error("Failed to shut down server gracefully")
		return fmt.Errorf("failed to shut down server gracefully")
	} else {
		f.logger.Info("Server shut down gracefully")
	}
	return nil
}

func (f *FiberServer) FallBack(function func(c *fiber.Ctx) error) {
	f.app.Use(function)
}

func NewServer(name string) *FiberServer {
	return &FiberServer{
		Name: name,
		app:  fiber.New(fiber.Config{}),
	}
}
