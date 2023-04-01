package db

import "database/sql"

var (
	Ldb = new(LocalDB)
)

type LocalDB struct {
	PDB
}

func NewLocalDB(db *sql.DB) {
	Ldb.DB = db
}

func init() {
	name := "db/local"
	RegisterDB(name, Ldb)
}
