package table

import (
	"strings"
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

	/**
	 * Default and maximal TEXT and BLOB columns sizes we can support for different DB systems.
	 */
	DEFAULT_TEXT_SIZE  = 1024
	MAX_TEXT_SIZE      = 2147483648
	MAX_VARBINARY_SIZE = 2147483648

	/**
	 * Default values for timestampses - fill with current timestamp on inserting record, on changing and both cases
	 */
	TIMESTAMP_INIT_UPDATE = "TIMESTAMP_INIT_UPDATE"
	TIMESTAMP_INIT        = "TIMESTAMP_INIT"
	TIMESTAMP_UPDATE      = "TIMESTAMP_UPDATE"

	/**
	 * Actions used for foreign keys
	 */
	ACTION_CASCADE     = "CASCADE"
	ACTION_SET_NULL    = "SET NULL"
	ACTION_NO_ACTION   = "NO ACTION"
	ACTION_RESTRICT    = "RESTRICT"
	ACTION_SET_DEFAULT = "SET DEFAULT"
)

type Table struct {
	tableName    string
	schemaName   string
	tableComment string
	columns      []*Column
	indexes      []*index
	foreignKeys  []*foreignKey
	options      string
}

func (t *Table) SetName(name string) *Table {
	t.tableName = name
	if t.tableComment == "" {
		t.tableComment = name
	}
	return t
}

func (t *Table) SetSchema(name string) *Table {
	t.schemaName = name
	return t
}

func (t *Table) SetComment(comment string) *Table {
	t.tableComment = comment
	return t
}

func (t *Table) GetName() string {
	if t.tableName == "" {
		panic("table name is not defined")
	}
	return t.tableName
}

func (t *Table) GetSchema() string {
	return t.schemaName
}

func (t *Table) GetComment() string {
	return t.tableComment
}

func (t *Table) NewColumn(name string) *Column {
	c := new(Column)
	c.name = strings.ToLower(name)
	return c
}

func (t *Table) AddColumn(c *Column) {
	c.assertDefinedType()
	if c.primary {
		if c.primaryPos > 0 {

		} else {
			pos := 0
			for _, v := range t.columns {
				if v.primary {
					pos++
				}
			}
			c.SetPrimaryPosition(pos)
		}
	}
	t.columns = append(t.columns, c)
}

func (t *Table) AddForeignKey(fkname, col, refTable, refColumn, onDelete, onUpdate string) *foreignKey {
	fkname = strings.ToLower(fkname)

	switch onDelete {
	case ACTION_CASCADE:
	case ACTION_RESTRICT:
	case ACTION_SET_DEFAULT:
	case ACTION_SET_NULL:
	default:
		onDelete = ACTION_NO_ACTION
	}

	switch onUpdate {
	case ACTION_CASCADE:
	case ACTION_RESTRICT:
	case ACTION_SET_DEFAULT:
	case ACTION_SET_NULL:
	default:
		onUpdate = ACTION_NO_ACTION
	}

	fk := &foreignKey{
		fkName:        fkname,
		columnName:    col,
		refTableName:  refTable,
		refColumnName: refColumn,
		onDelete:      onDelete,
		onUpdate:      onUpdate,
	}
	t.foreignKeys = append(t.foreignKeys, fk)
	return fk
}

func (t *Table) AddIndex(name string, fields []string, typ string) {
	// pos := 0
	// for _, field := range fields {

	// }

	// switch typ {
	// case INDEX_TYPE_FULLTEXT:
	// case INDEX_TYPE_PRIMARY:
	// case INDEX_TYPE_UNIQUE:
	// default:
	// 	typ = INDEX_TYPE_INDEX
	// }
	// idx := &index{
	// 	name:    name,
	// 	typ:     typ,
	// 	cloumns: fields,
	// }
	// t.indexes = append(t.indexes, idx)
	// idxType := INDEX_TYPE_INDEX
	// pos := 0

	// $idxType    = Varien_Db_Adapter_Interface::INDEX_TYPE_INDEX;
	//        $position   = 0;
	//        $columns    = array();
	//        if (!is_array($fields)) {
	//            $fields = array($fields);
	//        }

	//        foreach ($fields as $columnData) {
	//            $columnSize = null;
	//            $columnPos  = $position;
	//            if (is_string($columnData)) {
	//                $columnName = $columnData;
	//            } else if (is_array($columnData)) {
	//                if (!isset($columnData['name'])) {
	//                    throw new Zend_Db_Exception('Invalid index column data');
	//                }

	//                $columnName = $columnData['name'];
	//                if (!empty($columnData['size'])) {
	//                    $columnSize = (int)$columnData['size'];
	//                }
	//                if (!empty($columnData['position'])) {
	//                    $columnPos = (int)$columnData['position'];
	//                }
	//            } else {
	//                continue;
	//            }

	//            $columns[strtoupper($columnName)] = array(
	//                'NAME'      => $columnName,
	//                'SIZE'      => $columnSize,
	//                'POSITION'  => $columnPos
	//            );

	//            $position ++;
	//        }

	//        if (empty($columns)) {
	//            throw new Zend_Db_Exception('Columns for index are not defined');
	//        }

	//        if (!empty($options['type'])) {
	//            $idxType = $options['type'];
	//        }

	//        $this->_indexes[strtoupper($indexName)] = array(
	//            'INDEX_NAME'    => $indexName,
	//            'COLUMNS'       => $this->_normalizeIndexColumnPosition($columns),
	//            'TYPE'          => $idxType
	//        );

	//        return $this;
}

func (t *Table) GetColumns() []*Column {
	return t.columns
}

func (t *Table) GetIndexs() []*index {
	return t.indexes
}

func (t *Table) GetForeginKeys() []*foreignKey {
	return t.foreignKeys
}
