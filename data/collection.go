package data

type Collection struct {
	pageSize    int
	currentPage int
	items       []interface{}
}

func (c *Collection) GetItemByColumnValue(column string, value interface{}) {
	// for i := 0; i < len(c.items); i++ {
	// 	if c.items[i].GetData(column) == value {

	// 	}
	// }
}

func (c *Collection) AddItem(i interface{}) {
	c.items = append(c.items, i)
}

func (c *Collection) GetItems() []interface{} {
	return c.items
}

func (c *Collection) Count() int {
	return len(c.items)
}

func (c *Collection) Each() {
}
