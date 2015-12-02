package table

type foreignKey struct {
	fkName        string
	columnName    string
	refTableName  string
	refColumnName string
	onDelete      string
	onUpdate      string
}

func (fk *foreignKey) GetName() string {
	return fk.fkName
}

func (fk *foreignKey) GetColumnName() string {
	return fk.columnName
}

func (fk *foreignKey) GetRefTableName() string {
	return fk.refTableName
}

func (fk *foreignKey) GetRefColumnName() string {
	return fk.refColumnName
}

func (fk *foreignKey) GetOnDelete() string {
	return fk.onDelete
}

func (fk *foreignKey) GetOnUpdate() string {
	return fk.onUpdate
}
