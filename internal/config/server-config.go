package config

type ServerConfig struct {
	LogPath        string   `yaml:"logPath"`        // Đường dẫn file cấu hình Log
	ServerAddress  string   `yaml:"serverAddress"`  // Ví dụ ":8080"
	AllowOrigins   []string `yaml:"allowOrigins"`   // Danh sách origin cho phép (CORS)
	AllowHeaders   []string `yaml:"allowHeaders"`   // Các header được chấp nhận
}


