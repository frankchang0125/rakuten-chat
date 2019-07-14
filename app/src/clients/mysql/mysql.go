package mysql

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const timeout = "10s"
const readTimeout = "30s"
const writeTimeout = "30s"

var db *sql.DB

// InitClient initialize MySQL client.
// Note: This not thread-safe and should be called at server startup.
func InitClient() (*sql.DB, error) {
	if db == nil {
		dbName := "rakuten"
		host := viper.GetString("mysql.host")
		port := viper.GetInt("mysql.port")
		username := viper.GetString("mysql.username")
		password := viper.GetString("mysql.password")
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?timeout=%s&readTimeout=%s&writeTimeout=%s&parseTime=true&loc=Local",
			username,
			password,
			host,
			port,
			dbName,
			timeout,
			readTimeout,
			writeTimeout,
		)

		log.WithFields(log.Fields{
			"dsn": dsn,
		}).Info("Initialize MySQL client")

		_db, err := sql.Open("mysql", dsn)
		if err != nil {
			log.WithError(err).Error("Fail to initialize MySQL client")
			return nil, err
		}

		db = _db

		// Max open connections: 5
		// Max idle connections: 2 (Default value)
		db.SetMaxOpenConns(5)
	}

	return db, nil
}
