package db

import "database/sql"

var (
	otherDB = new(OtherDB)
)

type OtherDB struct {
	PDB
}

func NewOtherDB(db *sql.DB) {
	otherDB.DB = db
}

func init() {
	name := "db/other"
	RegisterDB(name, otherDB)
}
