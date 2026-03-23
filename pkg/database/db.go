package database

type returnData any

type BaseConnect struct {
	Dsn      string `yaml:"dsn,omitempty"`
	User     string
	Password string
	Database string
	Host     string
	Port     int
}

type Database interface {
	Init(BaseConnect) error
	Connect() (returnData, error)
	GetDB() returnData
	Close()
}

type Schema interface {
	Create(name string) error
	Remove() error
	Get() (returnData, error)
}
