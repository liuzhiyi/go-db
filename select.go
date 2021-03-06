package db

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/liuzhiyi/go-db/adapter"
	"github.com/liuzhiyi/utils/str"
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
	adapter     adapter.Adapter
	bind        map[string]interface{}
	parts       map[string]interface{}
	orderTables []string
	joinTypes   []string
	tableCols   []string
}

func NewSelect(a adapter.Adapter) *Select {
	s := new(Select)
	s._init(a)
	return s
}

func (s *Select) _init(a adapter.Adapter) {
	s.adapter = a
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
	s._initDistinctPart()
	s._initFromPart()
	s._initColumnsPart()
	s._initWherePart()
	s._initHavingPart()
	s._initGroupPart()
	s._initOrderPart()
	s._initCountPart()
	s._initOffsetPart()
	s._initUnionPart()
}

func (s *Select) _initDistinctPart() {
	s.parts[DISTINCT] = false
}

func (s *Select) _initFromPart() {
	s.parts[FROM] = make(map[string]map[string]string)
}

func (s *Select) _initColumnsPart() {
	s.parts[COLUMNS] = [][]string{}
}

func (s *Select) _initWherePart() {
	s.parts[WHERE] = []string{}
}

func (s *Select) _initHavingPart() {
	s.parts[HAVING] = []string{}
}

func (s *Select) _initGroupPart() {
	s.parts[GROUP] = []string{}
}

func (s *Select) _initOrderPart() {
	s.parts[ORDER] = []string{}
}

func (s *Select) _initUnionPart() {
	s.parts[UNION] = [][]string{}
}

func (s *Select) _initCountPart() {
	s.parts[LIMIT_COUNT] = int64(0)
}

func (s *Select) _initOffsetPart() {
	s.parts[LIMIT_OFFSET] = int64(0)
}

func (s *Select) Distinct(flag bool) *Select {
	s.parts[DISTINCT] = flag
	return s
}

/**
*@param name:"table1", alias is "table1"; "table1 as t", alias is "t"
*@param cols:"col1, col2 as c, o.col3 ..."
**/
func (s *Select) From(name, cols, schema string) *Select {
	return s._join(FROM, "", cols, schema, name)
}

/**
*@param name:"table1", alias is "table1"; "table1 as t", alias is "t"
*@param cols:"col1, col2 as c, o.col3 ..."
**/
func (s *Select) Join(name, cond, cols, schema string) {
	s._join(INNER_JOIN, cond, cols, schema, name)
}

/**
*@param name:"table1", alias is "table1"; "table1 as t", alias is "t"
*@param cols:"col1, col2 as c, o.col3 ..."
**/
func (s *Select) JoinLeft(name, cond, cols, schema string) {
	s._join(LEFT_JOIN, cond, cols, schema, name)
}

/**
*@param name:"table1", alias is "table1"; "table1 as t", alias is "t"
*@param cols:"col1, col2 as c, o.col3 ..."
**/
func (s *Select) JoinRight(name, cond, cols, schema string) {
	s._join(RIGHT_JOIN, cond, cols, schema, name)
}

/**
*@param name:"table1", alias is "table1"; "table1 as t", alias is "t"
*@param cols:"col1, col2 as c, o.col3 ..."
**/
func (s *Select) JoinFull(name, cond, cols, schema string) {
	s._join(FULL_JOIN, cond, cols, schema, name)
}

/**
*@param name:"table1", alias is "table1"; "table1 as t", alias is "t"
*@param cols:"col1, col2 as c, o.col3 ..."
**/
func (s *Select) JoinCross(name, cond, cols, schema string) {
	s._join(CROSS_JOIN, cond, cols, schema, name)
}

/**
*@param name:"table1", alias is "table1"; "table1 as t", alias is "t"
*@param cols:"col1, col2 as c, o.col3 ..."
**/
func (s *Select) JoinNatural(name, cond, cols, schema string) {
	s._join(NATURAL_JOIN, cond, cols, schema, name)
}

func (s *Select) FieldExpre(str string) string {
	set := s._prepareCols(str)
	for i, field := range set {
		field = strings.TrimSpace(field)
		col, alis := s.getAlisFromString(str)
		set[i] = fmt.Sprintf("{%s} as %s", col, alis)
	}

	return strings.Join(set, ",")
}

