package dbrepo

import (
	"database/sql"
	"github.com/indexcoder/bookings/internal/config"
	"github.com/indexcoder/bookings/internal/repository"
)

type postgresRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

type testDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

func NewPostgresRepo(app *config.AppConfig, db *sql.DB) repository.DatabaseRepo {
	return &postgresRepo{App: app, DB: db}
}

func NewTestingDBRepo(app *config.AppConfig) repository.DatabaseRepo {
	return &testDBRepo{App: app}
}
