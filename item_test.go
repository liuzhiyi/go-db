package db

import (
	"fmt"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestItem(t *testing.T) {
	i := new(Item)
	i.Load(3)
	fmt.Println(i)
}
