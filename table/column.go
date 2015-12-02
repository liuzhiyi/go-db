package table

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	UNIQUE = 1 << iota
	ATUO_INCREMENT
	PRIMARY_KEY
	KEY
)

type Column struct {
	name       string
	typ        string
	comment    string
	position   int
	flag       int
	dataType   string
	defaultVal string
	nullAble   bool
	length     int
	scale      int
	precision  int
	unsigned   bool
	primary    bool
	primaryPos int
	identity   bool
}

func (c *Column) SetType(typ string) *Column {
	switch typ {
	// case TYPE_CHAR, TYPE_VARCHAR, TYPE_LONGVARCHAR, TYPE_CLOB:
	// 	typ = TYPE_TEXT
	case TYPE_TINYINT:
		typ = TYPE_SMALLINT
	case TYPE_DOUBLE, TYPE_REAL:
		typ = TYPE_FLOAT
	case TYPE_TIME:
		typ = TYPE_TIMESTAMP
	case TYPE_BINARY, TYPE_LONGVARBINARY:
		typ = TYPE_BLOB
	}
	c.typ = typ
	return c
}

func (c *Column) assertDefinedType() {
	if c.typ == "" {
		panic("column type is not defined")
	}
}

func (c *Column) Unsigned() *Column {
	c.assertDefinedType()
	switch c.typ {
	case TYPE_SMALLINT, TYPE_INTEGER, TYPE_BIGINT:
	case TYPE_FLOAT:
	case TYPE_DECIMAL, TYPE_NUMERIC:
	default:
		return c
	}

	c.unsigned = true
	return c
}

func (c *Column) SetScaleAndPrecision(scale, precision int) *Column {
	c.assertDefinedType()
	switch c.typ {
	case TYPE_DECIMAL, TYPE_NUMERIC:
		c.scale = scale
		c.precision = precision
	}
	return c
}

func (c *Column) SetLength(size int) *Column {
	c.assertDefinedType()
	switch c.typ {
	case TYPE_CHAR, TYPE_VARCHAR,
		TYPE_LONGVARCHAR, TYPE_CLOB,
		TYPE_TEXT, TYPE_BLOB, TYPE_VARBINARY:

		c.length = size
	}
	return c
}

func (c *Column) Default(val string) *Column {
	c.defaultVal = val
	return c
}

func (c *Column) NullAble() *Column {
	c.nullAble = true
	return c
}

func (c *Column) Primary() *Column {
	c.primary = true
	return c
}

func (c *Column) SetPrimaryPosition(pos int) *Column {
	c.primaryPos = pos
	return c
}

func (c *Column) Comment(comment string) *Column {
	c.comment = comment
	return c
}

func (c *Column) AutoIncrement() *Column {
	c.identity = true
	return c
}

func (c *Column) Identity() *Column {
	c.identity = true
	return c
}

func (c *Column) GetType() string {
	switch c.typ {
	case TYPE_SMALLINT, TYPE_INTEGER, TYPE_BIGINT:
	case TYPE_DECIMAL, TYPE_NUMERIC:
		if c.precision == 0 {
			c.precision = 10
		}
		return fmt.Sprintf("%s (%d, %d)", c.typ, c.precision, c.scale)
	case TYPE_TEXT, TYPE_BLOB, TYPE_VARBINARY:
		if c.length == 0 {
			c.length = DEFAULT_TEXT_SIZE
		}
		if c.length <= 255 {
			if c.typ == TYPE_TEXT {
				return fmt.Sprintf("%s(%d)", TYPE_VARCHAR, c.length)
			} else {
				return fmt.Sprintf("%s(%d)", TYPE_VARBINARY, c.length)
			}
		}
	}
	return c.typ
}

func (c *Column) parseTextSize(size string) int64 {
	if size == "" {
		return DEFAULT_TEXT_SIZE
	}

	size = strings.TrimSpace(size)
	delta, err := strconv.ParseInt(size, 10, 0)
	if err != nil {
		panic(err)
	}
	last := strings.ToLower(size[len(size)-1 : len(size)])
	switch last[0] {
	case 'k':
		delta = delta * 1024
	case 'm':
		delta *= 1024 * 1024
	case 'g':
		delta *= 1024 * 1024 * 1024
	}

	if delta > MAX_TEXT_SIZE {
		return MAX_TEXT_SIZE
	}

	return delta
}

func (c *Column) GetUnsigned() string {
	if c.unsigned {
		return "UNSIGNED"
	}
	return ""
}

func (c *Column) GetNullAble() string {
	if c.unsigned {
		return "NULL"
	} else {
		return "NOT NULL"
	}
}

func (c *Column) GetDefault() string {
	if c.defaultVal != "" {
		return fmt.Sprintf("default %s", c.defaultVal)
	}
	return ""
}

func (c *Column) GetIdentity() string {
	if c.identity {
		return "auto_increment"
	}
	return ""
}

func (c *Column) GetComment() string {
	if c.comment == "" {
		c.comment = c.name
	}
	return c.comment
}

func (c *Column) GetPrimaryPostion() int {
	return c.primaryPos
}
