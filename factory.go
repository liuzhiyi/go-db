package db

import (
	"github.com/liuzhiyi/go-db/adapter"
)

var F *Factory

/**
*key is the table name, registry the singleton.
 */
type Factory struct {
	itemObject       map[string]*Item
	collectionObject map[string]*Collection
	resourceObject   map[string]*Resource
	connect          map[string]adapter.Adapter
	isInitDb         bool
}

func NewFactory() *Factory {
	f := new(Factory)
	f.itemObject = make(map[string]*Item)
	f.collectionObject = make(map[string]*Collection)
	f.resourceObject = make(map[string]*Resource)
	f.connect = make(map[string]adapter.Adapter)
	return f
}

func (f *Factory) IsInitDb() bool {
	return f.isInitDb
}

func (f *Factory) InitDb(driverName, readDsn, writeDsn string) {
	if readDsn == "" {
		panic("read connect must initialize!")
	}
	f.connect["read"] = adapter.NewAdapter(driverName, readDsn)
	if writeDsn != "" {
		f.connect["write"] = adapter.NewAdapter(driverName, writeDsn)
	}
	f.isInitDb = true
}

func (f *Factory) GetConnect(name string) adapter.Adapter {
	if !f.isInitDb {
		panic("you haven't initialize db connect")
	}
	return f.connect[name]
}

func (f *Factory) RegisterItem(item *Item) {
	table := item.GetResourceName()
	if table == "" {
		return
	}

	_, exists := f.itemObject[table]
	if !exists {
		f.itemObject[table] = item
	}
}

func (f *Factory) GetItemSingleton(table string) *Item {
	if c, ok := f.itemObject[table]; ok {
		return c
	} else {
		return nil
	}
}

func (f *Factory) RegisterCollection(collection *Collection) {
	table := collection.GetResourceName()
	if table == "" {
		return
	}

	_, exists := f.collectionObject[table]
	if !exists {
		f.collectionObject[table] = collection
	}
}

func (f *Factory) GetCollectionSingleton(table string) *Collection {
	if c, ok := f.collectionObject[table]; ok {
		return c
	} else {
		return nil
	}
}

func (f *Factory) GetResourceSingleton(table string, idField string) *Resource {
	if r, ok := f.resourceObject[table]; ok {
		return r
	} else {
		f.resourceObject[table] = NewResource(table, idField)
		return f.resourceObject[table]
	}
}

func (f *Factory) SetResourceSingleton(table, idField string) {
	if _, ok := f.resourceObject[table]; !ok {
		f.resourceObject[table] = NewResource(table, idField)
	}
}

func (f *Factory) Destroy() {
	for _, conn := range f.connect {
		conn.Close()
	}
}

func init() {
	F = NewFactory()
}
