package db

import (
	"github.com/liuzhiyi/go-db/adapter"
)

type Resource struct {
	adapter *adapter.Adapter
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

func (r *Resource) Load(item *Item, id int) {
	read := r.GetReadAdapter()
	field := r._getIdFieldName()
	sql := r._getLoadSelect(field, id)
	rows := read.Query(sql)
	defer rows.Close()
	for rows.Next() {
		cols := rows.Columns()
		contianers := make([]interface{}, len(cols))
		for i := 0; i < len(cols); i++ {
			var contianer interface{}
			contianers[i] = &contianer
		}
		rows.Scan(contianers...)
		for i := 0; i < len(cols); i++ {
			item.SetData(cols[i], *contianers[i])
		}
		return
	}
}

func (r *Resource) _getLoadSelect(field string, value interface{}) *Select {
    field := r.GetReadAdapter().QuoteIdentifier(fmt.Sprintf("%s.%s", r.GetMainTable(), field)
    sql := new(Select)
    sql.From(r.getMainTable(), "*", "")
    sql.Where(fmt.Sprintf("%s=?", field), value)
    return sql
}

func (r *Resource) Save(item *Item) {
    if item.GetInt("id") > 0 {
        condition := r.GetReadAdapter().QuoteInto(fmt.Sprintf("%s=?", r._getIdFieldName()), item.GetInt("id"))
    }
}

func (r *Resource) GetReadAdapter() *adapter.Adapter {
	return r.adapter
}
