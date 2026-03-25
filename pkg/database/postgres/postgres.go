package postgres

import (
	"fmt"

	mydb "github.com/SaraNguyen999/setup-server/pkg/database"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PGDBM struct {
	dsn      string
	user     string
	password string
	Database string
	host     string
	port     int
	db       *gorm.DB
}

func (p *PGDBM) Init(m mydb.BaseConnect) (*gorm.DB, error) {

	// nếu có DSN thì dùng luôn
	if m.Dsn != "" {
		return p.InitByString(m.Dsn)
	}

	p.user = m.User
	p.password = m.Password
	p.Database = m.Database
	p.host = m.Host
	p.port = m.Port

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=require",
		p.host,
		p.user,
		p.password,
		p.Database,
		p.port,
	)

	return p.InitByString(dsn)
}

func (p *PGDBM) InitByString(dsn string) (*gorm.DB, error) {

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("gorm connect error: %w", err)
	}

	// cấu hình connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(2)

	p.db = db

	return p.db, nil
}

func (p *PGDBM) GetDB() *gorm.DB {
	return p.db
}

func (p *PGDBM) Close() {

	if p.db == nil {
		return
	}

	sqlDB, err := p.db.DB()
	if err == nil {
		sqlDB.Close()
	}
}
