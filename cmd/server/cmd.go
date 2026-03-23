package main

import (
	_ "github.com/SaraNguyen999/setup-server/internal/distribute"
	"github.com/SaraNguyen999/setup-server/internal/server"
	_ "github.com/SaraNguyen999/setup-server/internal/software"
)

func main() {
	server.StartWSN()
}
