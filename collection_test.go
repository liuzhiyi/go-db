package db

import (
	"fmt"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/liuzhiyi/go-db/adapter"
)

func TestCollection(t *testing.T) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?strict=false", "root", "", "127.0.0.1:3306", "xing100")
	a := adapter.NewAdapter("mysql", dsn)
	r := NewResource(a)
	c := Collection{}
	c._init()
	c.curPage = 3
	c.SetMainTable("xing100b2c_users")
	c._initSelect()
	c.resource = r
	c.Join("xing100b2c_order_info as o", "m.user_id = o.user_id", "consignee")
	c.AddFieldToSelect("user_name as w, sex, o.user_id", c.GetMainAlias())
	// c.AddFieldToFilter("user_name", "eq", "liu")
	// c.AddFieldToNewFilter("user_name", "or eq", "13337311235")
	// c.AddFieldToFilter("user_name", "or eq", "13875175665")
	// c.AddFieldToFilter("m.email", "or eq", "123@d.com")
	c.Load()
	fmt.Println(c.GetSelect().Assemble())
	for _, item := range c.GetItems() {
		fmt.Println(item.(*Item).GetString("consignee"))
	}
	a.Close()
}
