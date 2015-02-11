go-db
====

它是一种数据库数据模型，目的是提高开发效率，代码重用等。
        item类----------->对应一张表的实体，也对应一个集合中的一条记录。
        collection类----->对应连表查询，或者对应多条记录的集合。
        select类--------->查询语句的构造类。
        resource类------->数据库资源。
        adapter类-------->统一处理各种数据库各种特性。
        fiter类---------->过滤器，有待实现。
        hook类----------->钩子，有待实现
        factory---------->统一分配资源。

# 安装/更新

```
go get -u https://github.com/liuzhiyi/go-db
```

# 使用


```go
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
```
