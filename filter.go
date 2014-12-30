package db

type Filter map[string]map[string]interface{}

func NewFilter() Filter {
	f := make(Filter)
	return f
}

func (f Filter) SetCondition(field, key string, value interface{}) {
	if condition, ok := f[field]; ok {
		condition[key] = value
		f[field] = condition
	} else {
		condition := make(map[string]interface{})
		condition[key] = value
		f[field] = condition
	}
}