func (s *Select) isExpre(str string) bool {
	return strings.HasPrefix(str, "{") && strings.HasSuffix(str, "}")
}

func (s *Select) ReExpre(str string) string {
	if len(str) > 2 {
		str = str[1 : len(str)-1]
	}

	return str
}

func (s *Select) Union(t string, set ...Select) {
	if t != SQL_UNION || t != SQL_UNION_ALL {
		t = SQL_UNION
	}
	unionPart := s.parts[UNION].([][]string)
	for i := 0; i < len(set); i++ {
		unionPart = append(unionPart, []string{set[i].Assemble(), t})
	}
	s.parts[SQL_UNION] = unionPart
}

func (s *Select) Columns(cols interface{}, correlationName string) *Select {
	if correlationName == "" && len(s.parts[FROM].(map[string]map[string]string)) > 0 {
		correlationName = ""
	} else if _, ok := s.parts[correlationName]; ok {
		panic("No table has been specified for the FROM clause")
	}
	s._tableCols(correlationName, cols)
	return s
}

func (s *Select) _prepareCols(cols string) []string {
	return strings.Split(cols, ",")
}

func (s *Select) Group(spec ...string) {
	groupPart := s.parts[GROUP].([]string)
	for i := 0; i < len(spec); i++ {
		groupPart = append(groupPart, spec[i])
	}
	s.parts[GROUP] = groupPart
}

func (s *Select) Where(cond string, values ...interface{}) {
	wherePart := s.parts[WHERE].([]string)
	wherePart = append(wherePart, s._where(cond, true, values...))
	s.parts[WHERE] = wherePart
}

func (s *Select) OrWhere(cond string, values ...interface{}) {
	wherePart := s.parts[WHERE].([]string)
	wherePart = append(wherePart, s._where(cond, false, values...))
	s.parts[WHERE] = wherePart
}

func (s *Select) Having(cond string, value interface{}) {
	havingPart := s.parts[HAVING].([]string)
	if value != nil {
		cond = s.adapter.QuoteInto(cond, value)
	}
	if len(havingPart) > 0 {
		havingPart = append(havingPart, fmt.Sprintf("%s (%s)", SQL_AND, cond))
	} else {
		havingPart = append(havingPart, fmt.Sprintf("(%s)", cond))
	}
	s.parts[HAVING] = havingPart
}

func (s *Select) Order(spec ...string) {
	orderPart := s.parts[ORDER].([]string)
	direction := SQL_ASC
	reg := regexp.MustCompile(`(.*)\W((?i:asc|desc))\b`)
	for i := 0; i < len(spec); i++ {
		col := spec[i]
		if m := reg.FindStringSubmatch(spec[i]); len(m) == 3 {
			col = m[1]
			direction = m[2]
		}
		orderPart = append(orderPart, s.adapter.QuoteIdentifier(col)+" "+strings.ToUpper(direction))
	}
	s.parts[ORDER] = orderPart
}

func (s *Select) Limit(count, offset int64) *Select {
	s.parts[LIMIT_COUNT] = count
	s.parts[LIMIT_OFFSET] = offset
	return s
}

func (s *Select) LimitPage(page, rowCount int64) {
	if page <= 0 {
		page = 1
	}
	if rowCount <= 0 {
		rowCount = 1
	}
	offset := rowCount * (page - 1)
	s.Limit(rowCount, offset)
}

func (s *Select) Assemble() string {
	sql := SQL_SELECT
	sql = s._renderColumns(sql)
	sql = s._renderFrom(sql)
	sql = s._renderWhere(sql)
	sql = s._renderGroup(sql)
	sql = s._renderHaving(sql)
	sql = s._renderOrder(sql)
	sql = s._renderLimit(sql)
	sql = s._renderUnion(sql)
	return sql
}

func (s *Select) Reset(part ...string) {
	if len(part) > 0 {
		for i := 0; i < len(part); i++ {
			switch part[i] {
			case DISTINCT:
				s._initDistinctPart()
			case FROM:
				s._initWherePart()
			case COLUMNS:
				s._initColumnsPart()
			case HAVING:
				s._initHavingPart()
			case WHERE:
				s._initWherePart()
			case GROUP:
				s._initGroupPart()
			case ORDER:
				s._initOrderPart()
			case LIMIT_COUNT:
				s._initCountPart()
			case LIMIT_OFFSET:
				s._initOffsetPart()
			case UNION:
				s._initUnionPart()
			}
		}
	} else {
		s._initPart()
	}

}

