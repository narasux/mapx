package mapx_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/narasux/mapx"
)

// SetItems success case
func TestSetItemsSuccessCase(t *testing.T) {
	// depth 1，val type int
	err := mapx.SetItems(deploySpec, "intKey4SetItem", 5)
	assert.Nil(t, err)
	ret, _ := mapx.GetItems(deploySpec, []string{"intKey4SetItem"})
	assert.Equal(t, 5, ret)

	// depth 2, val type string
	err = mapx.SetItems(deploySpec, "strategy.type", "Rolling")
	assert.Nil(t, err)
	ret, _ = mapx.GetItems(deploySpec, []string{"strategy", "type"})
	assert.Equal(t, "Rolling", ret)

	// depth 3, val type string
	err = mapx.SetItems(deploySpec, []string{"template", "spec", "restartPolicy"}, "Never")
	assert.Nil(t, err)
	ret, _ = mapx.GetItems(deploySpec, []string{"template", "spec", "restartPolicy"})
	assert.Equal(t, "Never", ret)

	// key noy exists
	err = mapx.SetItems(deploySpec, []string{"selector", "testKey"}, "testVal")
	assert.Nil(t, err)
	ret, _ = mapx.GetItems(deploySpec, "selector.testKey")
	assert.Equal(t, "testVal", ret)
}

// SetItems fail case
func TestSetItemsFailCase(t *testing.T) {
	// invalid paths type error
	err := mapx.SetItems(deploySpec, 0, 1)
	assert.True(t, errors.Is(err, mapx.ErrInvalidPathType))

	// not paths error
	err = mapx.SetItems(deploySpec, []string{}, 1)
	assert.NotNil(t, err)

	// not map[string]interface{} type error
	err = mapx.SetItems(deploySpec, []string{"replicas", "testKey"}, 1)
	assert.NotNil(t, err)

	// key not exist
	err = mapx.SetItems(deploySpec, []string{"templateKey", "spec"}, 1)
	assert.NotNil(t, err)

	err = mapx.SetItems(deploySpec, "templateKey.spec", 123)
	assert.NotNil(t, err)

	// paths type error
	err = mapx.SetItems(deploySpec, []int{123, 456}, 1)
	assert.NotNil(t, err)

	err = mapx.SetItems(deploySpec, 123, 1)
	assert.NotNil(t, err)
}
