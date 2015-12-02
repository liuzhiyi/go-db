package db

import (
	"database/sql"
	"fmt"
	"sort"
	"strings"

	"github.com/liuzhiyi/go-db/adapter"
	"github.com/liuzhiyi/go-db/table"
)

type Resource struct {
	idField   string
	mainTable string
	fields    []string
}

func NewResource(table, idField string) *Resource {
	r := new(Resource)
	r._setMainTable(table, idField)
	r.setFields()
	r.sortFields()

	return r
}

func (r *Resource) GetFields() []string {
	if len(r.fields) == 0 {
		r.setFields()
	}

	return r.fields
}

func (r *Resource) setFields() *Resource {
	sql := fmt.Sprintf("SHOW COLUMNS FROM `%s`", r.mainTable)
	rows, err := r.GetReadAdapter().RawQuery(sql)
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()

	clm, err := rows.Columns()
	if err != nil {
		panic(err.Error())
	}

	item := make([]interface{}, len(clm))
	for i, _ := range item {
		item[i] = new(string)
	}

	for rows.Next() {
		err = rows.Scan(item...)

		field := item[0].(*string)
		r.fields = append(r.fields, *field)
	}

	return r
}

func (r *Resource) sortFields() *Resource {
	a := sort.StringSlice(r.fields[0:])
	sort.Sort(a)
	r.fields = a

	return r
}

func (r *Resource) IsExistField(name string) bool {
	i := sort.SearchStrings(r.fields, name)

	return r.fields[i] == name
}

func (r *Resource) GetIdName() string {
	return r.idField
}

func (r *Resource) BeginTransaction() *adapter.Transaction {
	return r.GetWriteAdapter().BeginTransaction()

}

func (r *Resource) GetMainTable() string {
	return r.GetReadAdapter().GetTableName(r.mainTable)
}

func (r *Resource) _setMainTable(table, idField string) {
	r.mainTable = table
	if idField == "" {
		idField = fmt.Sprintf("%s_id", table)
	}
	r.idField = idField
}

func (r *Resource) GetTable(name string) string {
	return r.GetReadAdapter().GetTableName(name)
}

func (r *Resource) Load(item *Item, id int) {
	var (
		rows *sql.Rows
		err  error
	)

	read := r.GetReadAdapter()
	field := r.GetIdName()
	sql := r._getLoadSelect(field, id)

	transaction := item.GetTransaction()
	if transaction != nil {
		rows, err = transaction.RawQuery(sql.Assemble())
	} else {
		rows, err = read.RawQuery(sql.Assemble())
	}

	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		r._fetch(rows, item)
		return
	}
}

func (r *Resource) fetchOne(transaction *adapter.Transaction, sqlStr string, dest interface{}) {
	var (
		row *sql.Row
		err error
	)

	if transaction != nil {
		row, err = transaction.RawQueryRow(sqlStr)
	} else {
		row, err = r.GetReadAdapter().RawQueryRow(sqlStr)
	}

	if err != nil {
		return
	}

	row.Scan(dest)
}

func (r *Resource) FetchAll(c *Collection) {
	var (
		rows *sql.Rows
		err  error
	)

	sql := c.GetSelect().Assemble()

	transaction := c.GetTransaction()
	if transaction != nil {
		rows, err = transaction.RawQuery(sql)
	} else {
		rows, err = r.GetReadAdapter().RawQuery(sql)
	}

	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		item := NewItem(c.GetResourceName(), r.GetIdName())
		c.resource._fetch(rows, item)
		c.AddItem(item)
	}
}

func (r *Resource) FetchRow(item *Item) {
	var (
		rows *sql.Rows
		err  error
	)

	sql := NewSelect(r.GetReadAdapter())
	self := r.getSelfData(item)
	read := r.GetReadAdapter()

	sql.From(r.GetMainTable(), "*", "")

	for key, value := range self {
		field := r.GetReadAdapter().QuoteIdentifier(fmt.Sprintf("%s.%s", r.GetMainTable(), key))
		sql.Where(fmt.Sprintf("%s=?", field), value)
	}

	transaction := item.GetTransaction()
	if transaction != nil {
		rows, err = transaction.RawQuery(sql.Assemble())
	} else {
		rows, err = read.RawQuery(sql.Assemble())
	}

	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		r._fetch(rows, item)
		return
	}
}

func (r *Resource) _fetch(rows *sql.Rows, item *Item) {
	cols, _ := rows.Columns()
	contianers := make([]interface{}, len(cols))
	for i := 0; i < len(cols); i++ {
		var contianer interface{}
		contianers[i] = &contianer
	}
	rows.Scan(contianers...)
	item.SetRaw(contianers)
	for i := 0; i < len(cols); i++ {
		item.SetData(cols[i], contianers[i])
	}
}

func (r *Resource) _getLoadSelect(field string, value interface{}) *Select {
	field = r.GetReadAdapter().QuoteIdentifier(fmt.Sprintf("%s.%s", r.GetMainTable(), field))
	sql := NewSelect(r.GetReadAdapter())
	sql.From(r.GetMainTable(), "*", "")
	sql.Where(fmt.Sprintf("%s=?", field), value)
	return sql
}

