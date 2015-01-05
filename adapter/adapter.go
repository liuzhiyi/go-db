package adapter

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/liuzhiyi/utils/str"
)

/**
*
*注意：默认为mysql一类的适配器
**/
type Adapter struct {
	db               *sql.DB
	tx               *sql.Tx
	prefix           string
	driverName       string
	config           string
	transactionLevel int
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

func (a *Adapter) Close() {
	if err := a.db.Close(); err != nil {
		panic(err.Error())
	}
}

/**
*
*建议一般情况下开启事务机制
*****/
func (a *Adapter) BeginTransaction() {
	if a.transactionLevel == 0 {
		var err error
		if a.tx, err = a.db.Begin(); err != nil {
			panic(err.Error())
		}
	}
	a.transactionLevel++
}

func (a *Adapter) RollBack() {
	a.tx.Rollback()
	a.transactionLevel = 0
}

func (a *Adapter) Commit() {
	if a.transactionLevel == 1 {
		a.tx.Commit()
	}
	a.transactionLevel--
}

func (a *Adapter) GetTransactionLevel() int {
	return a.transactionLevel
}

func (a *Adapter) GetDb() *sql.DB {
	return a.db
}

func (a *Adapter) QueryRow(sql string, bind ...string) *sql.Row {
	stmt := a.prepare(sql)
	defer stmt.Close()
	row := stmt.QueryRow()
	return row
}

func (a *Adapter) Query(sql string, bind ...string) *sql.Rows {
	stmt := a.prepare(sql)
	defer stmt.Close()
	rows, _ := stmt.Query()
	return rows
}

func (a *Adapter) Exec(sql string, bind ...string) (sql.Result, error) {
	stmt := a.prepare(sql)
	defer stmt.Close()
	result, err := stmt.Exec(bind)
	return result, err
}

func (a *Adapter) prepare(sql string) *sql.Stmt {
	stmt, err := a.db.Prepare(sql)
	if err != nil {
		panic(err.Error())
	}
	return stmt
}

func (a *Adapter) Insert(table string, bind map[string]string) (int64, error) {
	var cols, vals, quotes []string
	for col, val := range bind {
		cols = append(cols, col)
		quotes = append(quotes, "?")
		vals = append(vals, val)
	}
	sql := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", table, strings.Join(cols, ","), strings.Join(quotes, ","))
	if result, err := a.Exec(sql, vals...); err != nil {
		return 0, err
	} else {
		return result.RowsAffected()
	}
}

func (a *Adapter) Update(table string, bind map[string]string, where string) (int64, error) {
	var sets, vals []string
	for col, val := range bind {
		sets = append(sets, fmt.Sprintf("%s = ?", col))
		vals = append(vals, val)
	}
	sql := fmt.Sprintf("UPDATE %s SET %s WHERE %s", table, strings.Join(sets, ","), where)
	if result, err := a.Exec(sql, vals...); err != nil {
		return 0, err
	} else {
		return result.RowsAffected()
	}
}

func (a *Adapter) Delete(table, where string) (int64, error) {
	sql := fmt.Sprintf("DELETE FROM %s WHERE %s", table, where)
	if result, err := a.Exec(sql); err != nil {
		return 0, err
	} else {
		return result.RowsAffected()
	}
}

func (a *Adapter) GetTableName(name string) string {
	if a.prefix != "" {
		return a.prefix + "_" + name
	}
	return name
}

func (a *Adapter) QuoteIdentifierAs(ident, alias string) string {
	as := " AS "
	idents := strings.Split(ident, ".")
	for i := 0; i < len(idents); i++ {
		if idents[i] == "*" {
			continue
		}
		idents[i] = a._quoteIdentifier(idents[i])
	}
	quoted := strings.Join(idents, ".")
	if alias != "" {
		quoted += as + a._quoteIdentifier(alias)
	}
	return quoted
}

func (a *Adapter) QuoteIdentifier(value string) string {
	return a.QuoteIdentifierAs(value, "")
}

func (a *Adapter) _quoteIdentifier(value string) string {
	q := a.GetQuoteIdentifierSymbol()
	return q + (strings.Replace(value, q, q+q, -1)) + q
}

func (a *Adapter) GetQuoteIdentifierSymbol() string {
	return "`"
}

func (a *Adapter) QuoteInto(text string, value interface{}) string {
	return strings.Replace(text, "?", a.Quote(value), -1)
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
		return "'" + str.AddSlashes(value.(string), "\000\n\r\\'\"\032") + "'"
	default:
		panic("Invalid value")
	}
}

func (a *Adapter) Limit(sql string, count, offset int64) string {
	if count <= 0 {
		panic(fmt.Sprintf("LIMIT argument count=%s is not valid", count))
	}
	if offset < 0 {
		panic(fmt.Sprintf("LIMIT argument offset=%s is not valid", offset))
	}
	sql += " LIMIT " + strconv.FormatInt(offset, 10)
	sql += ", " + strconv.FormatInt(count, 10)

	return sql
}
