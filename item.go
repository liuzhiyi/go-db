package db

import (
	"github.com/liuzhiyi/go-db/data"
)

type eventFunc func(*Item)

type Item struct {
	data.Item
	resource     *Resource
	events       map[string][]eventFunc
	resourceName string
}

func NewItem(resourceName string) *Item {
	i := new(Item)
	i.Init(resourceName)
	return i
}

func (i *Item) Init(resourceName string) {
	i.resourceName = resourceName
	i.Item.Init()
}

func (i *Item) GetResourceName() string {
	return i.resourceName
}

func (i *Item) GetResource() *Resource {
	if i.resource == nil {
		i.resource = F.GetResourceSingleton(i.GetResourceName(), "")
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
	if err := i.GetResource().Save(i); err != nil {
		i.GetResource().RollBack()
		return err
	} else {
		i.GetResource().Commit()
		return nil
	}
}

func (i *Item) GetCollection() *Collection {
	return F.GetCollectionObject(i.GetResourceName())
}
