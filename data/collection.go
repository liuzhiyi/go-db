package data

type Collection struct {
	pageSize    int
	currentPage int
	items       []*Item
}

func (c *Collection) GetItemByColumnValue(column string, value interface{}) {
	for i := 0; i < len(c.items); i++ {
		if c.items[i].GetData(column) == value {

		}
	}
}

func (c *Collection) AddItem(i *Item) {
	c.items = append(c.items, i)
}

func (c *Collection) GetItems() []*Item {
	return c.items
}

func (c *Collection) Count() int {
	return len(c.items)
}
