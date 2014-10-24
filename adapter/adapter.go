package adapter

import (
	"database/sql"
)

type adpaterAbstract struct {
	db *sql.DB
    tx *sql.Tx
}

func (a *adpaterAbstract) _query() {

}

func (a *adpaterAbstract) _beginTransaction() {
    var err error
    if a.tx, err = a.db.Begin(); err != nil {
        panic(err.Error())
    }
}

func (a *adpaterAbstract) _rollBack() {
    a.tx.Rollback()
}

func (a *adpaterAbstract) _commit() {
    a.tx.Commit()
}
