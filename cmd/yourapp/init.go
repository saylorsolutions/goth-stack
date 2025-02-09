package main

import (
	"database/sql"
	"errors"
	"github.com/saylorsolutions/x/env"
	"log"
	"yourapp/feature/audit"
	"yourapp/feature/auth"
)

func initData() (*sql.DB, error) {
	connectionURL := env.Val("DBURL", "")
	if len(connectionURL) == 0 {
		return nil, errors.New("no connection URL set")
	}
	return sql.Open("pgx", connectionURL)
}

func initAuth(logger *log.Logger, db *sql.DB) (*auth.Service, error) {
	auditLog := audit.NewLogger(db, audit.StdDelegate(logger, true))
	return auth.NewAuthService(auditLog, db)
}
