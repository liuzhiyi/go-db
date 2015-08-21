package adapter

import "database/sql"

/**
*
*注意：默认为mysql一类的适配器
**/
type Adapter interface {
	create(driverName, dsn string) Adapter
	Connect()
	Close()
	BeginTransaction() *Transaction
	SetTransaction(t *Transaction)
	GetDb() *sql.DB
	QueryRow(sql string, bind ...interface{}) *sql.Row
	Query(sql string, bind ...interface{}) *sql.Rows
	Exec(sql string, bind ...interface{}) (sql.Result, error)
	Prepare(sql string) *sql.Stmt
	Insert(table string, bind map[string]interface{}) (int64, error)
	Update(table string, bind map[string]interface{}, where string) (int64, error)
	Delete(table, where string) (int64, error)
	GetTableName(name string) string
	QuoteIdentifierAs(ident, alias string) string
	QuoteIdentifier(value string) string
	GetQuoteIdentifierSymbol() string
	QuoteInto(text string, values ...interface{}) string
	Quote(value interface{}) string
	Limit(sql string, count, offset int64) string
	MustExec(sql string, bind ...interface{})
}

var adapters = make(map[string]Adapter)

func register(name string, adapter Adapter) {
	if adapter == nil {
		panic("go-db: Register adapter is nil")
	}
	if _, dup := adapters[name]; dup {
		panic("go-db: Register called twice for adapter " + name)
	}
	adapters[name] = adapter
}

func NewAdapter(driverName, dsn string) Adapter {
	adapter, ok := adapters[driverName]
	if !ok {
		panic("unknown adapter:" + driverName)
	}

	return adapter.create(driverName, dsn)
}
