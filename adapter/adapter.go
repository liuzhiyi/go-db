package adapter

import (
	"database/sql"
	"fmt"
	"strings"

	db "github.com/liuzhiyi/go-db"
)

type Adapter struct {
	db         *sql.DB
	tx         *sql.Tx
	driverName string
	config     string
}

func NewAdapter(driverName, dsn string) *Adapter {
	a := new(Adapter)
	a.Init(driverName, dsn)
	return a
}

func (a *Adapter) Init(driverName, dsn string) {
	a.driverName = driverName
	a.config = dsn
	a._init()
}

func (a *Adapter) _init() {
	a.connect()
}

func (a *Adapter) _query() {

}

func (a *Adapter) connect() {
	if a.db != nil {
		return
	}
	var err error
	a.db, err = sql.Open(a.driverName, a.config)
	if err != nil {
		panic(err.Error())
	}
}

func (a *Adapter) BeginTransaction() {
	var err error
	if a.tx, err = a.db.Begin(); err != nil {
		panic(err.Error())
	}
}

func (a *Adapter) RollBack() {
	a.tx.Rollback()
}

func (a *Adapter) Commit() {
	a.tx.Commit()
}

func (a *Adapter) GetAdapter() *sql.DB {
	return a.db
}

func (a *Adapter) QueryRow(sql interface{}, bind ...[]string) *sql.Row {
	stmt := a.prepare(sql)
	defer stmt.Close()
	row := stmt.QueryRow(bind)
	return row
}

func (a *Adapter) Query(sql interface{}, bind ...[]string) *sql.Rows {
	stmt := a.prepare(sql)
	defer stmt.Close()
	rows, _ := stmt.Query()
	return rows
}

func (a *Adapter) Exec(sql interface{}, bind ...[]string) sql.Result {
	stmt := a.prepare(sql)
	defer stmt.Close()
	result, _ := stmt.Exec(bind)
	return result
}

func (a *Adapter) prepare(sql interface{}) *sql.Stmt {
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
	return stmt
}

func (a *Adapter) Insert(table string, bind map[string]string) int64 {
	var cols, vals, quotes []string
	for col, val := range bind {
		cols = append(cols, col)
		quotes = append(quotes, "?")
		vals = append(vals, val)
	}
	sql := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", table, strings.Join(cols, ","), strings.Join(quotes, ","))
	result := a.Exec(sql, vals)
	rows, err := result.RowsAffected()
	if err != nil {
		panic(err.Error())
	}
	return rows
}

func (a *Adapter) Update(table string, bind map[string]string, where string) int64 {
	var sets, vals []string
	for col, val := range bind {
		sets = append(sets, fmt.Sprintf("%s = ?", col))
		vals = append(vals, val)
	}
	sql := fmt.Sprintf("UPDATE %s SET %s WHERE %s", table, strings.Join(sets, ","), where)
	result := a.Exec(sql, vals)
	rows, err := result.RowsAffected()
	if err != nil {
		panic(err.Error())
	}
	return rows
}

func (a *Adapter) Delete(table, where string) int64 {
	sql := fmt.Sprintf("DELETE FROM %s WHERE %s", table, where)
	result := a.Exec(sql)
	rows, err := result.RowsAffected()
	if err != nil {
		panic(err.Error())
	}
	return rows
}

func (a *Adapter) QuoteInto(text string, value string) string {
	return strings.Replace(text, "?", a.Quote(value), 0)
}

func (a *Adapter) Quote(value interface{}) string {
	return a._quote(value)
}

/*
 Quote a raw string.
*/
func (a *Adapter) _quote(value interface{}) string {
	switch value.(type) {
	case int, int16, int32, int64, int8:
		return fmt.Sprintf("%d", value)
	case float32, float64:
		return fmt.Sprintf("%F", value)
	case string:
		return "'" + db.AddSlashes(value.(string), "\000\n\r\\'\"\032") + "'"
	default:
		panic("Invalid value")
	}
}
