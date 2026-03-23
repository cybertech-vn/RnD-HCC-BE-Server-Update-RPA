package postgres

import (
	"github.com/SaraNguyen999/setup-server/internal/config"
	pog "github.com/SaraNguyen999/setup-server/pkg/database/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func init() {
	DB = connectDB().GetDB()
}

func connectDB() *pog.PGDBM {
	var db = pog.PGDBM{}
	_, err := db.Init(config.CONFIG.DBConfig.BaseConnect)
	if err != nil {
		panic(err)
	}
	return &db
}

func ConnectDB() *pog.PGDBM {
	return connectDB()
}
