package db

/**
*key is the table name, registry the singleton.
 */
type Factory struct {
	itemObject       map[string]*Item
	collectionObject map[string]*Collection
	resourceObject   map[string]*Resource
}

func NewFactory() {
	f := new(Factory)
	f.itemObject = make(map[string]*Item)
	f.collectionObject = make(map[string]*Collection)
	f.resourceObject = make(map[string]*Resource)
}

func (f *Factory) GetItemObject(table string) *Item {
	r := f.GetResourceSingleton(table, "")
	return NewItem(r)
}

func (f *Factory) GetItemSingleton(table string) *Item {
	if c, ok := f.itemObject[table]; ok {
		return c
	} else {
		f.itemObject[table] = f.GetItemObject(table)
		return f.itemObject[table]
	}
}

func (f *Factory) GetCollectionObject(table string) *Collection {
	r := f.GetResourceSingleton(table, "")
	return NewCollection(r)
}

func (f *Factory) GetCollectionSingleton(table string) *Collection {
	if c, ok := f.collectionObject[table]; ok {
		return c
	} else {
		f.collectionObject[table] = f.GetCollectionObject(table)
		return f.collectionObject[table]
	}
}

func (f *Factory) GetResourceSingleton(table, idField string) *Resource {
	if r, ok := f.resourceObject[table]; ok {
		return r
	} else {
		f.resourceObject[table] = NewResource(table, idField)
		return f.resourceObject[table]
	}
}