func (s *Select) _join(joinType, cond, cols, schema string, name interface{}) *Select {
	if !str.InArray(joinType, s.joinTypes) && joinType != FROM {
		panic("Invalid join type " + joinType)
	}

	var correlationName, tableName string
	switch val := name.(type) {
	case string:
		reg := regexp.MustCompile(`^(.+)\s+[aA][sS]\s+(.+)$`)
		if val == "" {
			correlationName = ""
			tableName = ""
		} else if m := reg.FindStringSubmatch(val); len(m) == 3 {
			correlationName = m[2]
			tableName = m[1]
		} else {
			correlationName = s._uniqueCorrelation(val)
			tableName = val
		}
	case [2]string:
		correlationName = val[0]
		tableName = val[1]
	default:
		panic("Invalid params")
	}

	if strings.IndexByte(tableName, '.') > 0 && !s.isExpre(tableName) {
		tmp := strings.Split(tableName, ".")
		schema = tmp[0]
		tableName = tmp[1]
	}

	if correlationName != "" {
		fromPart := s.parts[FROM].(map[string]map[string]string)
		from := make(map[string]string)
		if _, ok := fromPart[correlationName]; ok {
			panic("You cannot define a correlation name " + correlationName + " more than once")
		}
		from["joinType"] = joinType
		from["schema"] = schema
		from["tableName"] = tableName
		from["joinCondition"] = cond
		fromPart[correlationName] = from
		s.parts[FROM] = fromPart
		s.orderTables = append(s.orderTables, correlationName)
		s._tableCols(correlationName, cols)
	}
	return s
}

func (s *Select) _uniqueCorrelation(name string) string {
	if pos := strings.IndexByte(name, '.'); pos > 0 {
		name = name[pos:]
	}
	fromPart := s.parts[FROM].(map[string]map[string]string)
	i := 2
	_, ok := fromPart[name]
	for ok {
		name = fmt.Sprintf("%s_%s", name, strconv.Itoa(i))
		_, ok = fromPart[name]
		i++
	}
	return name
}

func (s *Select) getAlisFromString(str string) (col, alias string) {
	re := regexp.MustCompile(`^(.+)\s+[aA][sS]\s+(.+)$`)
	if m := re.FindStringSubmatch(str); len(m) == 3 && m[1] != "" && m[2] != "" {
		col = m[1]
		alias = m[2]
	} else {
		col = str
		alias = ""
	}

	return
}

func (s *Select) getCorrelationNameFromString(str string) (col, correlationName string) {
	re := regexp.MustCompile(`(.+)\.(.+)`)
	if m := re.FindStringSubmatch(str); len(m) == 3 && m[1] != "" && m[2] != "" {
		correlationName = m[1]
		col = m[2]
	} else {
		col = str
		correlationName = ""
	}

	return
}

func (s *Select) _tableCols(correlationName string, express interface{}) {
	columnPart := s.parts[COLUMNS].([][]string)
	switch e := express.(type) {
	case string:
		cols := s._prepareCols(e)
		for _, col := range cols {
			var alias string
			if col == "" {
				continue
			}
			col = strings.TrimSpace(col)
			currentCorrelationName := correlationName
			col, alias = s.getAlisFromString(col)

			if !s.isExpre(col) {
				re := regexp.MustCompile(`(.+)\.(.+)`)
				if m := re.FindStringSubmatch(col); len(m) == 3 && m[1] != "" && m[2] != "" {
					currentCorrelationName = m[1]
					col = m[2]
				}
			} else {
				currentCorrelationName = ""
			}

			columnPart = append(columnPart, []string{currentCorrelationName, col, alias})
		}
	case Select:
		columnPart = append(columnPart, []string{"", fmt.Sprintf("(%s)", e.Assemble()), correlationName})
	default:
		panic("invalid col!")
	}
	s.parts[COLUMNS] = columnPart
}

func (s *Select) GetColumnPart() [][]string {
	columnPart := s.parts[COLUMNS].([][]string)
	return columnPart
}

