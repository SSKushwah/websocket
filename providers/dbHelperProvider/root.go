package dbHelperProvider

import (
	"websocket/providers"

	"github.com/jmoiron/sqlx"
)

type DBHelper struct {
	DB *sqlx.DB
}

func NewDBHelperProvider(db *sqlx.DB) providers.DBHelperProvider {
	return &DBHelper{
		DB: db,
	}
}
