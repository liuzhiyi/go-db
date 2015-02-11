package main

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	db "github.com/liuzhiyi/go-db"
)

func init() {
	db.F.InitDb("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?strict=false", "root", "", "127.0.0.1:3306", "weishop"), "")
}

type Api struct {
	db.Item
}

func NewApi() *Api {
	table := "core_api"
	a := new(Api)
	a.Init(table, "api_id")
	return a
}

func main() {
	createData()
	api := NewApi()
	api.Load(1)
	fmt.Println("api_name =>", api.GetString("api_name"))
	// fmt.Println(api.GetDate("time", "2006 01 02 15:04:05"))
	collection := api.GetCollection()
	collection.Join("core_website as w", "m.website_id = w.website_id", "name")
	collection.Load()
	for _, item := range collection.GetItems() {
		fmt.Println(item.GetString("name"))
	}
}
