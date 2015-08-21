package db

import (
	"fmt"

	"github.com/liuzhiyi/go-db/adapter"
	"github.com/liuzhiyi/go-db/data"
)

type eventFunc func(*Item)

type Item struct {
	data.Item
	transaction *adapter.Transaction
	resource    *Resource
	events      map[string][]eventFunc
	tableName   string
	idField     string
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
	if transaction := i.GetTransaction(); transaction != nil {
		transaction.Begin()
		defer transaction.Commit()
	}

	if err := i.GetResource().Delete(i); err != nil {
		return err
	} else {
		return nil
	}

}

func (i *Item) Save() error {
	if transaction := i.GetTransaction(); transaction != nil {
		transaction.Begin()
		defer transaction.Commit()
	}

	if err := i.GetResource().Save(i); err != nil {
		return err
	} else {
		return nil
	}
}

func (i *Item) GetCollection() *Collection {
	return F.GetCollectionObject(i.GetResourceName())
}

func (i *Item) SetTransaction(t *adapter.Transaction) error {
	if i.transaction != nil && !i.transaction.IsOver() {
		return fmt.Errorf("current transaction haven't overed")
	}

	i.transaction = t

	return nil
}

func (i *Item) GetTransaction() *adapter.Transaction {
	if i.transaction != nil {
		if !i.transaction.IsOver() {
			return i.transaction
		}
	}

	return nil
}
