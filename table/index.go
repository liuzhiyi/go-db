package table

const INDEX_TYPE_PRIMARY = "primary"
const INDEX_TYPE_UNIQUE = "unique"
const INDEX_TYPE_INDEX = "index"
const INDEX_TYPE_FULLTEXT = "fulltext"

type index struct {
	name    string
	cloumns []string
	typ     string
}
