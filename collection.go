package db

import (
	"fmt"
	"math"
	"regexp"
	"strings"

	"github.com/liuzhiyi/go-db/data"
)

type Collection struct {
	data.Collection
	resource     *Resource
	s            *Select
    lastSql   string
	orders       []string
	filter       Filter
	resourceName string
	whereFlag    bool
	isLoaded     bool
	isAllFields  bool //is query the main table all fields, default is true
	pageSize     int64
	totalSize    int64
	curPage      int64
}

func NewCollection(resourceName string) *Collection {
	c := new(Collection)
	c.Init(resourceName)
	return c
}

func (c *Collection) Init(resourceName string) {
	c.resourceName = resourceName
	c.totalSize = -1
	c.pageSize = 0
	c.curPage = 1
	c.whereFlag = false
	c.isLoaded = false
	c.isAllFields = true
	c._initSelect()
}

func (c *Collection) GetResource() *Resource {
	if c.resource == nil {
		c.resource = F.GetResourceSingleton(c.GetResourceName())
	}
	return c.resource
}

func (c *Collection) GetResourceName() string {
	return c.resourceName
}

func (c *Collection) GetMainTable() string {
	return c.GetResource().GetMainTable()
}

func (c *Collection) GetMainAlias() string {
	return "m"
}

func (c *Collection) SetMainTable(table string) {
	c.GetResource().mainTable = table
}

func (c *Collection) Load() {
	c._beforeLoad()

	if c.IsAllFields() {
		c.AddFieldToSelect("*", c.GetMainAlias())
	}
	c._where()
	c._renderOrders()
	c._renderLimit()
    if !c._isChanged() && c.IsLoaded() {
        return
    }

	c._fetchAll()
	c.ResetData()

    c.lastSql = c.GetSelect().Assemble()
	c._setIsLoaded(true)
	c._afterLoad()
}

func (c *Collection) _isChanged() bool {
    sql := c.GetSelect().Assemble()
    return sql != c.lastSql
}

func (c *Collection) _prepareSelect() {

}

func (c *Collection) GetItems() []*Item {
	var sets []*Item
	for _, item := range c.Collection.GetItems() {
		sets = append(sets, item.(*Item))
	}
	return sets
}

func (c *Collection) ResetData() {

}

func (c *Collection) GetSelect() *Select {
	if c.s == nil {
		c.s = NewSelect(c.GetResource().GetReadAdapter())
	}
	return c.s
}

func (c *Collection) _initSelect() {
	c.GetSelect().From(fmt.Sprintf("%s as %s", c.GetMainTable(), c.GetMainAlias()), "", "")
}

func (c *Collection) _initSelectFields() {

}

/**
*@param fields:"col1, col2 as c, o.col3 ..."
*@param correlation can equal ""
*the @parmam fields'correlation name > @param correlation
**/
func (c *Collection) AddFieldToSelect(fields string, correlation string) *Collection {
	c.GetSelect().Columns(fields, correlation)
	return c
}

func (c *Collection) IsAllFields() bool {
	if c.isAllFields {
		for _, column := range c.GetSelect().GetColumnPart() {
			if column[0] == c.GetMainAlias() {
				c.isAllFields = false
				break
			}
		}
	}
	return c.isAllFields
}

func (c *Collection) AddFieldToFilter(field, key string, value interface{}) *Collection {
	if len(c.filter) > 0 {
		c.filter.SetCondition(field, key, value)
	} else {
		f := NewFilter()
		f.SetCondition(field, key, value)
		c.filter = f
	}
	return c
}

func (c *Collection) AddFieldToNewFilter(field, key string, value interface{}) *Collection {
	c._where()
	if m := c._splitKey(key); len(m) == 2 {
		c.whereFlag = true
		key = m[1]
	}
	c.AddFieldToFilter(field, key, value)
	return c
}

func (c *Collection) _splitKey(key string) []string {
	reg := regexp.MustCompile(`^[oO][rR]\s+(.+)$`)
	return reg.FindStringSubmatch(key)
}

func (c *Collection) _renderFilter() string {
	result := "("
	for field, condition := range c.filter {
		for i := 0; i < len(condition); i++ {
			for key, value := range condition[i] {
				if result != "(" {
					if m := c._splitKey(key); len(m) == 2 {
						key = m[1]
						result += " " + SQL_OR + " "
					} else {
						result += " " + SQL_AND + " "
					}
				}
				result += c._getConditionSql(field, key, value)
			}
		}
	}
	result += ")"
	c.filter = NewFilter()
	return result
}

func (c *Collection) _where() {
	if len(c.filter) > 0 {
		if c.whereFlag {
			c.GetSelect().OrWhere(c._renderFilter(), nil)
		} else {
			c.GetSelect().Where(c._renderFilter(), nil)
		}
	}
}

func (c *Collection) _getConditionSql(fieldName, key string, value interface{}) string {
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
	if expre, ok := conditionKeyMap[key]; ok {
		query = c._prepareQuotedSqlCondition(expre, value, fieldName)
	}

	return query
}

func (c *Collection) _prepareQuotedSqlCondition(text string, value interface{}, fieldName string) string {
	sql := c.resource.GetReadAdapter().QuoteInto(text, value)
	sql = strings.Replace(sql, "{{fieldName}}", fieldName, -1)
	return sql
}

/**
*@param spec : "field direction"
**/
func (c *Collection) AddOrder(spec string) {
	c.orders = append(c.orders, spec)
}

func (c *Collection) Join(table, cond, cols string) *Collection {
	if cols == "" {
		cols = "*"
	}
	c.GetSelect().Join(table, cond, cols, "")
	return c
}

func (c *Collection) JoinLeft(table, cond, cols string) *Collection {
	if cols == "" {
		cols = "*"
	}
	c.GetSelect().JoinLeft(table, cond, cols, "")
	return c
}

func (c *Collection) _renderOrders() {
	c.GetSelect().Order(c.orders...)
}

func (c *Collection) _renderLimit() {
	if c.pageSize > 0 {
		c.GetSelect().LimitPage(c.GetCurPage(0), c.pageSize)
	}
}

func (c *Collection) _beforeLoad() {

}

func (c *Collection) _afterLoad() {

}

func (c *Collection) GetCurPage(offset int64) int64 {
	if c.curPage+offset <= 0 {
		return 1
	} else if c.curPage+offset > c.GetLastPage() {
		return c.GetLastPage()
	} else {
		return c.curPage + offset
	}
}

func (c *Collection) GetLastPage() int64 {
	count := c.GetSize()
	if count <= 0 {
		return 1
	} else {
		return int64(math.Ceil(float64(count / c.pageSize)))
	}
}

func (c *Collection) GetSize() int64 {
	if c.totalSize < 0 {
		sql := c.GetCountSql()
		fmt.Println(sql)
		c.resource.FetchOne(sql, &c.totalSize)
	}
	return c.totalSize
}

func (c *Collection) SetPageSize(size int64) {
    if size > 0 {
        c.pageSize = size
    }
}

func (c *Collection) SetCurPage(page int64) {
    if page < 0 {
        c.curPage = 1
    } else if page > c.GetLastPage() {
        c.curPage = c.GetLastPage()
    } else {
        c.curPage = page
    }
}

func (c *Collection) GetCountSql() string {
	c._where()

	return c.GetSelect().GetCountSql()
}

func (c *Collection) Save() error {
	var err error
	for _, item := range c.GetItems() {
		err = item.Save()
		if err != nil {
			break
		}
	}
	return err
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

func (c *Collection) Each(f func(*Item)) {
	for _, item := range c.GetItems() {
		f(item)
	}
}
