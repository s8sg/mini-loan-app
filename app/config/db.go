package config

import (
	"database/sql"
)

type DbConfig struct {
	User     string
	Password string
	Host     string
	DBName   string
}

func InitialiseDB(config DbConfig) (*sql.DB, error) {

	connectionUrl := getDbConnectionUrlFromConfig(config)

	db, err := sql.Open("postgres", connectionUrl)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func getDbConnectionUrlFromConfig(config DbConfig) string {
	return "postgres://" + config.User + ":" + config.Password + "@" + config.Host + "/" + config.DBName + "?sslmode=disable"
}
