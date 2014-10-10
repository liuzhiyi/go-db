package db

func inArray(val string, array []string) bool {
	for _, item := range array {
		if val == item {
			return true
		}
	}
	return false
}
