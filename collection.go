package db

import (
	"strings"

	"github.com/liuzhiyi/go-db/data"
)

type Collection struct {
	data.Collection
	resource  Resource
	s         *Select
	mainTable string
	orders    []string
	filters   []map[string]string
	isLoaded  bool
	pageSize  int
	totalSize int64
}

func (c *Collection) GetMainTable() string {
	return c.mainTable
}

func (c *Collection) SetMainTable(table string) {
	c.mainTable = table
}

func (c *Collection) Load() {
	if c.IsLoaded() {
		return
	}
	c._beforeLoad()
	c._renderFilters()
	c._renderOrders()
	c._renderLimit()
	c._fetchAll()
	c.ResetData()

	c._setIsLoaded(true)
	c._afterLoad()
}

func (c *Collection) _prepareSelect() {

}

func (c *Collection) ResetData() {

}

func (c *Collection) GetSelect() *Select {
	return c.s
}

func (c *Collection) _initSelect() {
	c.GetSelect().From(c.GetMainTable(), "", "")
}

func (c *Collection) _initSelectFields() {

}

func (c *Collection) AddFieldToSelect(field string) {
	if field == "*" {

	}
}

func (c *Collection) AddFieldToFilter(field []string, condition Filter) {
	var conditions []string
	for i := 0; i < len(field); i++ {
		conditions = append(conditions, c._getConditionSql(field[i], condition[field[i]]))
	}
	result := "(" + strings.Join(conditions, ") "+SQL_OR+" (") + ")"
	c.GetSelect().Where(result, nil)
}

func (c *Collection) _getConditionSql(fieldName string, condition map[string]interface{}) string {
	conditionKeyMap := make(map[string]string)
	conditionKeyMap["eq"] = "{{fieldName}} = ?"
	conditionKeyMap["neq"] = "{{fieldName}} != ?"
	conditionKeyMap["like"] = "{{fieldName}} LIKE ?"
	conditionKeyMap["nlike"] = "{{fieldName}} NOT LIKE ?"
	conditionKeyMap["in"] = "{{fieldName}} IN(?)"
	conditionKeyMap["nin"] = "{{fieldName}} NOT IN(?)"
	conditionKeyMap["is"] = "{{fieldName}} IS ?"
	conditionKeyMap["notnull"] = "{{fieldName}} IS NOT NULL"
	conditionKeyMap["null"] = "{{fieldName}} IS NULL"
	conditionKeyMap["gt"] = "{{fieldName}} > ?"
	conditionKeyMap["lt"] = "{{fieldName}} < ?"
	conditionKeyMap["gteq"] = "{{fieldName}} >= ?"
	conditionKeyMap["lteq"] = "{{fieldName}} <= ?"
	conditionKeyMap["finset"] = "FIND_IN_SET(?, {{fieldName}})"
	conditionKeyMap["regexp"] = "{{fieldName}} REGEXP ?"
	conditionKeyMap["from"] = "{{fieldName}} >= ?"
	conditionKeyMap["to"] = "{{fieldName}} <= ?"
	conditionKeyMap["seq"] = "null"
	conditionKeyMap["sneq"] = "null"

	query := ""
	for key, value := range condition {
		if key == "from" || key == "to" {
			if key == "from" {
				query = c._prepareQuotedSqlCondition(conditionKeyMap["from"], value, fieldName)
			}
			if key == "to" {
				if query != "" {
					query += c._prepareQuotedSqlCondition(conditionKeyMap["to"], value, fieldName)
				} else {
					query = c._prepareQuotedSqlCondition(conditionKeyMap["to"], value, fieldName)
				}

			}
		} else if expre, ok := conditionKeyMap[key]; ok {
			query = c._prepareQuotedSqlCondition(expre, value, fieldName)
		}
	}
	return query
}

func (c *Collection) _prepareQuotedSqlCondition(text string, value interface{}, fieldName string) string {
	sql := c.resource.GetReadAdapter().QuoteInto(text, value)
	sql = strings.Replace(sql, "{{fieldName}}", fieldName, -1)
	return sql
}

func (c *Collection) _renderFilters() {
	for _, filter := range c.filters {
		switch filter["type"] {
		case "or":
			c.GetSelect().OrWhere(filter["field"]+"=?", filter["value"])
		case "and":
			c.GetSelect().Where(filter["field"]+"=?", filter["value"])
		}
	}
}

func (c *Collection) _renderOrders() {
	c.GetSelect().Order(c.orders...)
}

func (c *Collection) _renderLimit() {
	if c.pageSize > 0 {
		c.GetSelect().Limit(c.GetCurPage(), c.pageSize)
	}
}

func (c *Collection) _beforeLoad() {

}

func (c *Collection) _afterLoad() {

}

func (c *Collection) GetSize() int64 {
	if c.totalSize <= 0 {
		sql := c.GetCountSql()
		c.resource.FetchOne(sql, c.totalSize)
	}
	return c.totalSize
}

func (c *Collection) GetCountSql() string {
	c._renderFilters()

	countSql := c.GetSelect().Clone()
	countSql.Reset(ORDER)
	countSql.Reset(LIMIT_COUNT)
	countSql.Reset(LIMIT_OFFSET)
	countSql.Reset(COLUMNS)

	return countSql.Assemble()
}

func (c *Collection) Save() {
	for _, item := range c.GetItems() {
		item.Save()
	}
}

func (c *Collection) IsLoaded() bool {
	return c.isLoaded
}

func (c *Collection) _setIsLoaded(flag bool) {
	c.isLoaded = flag
}

func (c *Collection) _reset() {
	c.s.Reset()
	c._initSelect()
	c._setIsLoaded(false)
}

func (c *Collection) _fetchAll() {
	c.resource.FetchAll(c)
}
