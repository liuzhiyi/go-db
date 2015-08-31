package main

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	db "github.com/liuzhiyi/go-db"
)

func init() {
	db.F.InitDb("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?strict=false", "root", "", "127.0.0.1:3306", "test"), "")
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

//数据一致性,注意此函数已经重载Item的delete,要调用item的delete，需显示调用。
func (a *Api) Delete() {
	transaction := a.GetResource().BeginTransaction()
	collection := a.GetCollection()
	collection.AddFieldToFilter("website_id", "eq", a.GetInt64("website_id")).Load()
	collection.Each(func(i *db.Item) {
		i.SetTransaction(transaction)
		err := i.Delete()

		if err != nil {
			fmt.Println(err)
		}
	})
	transaction.Commit()
}

type User struct {
	db.Item
}

func NewWebsite() *User {
	u := new(User)
	u.Init("core_website", "website_id")
	return u
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
	api.Delete()
}
