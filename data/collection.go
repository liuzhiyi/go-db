package data

type Collection struct {
	items []Item
}

func (c *Collection) GetItemByColumnValue(column string, value interface{}) {
	for i := 0; i < len(c.items); i++ {
		if c.items[i].GetData(column) == value {

		}
	}
}
