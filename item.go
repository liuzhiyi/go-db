package db

import (
	"github.com/liuzhiyi/go-db/data"
)

type eventFunc func(*Item)

type Item struct {
	data.Item
	resource *Resource
	events   map[string][]eventFunc
}

func NewItem() *Item {
	i := new(Item)
	i.Init()
	return i
}

func (i *Item) Init() {
	i.Item.Init()
}

func (i *Item) GetResource() *Resource {
	return i.resource
}

func (i *Item) GetIdName() string {
	return i.GetResource().GetIdName()
}

func (i *Item) GetId() int {
	idName := i.GetIdName()
	if idName != "" {
		//
	}
	return i.GetInt(idName)
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
	i.GetResource().BeginTransaction()
	if err := i.GetResource().Delete(i); err != nil {
		i.GetResource().RollBack()
		return err
	} else {
		i.GetResource().Commit()
		return nil
	}

}

func (i *Item) Save() error {
	i.GetResource().BeginTransaction()
	if err := i.GetResource().Delete(i); err != nil {
		i.GetResource().RollBack()
		return err
	} else {
		i.GetResource().Commit()
		return nil
	}
}

func (i *Item) GetCollection() *Collection {
	return new(Collection)
}
