package db

import (
	"regexp"
	"strings"
)

const (
	DISTINCT     = "distinct"
	COLUMNS      = "columns"
	FROM         = "from"
	UNION        = "union"
	WHERE        = "where"
	GROUP        = "group"
	HAVING       = "having"
	ORDER        = "order"
	LIMIT_COUNT  = "limitcount"
	LIMIT_OFFSET = "limitoffset"
	FOR_UPDATE   = "forupdate"

	INNER_JOIN   = "inner join"
	LEFT_JOIN    = "left join"
	RIGHT_JOIN   = "right join"
	FULL_JOIN    = "full join"
	CROSS_JOIN   = "cross join"
	NATURAL_JOIN = "natural join"

	SQL_WILDCARD   = "*"
	SQL_SELECT     = "SELECT"
	SQL_UNION      = "UNION"
	SQL_UNION_ALL  = "UNION ALL"
	SQL_FROM       = "FROM"
	SQL_WHERE      = "WHERE"
	SQL_DISTINCT   = "DISTINCT"
	SQL_GROUP_BY   = "GROUP BY"
	SQL_ORDER_BY   = "ORDER BY"
	SQL_HAVING     = "HAVING"
	SQL_FOR_UPDATE = "FOR UPDATE"
	SQL_AND        = "AND"
	SQL_AS         = "AS"
	SQL_OR         = "OR"
	SQL_ON         = "ON"
	SQL_ASC        = "ASC"
	SQL_DESC       = "DESC"
)

type Select struct {
	bind      map[string]interface{}
	parts     map[string]interface{}
	joinTypes []string
	tableCols []string
}

func (s *Select) _init() {
	//初始化可用连接类型
	s.joinTypes = append(s.joinTypes, INNER_JOIN)
	s.joinTypes = append(s.joinTypes, LEFT_JOIN)
	s.joinTypes = append(s.joinTypes, RIGHT_JOIN)
	s.joinTypes = append(s.joinTypes, FULL_JOIN)
	s.joinTypes = append(s.joinTypes, CROSS_JOIN)
	s.joinTypes = append(s.joinTypes, NATURAL_JOIN)

	s._initPart()
}

//初始化组装部分
func (s *Select) _initPart() {
	s.parts = make(map[string]interface{})
	s.parts[DISTINCT] = false
	s.parts[FROM] = make(map[string]map[string]string)
	s.parts[COLUMNS] = [][]string{}
	s.parts[WHERE] = ""
	s.parts[LIMIT_COUNT] = 0
	s.parts[LIMIT_OFFSET] = 0
}

func (s *Select) Distinct(flag bool) *Select {
	s.parts[DISTINCT] = flag
	return s
}

func (s *Select) From(name, cols, schema string) *Select {
	return s._join(FROM, "", name, cols, schema)
}

func (s *Select) Columns(cols, correlationName string) *Select {
	if correlationName == "" && len(s.parts[FROM].(map[string]interface{})) > 0 {
		correlationName = ""
	} else if _, ok := s.parts[correlationName]; ok {
		panic("No table has been specified for the FROM clause")
	}

	s._tableCols(correlationName, []string{cols})
	return s
}

func (s *Select) _join(joinType, cond, cols, schema string, name interface{}) *Select {
	if !inArray(joinType, s.joinTypes) && joinType != FROM {
		panic("Invalid join type " + joinType)
	}

	if _, ok := s.parts[UNION]; ok {
		panic("Invalid use of table with " + UNION)
	}

	var correlationName, tableName string
	switch name.(type) {
	case string:
		if name.(string) == "" {
			correlationName = ""
			tableName = ""
		} else {
			correlationName = name.(string)
			tableName = s._uniqueCorrelation(name.(string))
		}
	case [2]string:
		correlationName = name.([2]string)[0]
		tableName = name.([2]string)[1]
	default:
		panic("Invalid params")
	}

	if strings.IndexByte(tableName, '.') > 0 {
		tmp := strings.Split(tableName, ".")
		schema = tmp[0]
		tableName = tmp[1]
	}

	if correlationName != "" {
		fromPart := s.parts[FROM].(map[string]map[string]string)
		if _, ok := fromPart[correlationName]; ok {
			panic("You cannot define a correlation name " + correlationName + " more than once")
		}
		fromPart[correlationName]["joinType"] = joinType
		fromPart[correlationName]["schema"] = schema
		fromPart[correlationName]["tableName"] = tableName
		fromPart[correlationName]["joinCondition"] = cond
		s.parts[FROM] = fromPart
		s._tableCols(correlationName, []string{cols})
	}
	return s
}

func (s *Select) _uniqueCorrelation(name string) string {
	return name
}

func (s *Select) _tableCols(correlationName string, cols []string) {
	columnPart := s.parts[COLUMNS].([][]string)
	for _, col := range cols {
		var alias string
		currentCorrelationName := correlationName
		re := regexp.MustCompile(`/^(.+)\s+` + SQL_AS + `\s+(.+)$/i`)
		if m := re.FindStringSubmatch(col); m[1] != "" && m[2] != "" {
			col = m[1]
			alias = m[2]
		} else {
			alias = ""
		}
		re = regexp.MustCompile(`/(.+)\.(.+)/`)
		if m := re.FindStringSubmatch(col); m[1] != "" && m[2] != "" {
			currentCorrelationName = m[1]
			col = m[2]
		}
		columnPart = append(columnPart, []string{currentCorrelationName, col, alias})
	}
	s.parts[COLUMNS] = columnPart
}

func (s *Select) _where(condition string, flag bool) string {
	cond := ""
	if s.parts[WHERE].(string) != "" {
		if flag {
			cond = SQL_AND + " "
		} else {
			cond = SQL_OR + " "
		}
	}
	return cond + "(" + condition + ")"
}

func (s *Select) _renderDistinct(sql string) string {
	dis := s.parts[DISTINCT].(bool)
	if dis {
		sql += " " + SQL_DISTINCT
	}
	return sql
}

func (s *Select) _renderColumns(sql string) string {
	columnPart := s.parts[COLUMNS].([][]string)
	if len(columnPart) > 0 {
		return ""
	}
	for _, colEntity := range columnPart {
		//correlationName := colEntity[0]
		col := colEntity[1]
		//alias := colEntity[2]
		if col == SQL_WILDCARD {
			//alias = ""
		}
	}
	return sql
}

func (s *Select) Limit(count, offset int) *Select {
	s.parts[LIMIT_COUNT] = count
	s.parts[LIMIT_OFFSET] = offset
	return s
}
