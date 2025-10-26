package mapx

import (
	"errors"
	"fmt"
	"strings"
)

var ErrInvalidPathType = errors.New("paths's type must one of (string, []string)")

// GetItems gets the Map value of the nested definition
// The paths parameter supports the []string type, such as []string{"metadata", "namespace"}
// or string type (with '.' as the separator), such as "spec.template.spec.containers"
func GetItems(obj map[string]any, paths any) (any, error) {
	switch p := paths.(type) {
	case string:
		return getItems(obj, strings.Split(p, "."))
	case []string:
		return getItems(obj, p)
	default:
		return nil, ErrInvalidPathType
	}
}

func getItems(obj map[string]any, paths []string) (any, error) {
	if len(paths) == 0 {
		return nil, errors.New("paths is empty list")
	}
	ret, exists := obj[paths[0]]
	if !exists {
		return nil, fmt.Errorf("key %s not exist", paths[0])
	}
	if len(paths) == 1 {
		return ret, nil
	} else if subMap, ok := obj[paths[0]].(map[string]any); ok {
		return getItems(subMap, paths[1:])
	}
	return nil, fmt.Errorf("key %s, val not map[string]any type", paths[0])
}

// Get if the key does not exist, return the default value
func Get(obj map[string]any, paths any, defVal any) any {
	ret, err := GetItems(obj, paths)
	if err != nil {
		return defVal
	}
	return ret
}

// GetBool is Get for bool type, default value is false
func GetBool(obj map[string]any, paths any) bool {
	if r, ok := Get(obj, paths, false).(bool); ok {
		return r
	}
	return false
}

// GetInt64 is Get for int64 type, default value is int64(0)
func GetInt64(obj map[string]any, paths any) int64 {
	if r, ok := Get(obj, paths, 0).(int64); ok {
		return r
	}
	return 0
}

// GetFloat64 is Get for float64 type, default value is float64(0)
func GetFloat64(obj map[string]any, paths any) float64 {
	if r, ok := Get(obj, paths, 0).(float64); ok {
		return r
	}
	return 0
}

// GetStr is Get for string type, default value is ""
func GetStr(obj map[string]any, paths any) string {
	if r, ok := Get(obj, paths, "").(string); ok {
		return r
	}
	return ""
}

// GetList is Get for []any type, default value is []any{}
func GetList(obj map[string]any, paths any) []any {
	if r, ok := Get(obj, paths, nil).([]any); ok {
		return r
	}
	return nil
}

// GetMap is Get for map[string]any type, default value is map[string]any{}
func GetMap(obj map[string]any, paths any) map[string]any {
	if r, ok := Get(obj, paths, nil).(map[string]any); ok {
		return r
	}
	return nil
}
