package driver

import (
	"database/sql"
	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"time"
)

// DB holds the database connection pool
type DB struct {
	SQL *sql.DB
}

var dbConn = &DB{}

const maxOpenConn = 10
const maxIdleConn = 5
const maxConnLifetime = 5 * time.Minute

func ConnectSQL(dsn string) (*DB, error) {
	d, err := NewDatabase(dsn)
	if err != nil {
		panic(err)
	}

	d.SetMaxOpenConns(maxOpenConn)
	d.SetMaxIdleConns(maxIdleConn)
	d.SetConnMaxLifetime(maxConnLifetime)

	dbConn.SQL = d

	err = testDB(d)
	if err != nil {
		return nil, err
	}

	return dbConn, nil
}

func testDB(d *sql.DB) error {
	err := d.Ping()
	if err != nil {
		return err
	}
	return nil
}

func NewDatabase(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil

}
