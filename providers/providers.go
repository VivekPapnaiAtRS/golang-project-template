package providers

import "github.com/jmoiron/sqlx"

type PSQLProvider interface {
	DB() *sqlx.DB
}

type DBProvider interface {
	Ping() error
	PSQLProvider
}

type WebSocketHubProvider interface {
	Run()
	Get() interface{}
	Stop()
}
