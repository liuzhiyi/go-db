package adapter

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/liuzhiyi/utils/str"
)

func init() {
	register("mysql", &Mysql{})
}

/**
*
*注意：默认为mysql一类的适配器
**/
type Mysql struct {
	db              *sql.DB
	transaction     *Transaction
	isTransaction   bool
	transactionLock sync.Mutex
	prefix          string
	driverName      string
	config          string
}

func (m *Mysql) create(driverName, dsn string) Adapter {
	a := new(Mysql)
	a._init(driverName, dsn)
	return a
}

func (m *Mysql) _init(driverName, dsn string) {
	m.driverName = driverName
	m.config = dsn
	m.Connect()
}

func (m *Mysql) _query() {

}

func (m *Mysql) Connect() {
	if m.db != nil {
		return
	}
	var err error
	m.db, err = sql.Open(m.driverName, m.config)
	if err != nil {
		panic(err.Error())
	}
}

func (m *Mysql) Close() {
	if err := m.db.Close(); err != nil {
		panic(err.Error())
	}
}

/**
*
*建议一般情况下开启事务机制
*****/
func (m *Mysql) BeginTransaction() (t *Transaction) {
	if tx, err := m.db.Begin(); err != nil {
		panic(err.Error())
	} else {
		t = newTransaction(tx)
	}

	return
}

func (m *Mysql) SetTransaction(t *Transaction) {
	if t == nil {
		return
	}

	m.transactionLock.Lock()
	m.isTransaction = true
	m.transaction = t
}

func (m *Mysql) freeTransaction() {
	m.isTransaction = false

	m.transactionLock.Unlock()
}

func (m *Mysql) GetDb() *sql.DB {
	return m.db
}

func (m *Mysql) QueryRow(sql string, bind ...interface{}) *sql.Row {
	stmt := m.Prepare(sql)
	defer stmt.Close()
	row := stmt.QueryRow(bind...)
	return row
}

func (m *Mysql) Query(sql string, bind ...interface{}) *sql.Rows {
	stmt := m.Prepare(sql)
	defer stmt.Close()

	rows, err := stmt.Query(bind...)
	if err != nil {
		panic(err.Error())
	}
	return rows
}

func (m *Mysql) Exec(sql string, bind ...interface{}) (sql.Result, error) {
	stmt := m.Prepare(sql)
	defer stmt.Close()
	result, err := stmt.Exec(bind...)
	return result, err
}

func (m *Mysql) Prepare(query string) *sql.Stmt {
	var stmt *sql.Stmt
	var err error

	if m.isTransaction {
		stmt, err = m.transaction.Prepare(query)
		fmt.Println(err)
		defer m.freeTransaction()
	} else {
		stmt, err = m.db.Prepare(query)
	}

	if err != nil {
		panic(err.Error())
	}
	return stmt
}

func (m *Mysql) MustExec(query string, bind ...interface{}) {
	var err error
	if m.isTransaction {
		_, err = m.transaction.Exec(query, bind...)
		defer m.freeTransaction()
	} else {
		_, err = m.db.Exec(query, bind...)
	}

	if err != nil {
		panic(err.Error())
	}
}

func (m *Mysql) Insert(table string, bind map[string]interface{}) (int64, error) {
	var cols, quotes []string
	var vals []interface{}
	for col, val := range bind {
		cols = append(cols, col)
		quotes = append(quotes, "?")
		vals = append(vals, val)
	}
	sql := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", table, strings.Join(cols, ","), strings.Join(quotes, ","))
	if result, err := m.Exec(sql, vals...); err != nil {
		return 0, err
	} else {
		return result.LastInsertId()
	}
}

func (m *Mysql) Update(table string, bind map[string]interface{}, where string) (int64, error) {
	var sets []string
	var vals []interface{}
	for col, val := range bind {
		sets = append(sets, fmt.Sprintf("%s = ?", col))
		vals = append(vals, val)
	}
	sql := fmt.Sprintf("UPDATE %s SET %s WHERE %s", table, strings.Join(sets, ","), where)
	if result, err := m.Exec(sql, vals...); err != nil {
		return 0, err
	} else {
		return result.RowsAffected()
	}
}

func (m *Mysql) Delete(table, where string) (int64, error) {
	sql := fmt.Sprintf("DELETE FROM %s WHERE %s", table, where)
	if result, err := m.Exec(sql); err != nil {
		return 0, err
	} else {
		return result.RowsAffected()
	}
}

func (m *Mysql) GetTableName(name string) string {
	if m.prefix != "" {
		return m.prefix + "_" + name
	}
	return name
}

func (m *Mysql) QuoteIdentifierAs(ident, alias string) string {
	as := " AS "
	idents := strings.Split(ident, ".")
	for i := 0; i < len(idents); i++ {
		if idents[i] == "*" {
			continue
		}
		idents[i] = m._quoteIdentifier(idents[i])
	}
	quoted := strings.Join(idents, ".")
	if alias != "" {
		quoted += as + m._quoteIdentifier(alias)
	}
	return quoted
}

func (m *Mysql) QuoteIdentifier(value string) string {
	return m.QuoteIdentifierAs(value, "")
}

func (m *Mysql) _quoteIdentifier(value string) string {
	q := m.GetQuoteIdentifierSymbol()
	return q + (strings.Replace(value, q, q+q, -1)) + q
}

func (m *Mysql) GetQuoteIdentifierSymbol() string {
	return "`"
}

func (m *Mysql) QuoteInto(text string, values ...interface{}) string {
	for i := 0; i < len(values); i++ {
		text = strings.Replace(text, "?", m.Quote(values[i]), 1)
	}
	return text
}

func (m *Mysql) Quote(value interface{}) string {
	return m._quote(value)
}

/*
 Quote m raw string.
*/
func (m *Mysql) _quote(value interface{}) string {
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

func (m *Mysql) Limit(sql string, count, offset int64) string {
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
