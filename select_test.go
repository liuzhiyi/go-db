package db

import (
	"fmt"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/liuzhiyi/go-db/adapter"
)

func TestSelect(t *testing.T) {
	a := adapter.Mysql{}
	s := NewSelect(a)
	s.From("table1 as 1", "*", "db1")
	s.Join("table2 as 2", "1.id = 2.id", "*", "schema")
	s.Where("2.name=?", "liming")
	s.Columns("select from user where id = 1", "")
	s.OrWhere("1.name", "xiao")
	s.Group("1.id")
	s.Group("2.id")
	s.Order("1.id desc")
	s.Having("2.name = ?", "wudao")
	s.Limit(1, 30)
	fmt.Println(s.Assemble())
}
