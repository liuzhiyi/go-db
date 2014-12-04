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

/**
*
*建议一般情况下开启事物机制
*****/
func (i *Item) Delete() {
    defer func() {
        e := recover()
    }()
    i.GetResource().BeginTransaction()
    i.GetResource().

}

func (i *Item) GetCollection() *Collection {

}

