package dbrepo

import (
	"database/sql"
	"web3/pkg/config"
	"web3/pkg/repository"
)

type postgresDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

func NewPostgresRepo(conn *sql.DB, ac *config.AppConfig) repository.DatabaseRepo {
	return &postgresDBRepo{
		App: ac,
		DB:  conn,
	}
}
