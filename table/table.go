package Table

import (
	"strconv"
	"strings"

	"github.com/liuzhiyi/utils/str"
)

const (
	TYPE_BOOLEAN   = "boolean"
	TYPE_SMALLINT  = "smallint"
	TYPE_INTEGER   = "integer"
	TYPE_BIGINT    = "bigint"
	TYPE_FLOAT     = "float"
	TYPE_NUMERIC   = "numeric"
	TYPE_DECIMAL   = "decimal"
	TYPE_DATE      = "date"
	TYPE_TIMESTAMP = "timestamp" // Capable to support date-time from 1970 + auto-triggers in some RDBMS
	TYPE_DATETIME  = "datetime"  // Capable to support long date-time before 1970
	TYPE_TEXT      = "text"
	TYPE_BLOB      = "blob"      // Used for back compatibility, when query param can"t use statement options
	TYPE_VARBINARY = "varbinary" // A real blob, stored as binary inside DB

	// Deprecated column types, support is left only in MySQL adapter.
	TYPE_TINYINT       = "tinyint"     // Internally converted to TYPE_SMALLINT
	TYPE_CHAR          = "char"        // Internally converted to TYPE_TEXT
	TYPE_VARCHAR       = "varchar"     // Internally converted to TYPE_TEXT
	TYPE_LONGVARCHAR   = "longvarchar" // Internally converted to TYPE_TEXT
	TYPE_CLOB          = "cblob"       // Internally converted to TYPE_TEXT
	TYPE_DOUBLE        = "double"      // Internally converted to TYPE_FLOAT
	TYPE_REAL          = "real"        // Internally converted to TYPE_FLOAT
	TYPE_TIME          = "time"        // Internally converted to TYPE_TIMESTAMP
	TYPE_BINARY        = "binary"      // Internally converted to TYPE_BLOB
	TYPE_LONGVARBINARY = "longvarbinary"
)

type Table struct {
	tableName    string
	schemaName   string
	tableComment string
	columns      map[string]map[string]string
	indexes      []string
	foreignKeys  []string
	options      string
}

func (t *Table) AddColumn(name, ft, size, comment string, options ...string) {
	position := strconv.Itoa(len(t.columns))
	isDefault := "false"
	isnullable := "true"
	length := "0"
	scale := "0"
	precision := "0"
	unsigned := "false"
	primary := "false"
	primaryPosition := "0"
	identity := "false"

	// Convert deprecated types
	switch ft {
	// case TYPE_CHAR, TYPE_VARCHAR, TYPE_LONGVARCHAR, TYPE_CLOB:
	// 	ft = TYPE_TEXT
	case TYPE_TINYINT:
		ft = TYPE_SMALLINT
	case TYPE_DOUBLE, TYPE_REAL:
		ft = TYPE_FLOAT
	case TYPE_TIME:
		ft = TYPE_TIMESTAMP
	case TYPE_BINARY, TYPE_LONGVARBINARY:
		ft = TYPE_BLOB
	}

	// Prepare different properties
	switch ft {
	case TYPE_BOOLEAN:
		break
	case TYPE_SMALLINT, TYPE_INTEGER, TYPE_BIGINT:
		if str.InArray("unsigned", options) {
			unsigned = "true"
		}

	case TYPE_FLOAT:
		if str.InArray("unsigned", options) {
			unsigned = "true"
		}

	case TYPE_DECIMAL, TYPE_NUMERIC:
		scale = "10"
		precision = "0"
		// parse size value
		vals := strings.Split(size, ".")
		precision = vals[0]
		if len(vals) == 2 {
			scale = vals[1]
		}

		if str.InArray("unsigned", options) {
			unsigned = "true"
		}
		break
	case TYPE_DATE, TYPE_DATETIME, TYPE_TIMESTAMP:
		break
	case TYPE_CHAR, TYPE_VARCHAR,
		TYPE_LONGVARCHAR, TYPE_CLOB,
		TYPE_TEXT, TYPE_BLOB, TYPE_VARBINARY:

		length = size
		break
	default:
		panic("Invalid column data type" + ft)
	}

	if str.InArray("default", options) {
		isDefault = "true"
	}
	if str.InArray("nullable", options) {
		isnullable = "true"
	}
	if str.InArray("primary", options) {
		primary = "true"
	}
	if str.InArray("identity", options) || str.InArray("auto_increment", options) {
		identity = "true"
	}

	upperName := strings.ToUpper(name)
	column := make(map[string]string)
	column["COLUMN_NAME"] = name
	column["COLUMN_TYPE"] = ft
	column["COLUMN_POSITION"] = position
	column["DATA_TYPE"] = ft
	column["DEFAULT"] = isDefault
	column["NULLABLE"] = isnullable
	column["LENGTH"] = length
	column["SCALE"] = scale
	column["PRECISION"] = precision
	column["UNSIGNED"] = unsigned
	column["PRIMARY"] = primary
	column["PRIMARY_POSITION"] = primaryPosition
	column["IDENTITY"] = identity
	column["COMMENT"] = comment
	t.columns[upperName] = column
}
