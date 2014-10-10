package adapter

import (
	"database/sql"
)

type adpaterAbstract struct {
	db *sql.DB
}
