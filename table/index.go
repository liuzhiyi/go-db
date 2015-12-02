package table

const INDEX_TYPE_PRIMARY = "primary"
const INDEX_TYPE_UNIQUE = "unique"
const INDEX_TYPE_INDEX = "index"
const INDEX_TYPE_FULLTEXT = "fulltext"

type index struct {
	name    string
	columns []string
	typ     string
}

func (idx *index) GetName() string {
	return idx.name
}

func (idx *index) GetType() string {
	return idx.typ
}

func (idx *index) GetColumns() []string {
	return idx.columns
}
