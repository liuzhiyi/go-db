package db

import (
	"sync"

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
	collectionLock   sync.RWMutex
	itemLock         sync.RWMutex
	resourceLock     sync.RWMutex
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

	f.itemLock.RLock()
	_, exists := f.itemObject[table]
	if !exists {
		f.itemLock.RUnlock()
		f.itemLock.Lock()
		f.itemObject[table] = item
		f.itemLock.Unlock()
	} else {
		f.itemLock.RUnlock()
	}
}

func (f *Factory) GetItemSingleton(table string) *Item {
	f.itemLock.RLock()
	defer f.itemLock.RUnlock()
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

	f.collectionLock.RLock()
	_, exists := f.collectionObject[table]
	if !exists {
		f.collectionLock.RUnlock()
		f.collectionLock.Lock()
		f.collectionObject[table] = collection
		f.collectionLock.Unlock()
	} else {
		f.collectionLock.RUnlock()
	}
}

func (f *Factory) GetCollectionSingleton(table string) *Collection {
	f.collectionLock.RLock()
	defer f.collectionLock.Unlock()
	if c, ok := f.collectionObject[table]; ok {
		return c
	} else {
		return nil
	}
}

func (f *Factory) GetResourceSingleton(table string, idField string) *Resource {
	f.resourceLock.RLock()
	r, ok := f.resourceObject[table]
	f.resourceLock.RUnlock()
	if ok {
		return r
	}

	return f.SetResourceSingleton(table, idField)

}

func (f *Factory) SetResourceSingleton(table, idField string) *Resource {
	var r *Resource

	f.resourceLock.RLock()
	if _, ok := f.resourceObject[table]; !ok {
		f.resourceLock.RUnlock()
		f.resourceLock.Lock()
		defer f.resourceLock.Unlock()
		r = NewResource(table, idField)
		f.resourceObject[table] = r
	} else {
		f.resourceLock.RUnlock()
	}

	return r
}

func (f *Factory) Destroy() {
	for _, conn := range f.connect {
		conn.Close()
	}
}

func init() {
	F = NewFactory()
}
