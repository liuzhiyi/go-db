package adapter

import (
	"fmt"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestAdapter(t *testing.T) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?strict=false", "root", "", "127.0.0.1:3306", "magento")
	a := NewAdapter("mysql", dsn)
	rows := a.Query("select role_id from admin_role")
	defer rows.Close()
	fmt.Println(rows.Columns())
	c := make([]interface{}, 1)
	var tmp interface{}
	c[0] = &tmp
	if rows.Next() {
		fmt.Println(rows.Scan(c...), c)
	}
}
