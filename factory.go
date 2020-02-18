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
	f.itemLock.RUnlock()
	if !exists {
		f.itemLock.Lock()
		f.itemObject[table] = item
		f.itemLock.Unlock()
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
	f.collectionLock.RUnlock()
	if !exists {
		f.collectionLock.Lock()
		f.collectionObject[table] = collection
		f.collectionLock.Unlock()
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
	return nil
}

func (f *Factory) SetResourceSingleton(table, idField string) (*Resource,error) {
	f.resourceLock.RLock()
	r, ok := f.resourceObject[table]
	f.resourceLock.RUnlock()
	if !ok {
		r, err := NewResource(table, idField)
		if err != nil {
			return nil, err
		}

		f.resourceLock.Lock()
		defer f.resourceLock.Unlock()
		f.resourceObject[table] = r
		return r, nil
	} else {
		return r, nil
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
