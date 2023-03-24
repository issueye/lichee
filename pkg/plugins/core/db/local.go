package db

import "database/sql"

var (
	localDB = new(LocalDB)
)

type LocalDB struct {
	PDB
}

func NewLocalDB(db *sql.DB) {
	localDB.DB = db
}

func init() {
	name := "db/local"
	RegisterDB(name, localDB)
}
