package mapx_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/narasux/mapx"
)

func TestExists(t *testing.T) {
	obj := map[string]any{
		"key_1": "val_1",
		"key_3": "val_3",
	}
	assert.True(t, mapx.Exists(obj, "key_1"))
	assert.False(t, mapx.Exists(obj, "key_2"))
	assert.True(t, mapx.Exists(obj, "key_3"))
}
