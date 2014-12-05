package db

import (
	"github.com/liuzhiyi/go-db/data"
)

type eventFunc func(*Item)

type Item struct {
	data.Item
	resource interface{}
    events map[string][]eventFunc
}

func (i *Item) GetResource() {

}

func (i *Item) GetIdName() string {
    return i.GetResource()->GetIdName()
}

func (i *Item) GetId() {
    idName = i.GetIdName()
    if idName != "" {
        return i.GetData(idName)
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
    if _, err := i.GetResource().Delete(i); err != nil {
        i.GetResource().RollBack()
        return err
    } else {
        i.GetResource().Commit()
        return nil
    }

}

func (i *Item) Save() error {
    i.GetResource().BeginTransaction()
    if _, err := i.GetResource().Delete(i); err != nil {
        i.GetResource().RollBack()
        return err
    } else {
        i.GetResource().Commit()
        return nil
    }
}

func (i *Item) GetCollection() *Collection {

}

