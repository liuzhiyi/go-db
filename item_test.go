package db

import (
	"fmt"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/liuzhiyi/go-db/adapter"
)

func TestItem(t *testing.T) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?strict=false", "root", "", "127.0.0.1:3306", "weishop")
	a := adapter.NewAdapter("mysql", dsn)
	r := NewResource(a)
	r.idName = "id"
	r.mainTable = "tp_api"
	i := new(Item)
	i.Init()
	i.resource = &r
	i.SetId(1)
	i.SetData("uid", 1)
	i.SetData("token", "487524277")
	i.Delete()
	fmt.Println(i.GetMap())
}
