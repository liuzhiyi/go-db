package main

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	db "github.com/liuzhiyi/go-db"
)

func init() {
	db.F.InitDb("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?strict=false", "root", "", "127.0.0.1:3306", "weishop"), "")
	db.F.GetResourceSingleton("tp_api", "id")
}

type Api struct {
	db.Item
}

func NewApi() *Api {
	table := "tp_api"
	a := new(Api)
	a.Init(table)
	return a
}

func main() {
	api := NewApi()
	api.Load(1)
	fmt.Println(api.GetMap())
	fmt.Println("token =>", api.GetString("token"))
}
