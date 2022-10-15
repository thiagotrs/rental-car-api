package database

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/thiagotrs/rentalcar-ddd/internal/pkg/config"
)

func GetDBConn(dbConf config.DBConfig) *sqlx.DB {
	connStr := fmt.Sprintf(`%s://%s:%s@%s:%d/%s?sslmode=disable`, dbConf.Type, dbConf.User, dbConf.Pass, dbConf.Host, dbConf.Port, dbConf.Name)
	db, err := sqlx.Open(dbConf.Type, connStr)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	return db
}
