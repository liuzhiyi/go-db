package adapter

import (
	"database/sql"
)

type Transaction struct {
	tx    *sql.Tx
	level int
}

func newTransaction(tx *sql.Tx) *Transaction {
	return &Transaction{
		tx:    tx,
		level: 1,
	}
}

func (t *Transaction) Begin() *Transaction {
	t.level++

	return t
}

func (t *Transaction) Commit() error {
	if t.level == 0 {
		return t.tx.Commit()
	}

	t.level--

	return nil
}

func (t *Transaction) Rollback() error {
	t.level = 0
	return t.tx.Rollback()
}

func (t *Transaction) Query(sql string, args ...interface{}) (*sql.Rows, error) {
	return t.tx.Query(sql, args...)
}

func (t *Transaction) Prepare(sql string) (*sql.Stmt, error) {
	return t.tx.Prepare(sql)
}

func (t *Transaction) QueryRow(sql string, args ...interface{}) *sql.Row {
	return t.tx.QueryRow(sql, args...)
}

func (t *Transaction) Exec(sql string, args ...interface{}) (sql.Result, error) {
	return t.tx.Exec(sql, args...)
}

func (t *Transaction) IsOver() bool {
	return t.level == 0
}
