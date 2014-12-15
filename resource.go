package db

import (
	"fmt"

	"github.com/liuzhiyi/go-db/adapter"
)

type Resource struct {
	idName    string
	mainTable string
	adapter   *adapter.Adapter
}

func NewResource(a *adapter.Adapter) Resource {
	return Resource{
		adapter: a,
	}
}

func (r *Resource) GetIdName() string {
	return r.idName
}

func (r *Resource) BeginTransaction() {
	r.adapter.BeginTransaction()
}

func (r *Resource) Commit() {
	r.adapter.Commit()
}

func (r *Resource) RollBack() {
	r.adapter.RollBack()
}

func (r *Resource) GetMainTable() string {
	return r.GetTable(r.mainTable)
}

func (r *Resource) GetTable(name string) string {
	return r.adapter.GetTableName(name)
}

func (r *Resource) Load(item *Item, id int) {
	read := r.GetReadAdapter()
	field := r.GetIdName()
	sql := r._getLoadSelect(field, id)
	fmt.Println(sql.Assemble())
	rows := read.Query(sql.Assemble())
	defer rows.Close()
	cols, _ := rows.Columns()
	for rows.Next() {
		contianers := make([]interface{}, len(cols))
		for i := 0; i < len(cols); i++ {
			var contianer interface{}
			contianers[i] = &contianer
		}
		rows.Scan(contianers...)
		for i := 0; i < len(cols); i++ {
			item.SetData(cols[i], contianers[i])
		}
		return
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
	return r.adapter
}
