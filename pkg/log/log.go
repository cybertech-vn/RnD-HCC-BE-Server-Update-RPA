package logs

import (
	"os"

	"github.com/sirupsen/logrus"
)

// LoggerOption định nghĩa kiểu function để tùy chỉnh logger
type LoggerOption func(*logrus.Logger)

// setupLogger với variadic options pattern
func SetupLogger(options ...LoggerOption) *logrus.Logger {
	// Tạo một instance mới của Logger
	logger := logrus.New()

	// Áp dụng từng option cho logger
	for _, option := range options {
		option(logger)
	}

	return logger
}

// Option để log ra file thay vì console
func WithFileOutput(filePath string) LoggerOption {
	return func(logger *logrus.Logger) {
		file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			logger.Fatal("Failed to log to file, using default stderr")
		}
		logger.SetOutput(file)
	}
}

// Option để set mức độ log (Info, Debug, Warn,...)
func WithLogLevel(level logrus.Level) LoggerOption {
	return func(logger *logrus.Logger) {
		logger.SetLevel(level)
	}
}

// Option để sử dụng định dạng JSON thay vì text
func WithJSONFormat() LoggerOption {
	return func(logger *logrus.Logger) {
		logger.SetFormatter(&logrus.JSONFormatter{})
	}
}

// Option để sử dụng định dạng Text với timestamp đầy đủ
func WithTextFormat() LoggerOption {
	return func(logger *logrus.Logger) {
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
	}
}
