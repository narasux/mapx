package mapx

// Exists check if key exists in map
func Exists(obj map[string]any, key string) bool {
	_, ok := obj[key]
	return ok
}
