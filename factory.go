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
	connect          map[string]*adapter.Adapter
	isInitDb         bool
}

func NewFactory() *Factory {
	f := new(Factory)
	f.itemObject = make(map[string]*Item)
	f.collectionObject = make(map[string]*Collection)
	f.resourceObject = make(map[string]*Resource)
	f.connect = make(map[string]*adapter.Adapter)
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

func (f *Factory) getConnect(name string) *adapter.Adapter {
	if !f.isInitDb {
		panic("you haven't initialize db connect")
	}
	return f.connect[name]
}

func (f *Factory) GetItemObject(table string) *Item {
	return NewItem(table)
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
	return NewCollection(table)
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

func init() {
	F = NewFactory()
}
