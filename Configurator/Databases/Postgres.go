package Databases

import (
	"database/sql"
	_ "github.com/lib/pq"
)

type Postgres struct {
	Driver       string `json:"driver"`
	Host         string `json:"host"`
	Port         string `json:"port"`
	Replica      string `json:"replica"`
	DatabaseName string `json:"name"`
	User         string `json:"user"`
	Password     string `json:"password"`
	Connection   *sql.DB
}

func (pdb *Postgres) IsEmpty() bool {
	return pdb.DatabaseName == ""
}

func (pdb *Postgres) Initial() error {
	if !pdb.IsEmpty() {
		var err error
		connStr := pdb.Driver + "://" + pdb.User + ":" + pdb.Password + "@"
		if pdb.Port != "" {
			connStr = connStr + ":" + pdb.Port
		}
		connStr = connStr + "/" + pdb.DatabaseName + "?sslmode=disable"
		pdb.Connection, err = sql.Open("postgres", connStr)
		if err != nil {
			return err
		}
		if err = pdb.Connection.Ping(); err != nil {
			return err
		}
	}
	return nil
}

func (pdb *Postgres) GetConnection() (*sql.DB, error) {
	return pdb.Connection, nil
}

func (pdb *Postgres) GetDatabase() (*sql.DB, error) {
	return pdb.Connection, nil
}

func (pdb *Postgres) Close() error {
	return pdb.Connection.Close()
}
