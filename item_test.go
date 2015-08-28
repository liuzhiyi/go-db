package db

import (
	"fmt"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestItem(t *testing.T) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?strict=false", "root", "", "127.0.0.1:3306", "weishop")
	F.InitDb("mysql", dsn, "")
	i := NewItem("tp_api", "id")
	i.SetId(1)
	i.SetData("uid", 1)
	i.SetData("token", "487524277")
	i.Delete()
	fmt.Println(i.GetMap())
}
