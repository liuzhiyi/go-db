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

const maxTransationId = 0xffffffff

/**
*
*注意：默认为mysql一类的适配器
**/
type Mysql struct {
	db              *sql.DB
	transaction     map[uint64]*Transaction
	transactionId   uint64
	transactionLock sync.RWMutex
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
	m.transactionId = 1
	m.transaction = make(map[uint64]*Transaction)
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
	for _, t := range m.transaction {
		t.Rollback()
	}

	if err := m.db.Close(); err != nil {
		panic(err.Error())
	}
}

func (m *Mysql) getTransactionId() uint64 {
	var id uint64

	m.transactionLock.RLock()
	defer m.transactionLock.RUnlock()
	for {
		id = m.transactionId

		if m.transactionId == maxTransationId {
			m.transactionId = 1
		} else {
			m.transactionId++
		}

		if _, exists := m.transaction[id]; !exists {
			break
		}

	}

	return id
}

/**
*
*建议一般情况下开启事务机制
*****/
func (m *Mysql) BeginTransaction() (t *Transaction) {

	if tx, err := m.db.Begin(); err != nil {
		panic(err.Error())
	} else {
		id := m.getTransactionId()
		t = newTransaction(tx, m, id)
		m.transactionLock.Lock()
		m.transaction[id] = t
		m.transactionLock.Unlock()
	}

	return t
}

func (m *Mysql) GetDb() *sql.DB {
	return m.db
}

func (m *Mysql) QueryRow(sql string, bind ...interface{}) (*sql.Row, error) {
	stmt, err := m.Prepare(sql)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	row := stmt.QueryRow(bind...)
	return row, nil
}

func (m *Mysql) Query(sql string, bind ...interface{}) (*sql.Rows, error) {
	stmt, err := m.Prepare(sql)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	return stmt.Query(bind...)
}

func (m *Mysql) Exec(sqlStr string, bind ...interface{}) (sql.Result, error) {
	stmt, err := m.Prepare(sqlStr)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(bind...)
	return result, err
}

func (m *Mysql) RawQuery(sql string, args ...interface{}) (*sql.Rows, error) {

	return m.db.Query(sql, args...)
}

func (m *Mysql) RawQueryRow(sql string, args ...interface{}) (*sql.Row, error) {

	return m.db.QueryRow(sql, args...), nil
}

func (m *Mysql) RawExec(sql string, args ...interface{}) (sql.Result, error) {

	return m.db.Exec(sql, args...)
}

func (m *Mysql) Prepare(query string) (*sql.Stmt, error) {
	return m.db.Prepare(query)
}

func (m *Mysql) MustExec(query string, bind ...interface{}) {
	var err error

	_, err = m.db.Exec(query, bind...)

	if err != nil {
		panic(err.Error())
	}
}

func (m *Mysql) Insert(table string, bind map[string]interface{}) (int64, error) {
	var cols, quotes []string
	var vals []interface{}
	for col, val := range bind {
		col = m.QuoteIdentifier(col)
		cols = append(cols, col)
		quotes = append(quotes, "?")
		vals = append(vals, val)
	}
	sql := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", table, strings.Join(cols, ","), strings.Join(quotes, ","))
	if result, err := m.RawExec(sql, vals...); err != nil {
		return 0, err
	} else {
		return result.LastInsertId()
	}
}

func (m *Mysql) Update(table string, bind map[string]interface{}, where string) (int64, error) {
	var sets []string
	var vals []interface{}
	for col, val := range bind {
		col = m.QuoteIdentifier(col)
		sets = append(sets, fmt.Sprintf("%s = ?", col))
		vals = append(vals, val)
	}
	sql := fmt.Sprintf("UPDATE %s SET %s WHERE %s", table, strings.Join(sets, ","), where)
	if result, err := m.RawExec(sql, vals...); err != nil {
		return 0, err
	} else {
		return result.RowsAffected()
	}
}

func (m *Mysql) Delete(table, where string) (int64, error) {
	sql := fmt.Sprintf("DELETE FROM %s WHERE %s", table, where)
	if result, err := m.RawExec(sql); err != nil {
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
	case int, int16, int32, int64, int8, uint32, uint64, uint16, uint, uint8:
		return fmt.Sprintf("%d", value)
	case float32, float64:
		return fmt.Sprintf("%F", value)
	case string:
		return "'" + str.AddSlashes(value.(string), "\000\n\r\\'\"\032") + "'"
	default:
		fmt.Println("%T",value)
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
