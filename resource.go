package db

import (
	"database/sql"
	"fmt"

	"github.com/liuzhiyi/go-db/adapter"
)

type Resource struct {
	idName    string
	mainTable string
	read      *adapter.Adapter
	write     *adapter.Adapter
}

func NewResource(a *adapter.Adapter) Resource {
	return Resource{
		read: a,
	}
}

func (r *Resource) GetIdName() string {
	return r.idName
}

func (r *Resource) BeginTransaction() {
	r.GetReadAdapter().BeginTransaction()
}

func (r *Resource) Commit() {
	r.GetReadAdapter().Commit()
}

func (r *Resource) RollBack() {
	r.GetReadAdapter().RollBack()
}

func (r *Resource) GetMainTable() string {
	return r.GetReadAdapter().GetTableName(r.mainTable)
}

func (r *Resource) GetTable(name string) string {
	return r.GetReadAdapter().GetTableName(name)
}

func (r *Resource) Load(item *Item, id int) {
	read := r.GetReadAdapter()
	field := r.GetIdName()
	sql := r._getLoadSelect(field, id)
	fmt.Println(sql.Assemble())
	rows := read.Query(sql.Assemble())
	defer rows.Close()
	for rows.Next() {
		r._fetch(rows, item)
		return
	}
}

func (r *Resource) FetchOne(sql string, dest interface{}) {
	row := r.GetReadAdapter().QueryRow(sql)
	row.Scan(dest)
}

func (r *Resource) FetchAll(c *Collection) {
	sql := c.GetSelect().Assemble()
	rows := r.GetReadAdapter().Query(sql)
	for rows.Next() {
		item := NewItem()
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
	for i := 0; i < len(cols); i++ {
		item.SetData(cols[i], contianers[i])
	}
}

func (r *Resource) _getLoadSelect(field string, value interface{}) *Select {
	field = r.GetReadAdapter().QuoteIdentifier(fmt.Sprintf("%s.%s", r.GetMainTable(), field))
	sql := new(Select)
	sql._init()
	sql.From(r.GetMainTable(), "*", "")
	sql.Where(fmt.Sprintf("%s=?", field), value)
	return sql
}

func (r *Resource) Save(item *Item) error {
	if item.GetInt("id") > 0 {
		condition := r.GetReadAdapter().QuoteInto(fmt.Sprintf("%s=?", r.GetIdName()), item.GetInt("id"))
		fmt.Println(condition)
	}
	return nil
}

func (r *Resource) Delete(item *Item) error {
	return nil
}

func (r *Resource) GetReadAdapter() *adapter.Adapter {
	return r.read
}

func (r *Resource) GetWriteAdapter() *adapter.Adapter {
	if r.write == nil {
		r.write = r.read
	}
	return r.write
}
