package db

import (
	"fmt"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestCollection(t *testing.T) {
	//dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?strict=false", "root", "", "127.0.0.1:3306", "magento")
	//a := adapter.NewAdapter("mysql", dsn)
	c := Collection{}
	c._initSelect()
	c.AddFieldToSelect("name as w", "n")
	f := NewFilter()
	f.SetCondition("id", "eq", "4")
	c.AddFieldToFilter([]string{"id"}, f)
	fmt.Println(c.GetSelect().Assemble())
}
