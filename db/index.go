package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"time"
)

type Connection struct {
	Host        string
	DbName      string
	User        string
	Password    string
	Port        string
	Location    *time.Location
	SslMode     string
	SslCert     string
	SslKey      string
	SslRootCert string
	DBMaxConnection         int
	DBMAxIdleConnection     int
	DBMaxLifeTimeConnection int
}

func (c Connection) DbConnect() (*sql.DB, error) {
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:5432/%s?sslmode=%s&TimeZone=UTC", c.User, c.Password, c.Host, c.DbName, c.SslMode,
	)

	if c.SslMode == "require" {
		connStr = fmt.Sprintf(
			"postgres://%s:%s@%s:%d/%s?sslmode=%s&TimeZone=UTC&sslcert=%s&sslkey=%s&sslrootcert=%s",
			c.User, c.Password, c.Host, c.Port, c.DbName, c.SslMode, c.SslCert, c.SslKey, c.SslRootCert,
		)
	}
	db, err := sql.Open("postgres", connStr)
	db.SetMaxOpenConns(c.DBMaxConnection)
	db.SetMaxIdleConns(c.DBMAxIdleConnection)
	db.SetConnMaxLifetime(time.Duration(c.DBMaxLifeTimeConnection)*time.Second)

	// migrations := &migrate.FileMigrationSource{
	// 	Dir: "../db/migrations",
	// }
	// n, err := migrate.Exec(db, "postgres", migrations, migrate.Up)
	// if err != nil {
	// 	log.Fatal("Error migration := ", err.Error())
	// }
	// fmt.Printf("Applied %d migrations!\n", n)

	return db, err
}
