package mapx

import (
	"fmt"
	"strings"
)

// SetItems assigns values to nested Maps
// The paths parameter supports []string type, such as []string{"metadata", "namespace"}
// or string type (with '.' as separator), such as "spec.template.spec.containers"
func SetItems(obj map[string]any, paths any, val any) error {
	// check paths type
	switch p := paths.(type) {
	case string:
		if err := setItems(obj, strings.Split(p, "."), val); err != nil {
			return err
		}
	case []string:
		if err := setItems(obj, p, val); err != nil {
			return err
		}
	default:
		return ErrInvalidPathType
	}
	return nil
}

func setItems(obj map[string]any, paths []string, val any) error {
	if len(paths) == 0 {
		return fmt.Errorf("paths is empty list")
	}
	if len(paths) == 1 {
		obj[paths[0]] = val
	} else if subMap, ok := obj[paths[0]].(map[string]any); ok {
		return setItems(subMap, paths[1:], val)
	} else {
		return fmt.Errorf("key %s not exists or obj[key] not map[string]any type", paths[0])
	}
	return nil
}
