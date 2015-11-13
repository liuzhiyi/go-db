package table

type foreignKey struct {
	fkName        string
	columnName    string
	refTableName  string
	refColumnName string
	onDelete      string
	onUpdate      string
}
