package db

import (
	"github.com/liuzhiyi/go-db/data"
)

type Collection struct {
	data.Collection
	s Select
}

func (c *Collection) Load() {
	if c.IsLoaded() {
		return
	}
	c._beforeLoad()
	c._renderFilters()
	c._renderOrders()
	c._renderLimit()
	data := c.GetData()
	c.ResetData()
	for _, row := range data {
		item := NewItem()
		item.AddData(row)
		c.AddItem(item)
	}
	c._setIsLoaded()
	c._afterLoad()
}

func (c *Collection) GetData() {
	if c.GetItems() == nil {
		c._beforeLoad()
		c._renderFilters()
		c._renderOrders()
		c._renderLimit()
		c._fetchAll(c.s)
	}
}

func (c *Collection) ResetData() {
	c._data = nil
}

func (c *Collection) _beforeLoad() {

}

func (c *Collection) _afterLoad() {

}

func (c *Collection) _reset() {
	c.s.Reset()
	c._initSelect()
	c._setIsloaded(false)
	c._data = nil
}

func (c *Collection) _fetchAll() {

}