func (s *Select) _where(condition string, flag bool, values ...interface{}) string {
	cond := ""
	if len(values) > 0 {
		condition = s.adapter.QuoteInto(condition, values...)
	}
	if len(s.parts[WHERE].([]string)) > 0 {
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
	var columns []string
	for _, colEntity := range columnPart {
		correlationName := colEntity[0]
		col := colEntity[1]
		alias := colEntity[2]
		if col == SQL_WILDCARD {
			alias = ""
		}

		if s.isExpre(col) {
			col = s.ReExpre(col)
			columns = append(columns, fmt.Sprintf("%s AS %s", col, s.adapter.QuoteIdentifier(alias)))
		} else {
			if correlationName != "" {
				columns = append(columns, s.adapter.QuoteIdentifierAs(fmt.Sprintf("%s.%s", correlationName, col), alias))
			} else {
				columns = append(columns, s.adapter.QuoteIdentifierAs(col, alias))
			}
		}

	}
	return sql + " " + strings.Join(columns, ", ")
}

func (s *Select) _renderFrom(sql string) string {
	fromPart := s.parts[FROM].(map[string]map[string]string)

	var from []string
	var table map[string]string
	for _, correlationName := range s.orderTables {
		table = fromPart[correlationName]
		tmp := ""
		joinType := table["joinType"]

		if joinType == FROM {
			joinType = INNER_JOIN
		}

		if len(from) > 0 {
			tmp += fmt.Sprintf(" %s ", strings.ToUpper(joinType))
		}

		tmp += s._getQuotedTable(table["tableName"], correlationName)
		if len(from) > 0 && table["joinCondition"] != "" {
			tmp += fmt.Sprintf(" %s %s ", SQL_ON, table["joinCondition"])
		}

		from = append(from, tmp)
	}

	if len(from) > 0 {
		sql += fmt.Sprintf(" %s %s", SQL_FROM, strings.Join(from, "\n"))
	}
	return sql
}

func (s *Select) _getQuotedTable(tableName, correlationName string) string {
	if s.isExpre(tableName) {
		tableName = s.ReExpre(tableName)
	} else {
		tableName = s.adapter.QuoteIdentifier(tableName)
	}
	correlationName = s.adapter.QuoteIdentifier(correlationName)
	return fmt.Sprintf("%s %s %s", tableName, SQL_AS, correlationName)
}

func (s *Select) _renderUnion(sql string) string {
	unionPart := s.parts[UNION].([][]string)
	l := len(unionPart)
	for i := 0; i < l; i++ {
		sql += " " + unionPart[i][1] + " "
		sql += unionPart[i][0]
	}
	return sql
}

func (s *Select) _renderWhere(sql string) string {
	wherePart := s.parts[WHERE].([]string)
	if len(wherePart) > 0 {
		sql += fmt.Sprintf(" %s %s ", SQL_WHERE, strings.Join(wherePart, " "))
	}
	return sql
}

func (s *Select) _renderGroup(sql string) string {
	groupPart := s.parts[GROUP].([]string)
	if len(groupPart) > 0 {
		sql += fmt.Sprintf(" %s %s ", SQL_GROUP_BY, strings.Join(groupPart, ",\n\t"))
	}
	return sql
}

func (s *Select) _renderHaving(sql string) string {
	havingPart := s.parts[HAVING].([]string)
	if len(havingPart) > 0 {
		sql += " " + SQL_HAVING + " " + strings.Join(havingPart, " ")
	}
	return sql
}

func (s *Select) _renderOrder(sql string) string {
	orderPart := s.parts[ORDER].([]string)
	if len(orderPart) > 0 {
		sql += " " + SQL_ORDER_BY + " " + strings.Join(orderPart, ", ")
	}
	return sql
}

func (s *Select) _renderLimit(sql string) string {
	count := s.parts[LIMIT_COUNT].(int64)
	offset := s.parts[LIMIT_OFFSET].(int64)
	if count > 0 {
		sql = s.adapter.Limit(sql, count, offset)
	}
	return sql
}

func (s *Select) GetCountSql() string {
	sql := SQL_SELECT
	sql += " COUNT(*) "
	sql = s._renderFrom(sql)
	sql = s._renderWhere(sql)
	return sql
}

func (s *Select) Clone() *Select {
	n := new(Select)
	n._init(s.adapter)
	/***
	 *拷贝数据
	 **/
	return n
}
