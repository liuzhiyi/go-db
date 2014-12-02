package adapter

import (
	"database/sql"
	"fmt"
	"strings"

	db "github.com/liuzhiyi/go-db"
)

type Abstract struct {
	db *sql.DB
	tx *sql.Tx
}

func (a *Abstract) _query() {

}

func (a *Abstract) BeginTransaction() {
	var err error
	if a.tx, err = a.db.Begin(); err != nil {
		panic(err.Error())
	}
}

func (a *Abstract) RollBack() {
	a.tx.Rollback()
}

func (a *Abstract) Commit() {
	a.tx.Commit()
}

func (a *Abstract) GetAdapter() *sql.DB {
	return a.db
}

func (a *Abstract) Query(sql interface{}, bind ...[]string) sql.Result {
	var str string
	switch sql.(type) {
	case string:
		str = sql.(string)
	case *db.Select:
		str = (sql.(*db.Select)).Assemble()
	default:
		panic("invalid sql type")
	}
	stmt, err := a.db.Prepare(str)
	if err != nil {
		panic(err.Error())
	}
	result, _ := stmt.Exec(bind)
	return result
}

func (a *Abstract) Insert(table string, bind map[string]string) int64 {
	var cols, vals, quotes []string
	for col, val := range bind {
		cols = append(cols, col)
		quotes = append(quotes, "?")
		vals = append(vals, val)
	}
	sql := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", table, strings.Join(cols, ","), strings.Join(quotes, ","))
	result := a.Query(sql, vals)
	rows, err := result.RowsAffected()
	if err != nil {
		panic(err.Error())
	}
	return rows
}

func (a *Abstract) Update(table string, bind map[string]string, where string) int64 {
	var sets, vals []string
	for col, val := range bind {
		sets = append(sets, fmt.Sprintf("%s = ?", col))
		vals = append(vals, val)
	}
	sql := fmt.Sprintf("UPDATE %s SET %s WHERE %s", table, strings.Join(sets, ","), where)
	result := a.Query(sql, vals)
	rows, err := result.RowsAffected()
	if err != nil {
		panic(err.Error())
	}
	return rows
}

func (a *Abstract) Delete(table, where string) int64 {
	sql := fmt.Sprintf("DELETE FROM %s WHERE %s", table, where)
	result := a.Query(sql)
	rows, err := result.RowsAffected()
	if err != nil {
		panic(err.Error())
	}
	return rows
}

func (a *Abstract) QuoteInto(text string, value string) string {
	return strings.Replace(text, "?", a.Quote(value), 0)
}

func (a *Abstract) Quote(value interface{}) string {
	return a._quote(value)
}

/*
 Quote a raw string.
*/
func (a *Abstract) _quote(value interface{}) string {
	switch value.(type) {
	case int, int16, int32, int64, int8:
		return fmt.Sprintf("%d", value)
	case float32, float64:
		return fmt.Sprintf("%F", value)
	case string:
		return "'" + AddSlashes(value, "\000\n\r\\'\"\032") + "'"
	default:
		panic("Invalid value")
	}
}
