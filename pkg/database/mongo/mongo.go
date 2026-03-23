package mongo

import (
	"context"
	"fmt"

	mydb "github.com/SaraNguyen999/setup-server/pkg/database"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MGDBM struct {
	user     string
	password string
	Database string
	host     string
	port     int
	db       *mongo.Client
}

func (p *MGDBM) Init(m mydb.BaseConnect) error {
	p.user = m.User
	p.password = m.Password
	p.Database = m.Database
	p.host = m.Host
	p.port = m.Port

	db, err := p.Connect()
	if err != nil {
		return err
	}
	p.db = db

	return nil
}

func (p *MGDBM) Connect() (*mongo.Client, error) {
	uri := fmt.Sprintf("mongodb://%s:%d",
		p.host,
		p.port,
	)

	var clientOptions *options.ClientOptions
	if p.user != "" {
		credential := options.Credential{
			Username: p.user,
			Password: p.password,
		}
		clientOptions = options.Client().ApplyURI(uri).SetAuth(credential)
	} else {
		clientOptions = options.Client().ApplyURI(uri)
	}

	// Kết nối đến MongoDB
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	}

	// Ping để kiểm tra kết nối
	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (p *MGDBM) GetDB() *mongo.Database {
	return p.db.Database(p.Database)
}

func (p *MGDBM) Close() {
	p.db.Disconnect(context.Background())
}
