package db

import (
	"fmt"
	"regexp"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestSelect(t *testing.T) {
	s := new(Select)
	s._init()
	s.From("table1 as 1", "*", "db1")
	s.Join("table2 as 2", "1.id = 2.id", "*", "schema")
	s.Where("2.name=?", "liming")
	s.OrWhere("1.name", "xiao")
	s.Group("1.id")
	s.Group("2.id")
	s.Order("1.id desc")
	fmt.Println(s.Assemble())
	reg := regexp.MustCompile(`(.*\W)((?i:asc|desc))\b`)
	m := reg.FindStringSubmatch("1.id desC")
	fmt.Println(len(m), m)
}
