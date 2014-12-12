package adapter

import (
	"fmt"
	"strconv"
)

type pgAdapter struct {
	Adapter
}

func (p *pgAdapter) Limit(sql string, count, offset int64) string {
	if count <= 0 {
		panic(fmt.Sprintf("LIMIT argument count=%s is not valid", count))
	}
	if offset < 0 {
		panic(fmt.Sprintf("LIMIT argument offset=%s is not valid", offset))
	}
	sql += " LIMIT " + strconv.FormatInt(count, 10)
	if offset > 0 {
		sql += " OFFSET " + strconv.FormatInt(offset, 10)
	}
	return sql
}
