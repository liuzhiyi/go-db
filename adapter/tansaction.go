package adapter

import (
	"database/sql"
	"fmt"
	"strings"
)

var TransactionOveredErr = fmt.Errorf("using a transaction has been overed")

type Transaction struct {
	tx     *sql.Tx
	adpter *Mysql
	id     uint64
	level  int
}

func newTransaction(tx *sql.Tx, adpter *Mysql, id uint64) *Transaction {
	return &Transaction{
		tx:     tx,
		adpter: adpter,
		id:     id,
		level:  1,
	}
}

func (t *Transaction) GetId() uint64 {
	return t.id
}

func (t *Transaction) Begin() *Transaction {
	t.level++

	return t
}

func (t *Transaction) Commit() error {
	if t.IsOver() {
		return TransactionOveredErr
	}

	t.level--

	if t.level == 0 {
		t.close()
		return t.tx.Commit()
	}

	return nil
}

func (t *Transaction) Rollback() error {
	if t.IsOver() {
		return TransactionOveredErr
	}

	t.level = 0
	t.close()
	return t.tx.Rollback()
}

func (t *Transaction) RawQuery(sql string, args ...interface{}) (*sql.Rows, error) {
	if t.IsOver() {
		return nil, TransactionOveredErr
	}

	return t.tx.Query(sql, args...)
}

func (t *Transaction) RawQueryRow(sql string, args ...interface{}) (*sql.Row, error) {
	if t.IsOver() {
		return nil, TransactionOveredErr
	}

	return t.tx.QueryRow(sql, args...), nil
}

func (t *Transaction) RawExec(sql string, args ...interface{}) (sql.Result, error) {
	if t.IsOver() {
		return nil, TransactionOveredErr
	}

	return t.tx.Exec(sql, args...)
}

func (t *Transaction) Query(sql string, args ...interface{}) (*sql.Rows, error) {
	if t.IsOver() {
		return nil, TransactionOveredErr
	}

	stmt, err := t.Prepare(sql)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(args...)
	return rows, err
}

func (t *Transaction) Prepare(sql string) (*sql.Stmt, error) {
	if t.IsOver() {
		return nil, TransactionOveredErr
	}

	return t.tx.Prepare(sql)
}

func (t *Transaction) QueryRow(sql string, args ...interface{}) (*sql.Row, error) {
	if t.IsOver() {
		return nil, TransactionOveredErr
	}

	stmt, err := t.Prepare(sql)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	row := stmt.QueryRow(args...)
	return row, nil
}

func (t *Transaction) Exec(sqlStr string, args ...interface{}) (sql.Result, error) {
	if t.IsOver() {
		return nil, TransactionOveredErr
	}

	stmt, err := t.Prepare(sqlStr)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(args...)
	return result, err
}

func (t *Transaction) Insert(table string, bind map[string]interface{}) (int64, error) {
	var cols, quotes []string
	var vals []interface{}
	for col, val := range bind {
		cols = append(cols, col)
		quotes = append(quotes, "?")
		vals = append(vals, val)
	}
	sql := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", table, strings.Join(cols, ","), strings.Join(quotes, ","))
	if result, err := t.RawExec(sql, vals...); err != nil {
		return 0, err
	} else {
		return result.LastInsertId()
	}
}

func (t *Transaction) Update(table string, bind map[string]interface{}, where string) (int64, error) {
	var sets []string
	var vals []interface{}
	for col, val := range bind {
		sets = append(sets, fmt.Sprintf("%s = ?", col))
		vals = append(vals, val)
	}
	sql := fmt.Sprintf("UPDATE %s SET %s WHERE %s", table, strings.Join(sets, ","), where)
	if result, err := t.RawExec(sql, vals...); err != nil {
		return 0, err
	} else {
		return result.RowsAffected()
	}
}

func (t *Transaction) Delete(table, where string) (int64, error) {
	sql := fmt.Sprintf("DELETE FROM %s WHERE %s", table, where)
	if result, err := t.RawExec(sql); err != nil {
		return 0, err
	} else {
		return result.RowsAffected()
	}
}

func (t *Transaction) IsOver() bool {
	return t.level == 0
}

func (t *Transaction) close() {
	delete(t.adpter.transaction, t.GetId())
}
