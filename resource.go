package db

import (
	"database/sql"
	"fmt"
	"sort"

	"github.com/liuzhiyi/go-db/adapter"
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
	rows, err := r.GetReadAdapter().Query(sql)
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
	read := r.GetReadAdapter()
	field := r.GetIdName()
	sql := r._getLoadSelect(field, id)
	rows, err := read.Query(sql.Assemble())
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		r._fetch(rows, item)
		return
	}
}

func (r *Resource) FetchOne(sql string, dest interface{}) {
	row, err := r.GetReadAdapter().QueryRow(sql)
	if err != nil {
		return
	}

	row.Scan(dest)
}

func (r *Resource) FetchAll(c *Collection) {
	sql := c.GetSelect().Assemble()
	rows, err := r.GetReadAdapter().Query(sql)
	if err != nil {
		return
	}

	for rows.Next() {
		item := NewItem(c.GetResourceName(), r.GetIdName())
		c.resource._fetch(rows, item)
		c.AddItem(item)
	}
}

func (r *Resource) FetchRow(item *Item) {

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

func (r *Resource) Save(item *Item) error {
	var err error

	newMap := make(map[string]interface{})
	for key, val := range item.GetMap() {
		if r.IsExistField(key) {
			newMap[key] = val
		}
	}

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