func (r *Resource) getSelfData(item *Item) map[string]interface{} {
	newMap := make(map[string]interface{})
	for key, val := range item.GetMap() {
		if r.IsExistField(key) {
			newMap[key] = val
		}
	}

	return newMap
}

func (r *Resource) Save(item *Item) error {
	var err error

	newMap := r.getSelfData(item)

	if item.GetId() > 0 {
		condition := r.GetReadAdapter().QuoteInto(fmt.Sprintf("%s=?", r.GetIdName()), item.GetId())

		transaction := item.GetTransaction()
		if transaction != nil {
			_, err = transaction.Update(r.GetMainTable(), newMap, condition)
		} else {
			_, err = r.GetWriteAdapter().Update(r.GetMainTable(), newMap, condition)
		}

	} else {
		var lastId int64

		transaction := item.GetTransaction()
		if transaction != nil {
			lastId, err = transaction.Insert(r.GetMainTable(), newMap)
		} else {
			lastId, err = r.GetWriteAdapter().Insert(r.GetMainTable(), newMap)
		}

		item.SetId(lastId)
	}
	return err
}

func (r *Resource) Delete(item *Item) error {
	var err error

	condition := r.GetReadAdapter().QuoteInto(fmt.Sprintf("%s=?", r.GetIdName()), item.GetId())

	transation := item.GetTransaction()
	if transation != nil {
		_, err = transation.Delete(r.GetMainTable(), condition)
	} else {
		_, err = r.GetWriteAdapter().Delete(r.GetMainTable(), condition)
	}

	return err
}

func (r *Resource) GetReadAdapter() adapter.Adapter {
	return F.GetConnect("read")
}

func (r *Resource) GetWriteAdapter() adapter.Adapter {

	write := F.GetConnect("write")
	if write == nil {
		write = r.GetReadAdapter()
	}

	return write
}

func (r *Resource) CreateTable(tb *table.Table) (sql.Result, error) {
	sqlFragment := strings.Join(r.getColumnsDefinition(tb), ",\n") +
		strings.Join(r.getIndexsDefinition(tb), ",\n") +
		strings.Join(r.getForeignKeysDefinition(tb), ",\n")

	sql := fmt.Sprintf("CREATE TABLE %s (\n%s\n) %s",
		tb.GetName(),
		sqlFragment,
		"tableOptions")
	return r.GetWriteAdapter().RawExec(sql)
}

func (r *Resource) getIndexsDefinition(tb *table.Table) []string {
	var definition []string

	indexs := tb.GetIndexs()
	if len(indexs) > 0 {
		for _, index := range indexs {
			indexType := "KEY"
			switch index.GetType() {
			case table.INDEX_TYPE_PRIMARY:
				indexType = "PRIMARY KEY"
			default:
				indexType = strings.ToUpper(index.GetType())
			}
			columnJoin := strings.Join(index.GetColumns(), ", ")

			definition = append(definition, fmt.Sprintf(
				"  %s %s (%s)",
				indexType,
				index.GetName(),
				columnJoin,
			))
		}
	}
	return definition
}

func (r *Resource) getForeignKeysDefinition(tb *table.Table) []string {
	var definition []string

	relations := tb.GetForeginKeys()
	if len(relations) > 0 {
		for _, fk := range relations {
			definition = append(definition, fmt.Sprintf(
				"  CONSTRAINT %s FOREIGN KEY (%s) REFERENCES %s (%s) ON DELETE %s ON UPDATE %s",
				fk.GetName(),
				fk.GetColumnName(),
				fk.GetRefTableName(),
				fk.GetRefColumnName(),
				fk.GetOnDelete(),
				fk.GetOnUpdate(),
			))
		}
	}
	return definition
}

func (r *Resource) getColumnsDefinition(tb *table.Table) []string {

	columns := tb.GetColumns()
	primaryKeys := make([]string, len(columns))
	if len(columns) == 0 {
		panic("table columns are not defined")
	}

	var definition []string
	for _, column := range columns {
		definition = append(definition,
			fmt.Sprintf(" %s %s",
				"column.GetName()",
				r.getColumnDefinition(column)),
		)
		// primaryKeys[column.GetPrimaryPostion()]
	}

	keyString := ""
	for _, v := range primaryKeys {
		if v != "" {
			if keyString == "" {
				keyString += v
			} else {
				keyString += ", " + v
			}

		}
	}
	if keyString != "" {
		definition = append(definition, fmt.Sprintf("  PRIMARY KEY (%s)", keyString))
	}
	return definition
}

func (r *Resource) getColumnDefinition(c *table.Column) string {
	return fmt.Sprintf("%s%s%s%s%s COMMENT %s",
		c.GetType(),
		c.GetUnsigned(),
		c.GetNullAble(),
		c.GetDefault(),
		c.GetIdentity(),
		c.GetComment(),
	)
}
