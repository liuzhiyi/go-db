package db

import (
	"fmt"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/liuzhiyi/go-db/adapter"
)

func TestItem(t *testing.T) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?strict=false", "root", "", "127.0.0.1:3306", "magento")
	a := adapter.NewAdapter("mysql", dsn)
	r := NewResource(a)
	r.idName = "user_id"
	r.mainTable = "admin_user"
	i := new(Item)
	i.Init()
	i.resource = &r
	i.Load(1)
	fmt.Println(i.GetString("created"))
}
