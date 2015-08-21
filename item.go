package db

import (
	"github.com/liuzhiyi/go-db/data"
)

type eventFunc func(*Item)

type Item struct {
	data.Item
	resource  *Resource
	events    map[string][]eventFunc
	tableName string
	idField   string
}

func NewItem(tableName string, idField string) *Item {
	i := new(Item)
	i.Init(tableName, idField)
	return i
}

func (i *Item) Init(tableName string, idField string) {
	i.tableName = tableName
	i.idField = idField
	F.SetResourceSingleton(tableName, idField)
	i.Item.Init()
}

func (i *Item) GetResourceName() string {
	return i.tableName
}

func (i *Item) GetResource() *Resource {
	if i.resource == nil {
		i.resource = F.GetResourceSingleton(i.GetResourceName())
	}
	return i.resource
}

func (i *Item) GetIdName() string {
	return i.GetResource().GetIdName()
}

func (i *Item) GetId() int {
	idName := i.GetIdName()
	if idName == "" {
		idName = "id"
	}
	return i.GetInt(idName)
}

func (i *Item) SetId(id int64) {
	if i.GetIdName() != "" {
		i.SetData(i.GetIdName(), id)
	} else {
		i.SetData("id", id)
	}
}

func (i *Item) Load(id int) {
	for _, f := range i.events["beforeLoad"] {
		f(i)
	}
	i.GetResource().Load(i, id)
	for _, f := range i.events["afterLoad"] {
		f(i)
	}
}

func (i *Item) Delete() error {
	if err := i.GetResource().Delete(i); err != nil {
		return err
	} else {
		return nil
	}

}

func (i *Item) Save() error {

	if err := i.GetResource().Save(i); err != nil {
		return err
	} else {
		return nil
	}
}

func (i *Item) GetCollection() *Collection {
	return F.GetCollectionObject(i.GetResourceName())
}
