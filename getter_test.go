package mapx_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/narasux/mapx"
)

var deploySpec = map[string]any{
	"testKey":              "testValue",
	"replicas":             3,
	"revisionHistoryLimit": 10,
	"intKey4SetItem":       8,
	"selector": map[string]any{
		"matchLabels": map[string]any{
			"app": "nginx",
		},
	},
	"strategy": map[string]any{
		"rollingUpdate": map[string]any{
			"maxSurge":       "25%",
			"maxUnavailable": "25%",
		},
		"type": "RollingUpdate",
	},
	"template": map[string]any{
		"metadata": map[string]any{
			"creationTimestamp": nil,
			"int64Key4GetInt64": int64(10),
			"labels": map[string]any{
				"app":           "nginx",
				"strKey4GetStr": "value",
			},
		},
		"spec": map[string]any{
			"boolKey4GetBool": true,
			"containers": []map[string]any{
				{
					"image":           "nginx:latest",
					"imagePullPolicy": "IfNotPresent",
					"name":            "nginx",
					"ports": map[string]any{
						"containerPort": 80,
						"protocol":      "TCP",
					},
					"resources": map[string]any{},
				},
			},
			"dnsPolicy":                     "ClusterFirst",
			"restartPolicy":                 "Always",
			"schedulerName":                 "default-scheduler",
			"securityContext":               map[string]any{},
			"terminationGracePeriodSeconds": 30,
		},
		"interfaceList": []any{
			map[string]any{"key": "value"},
			"key-value",
		},
	},
}

// paths 为以 '.' 连接的字符串
func TestGetItems(t *testing.T) {
	// depth 1，val type int
	ret, _ := mapx.GetItems(deploySpec, "replicas")
	assert.Equal(t, 3, ret)

	// depth 2, val type string
	ret, _ = mapx.GetItems(deploySpec, "strategy.type")
	assert.Equal(t, "RollingUpdate", ret)

	// depth 3, val type string
	ret, _ = mapx.GetItems(deploySpec, "template.spec.restartPolicy")
	assert.Equal(t, "Always", ret)
}

// paths 为 []string，成功的情况
func TestGetItemsSuccessCase(t *testing.T) {
	// depth 1，val type int
	ret, _ := mapx.GetItems(deploySpec, []string{"replicas"})
	assert.Equal(t, 3, ret)

	// depth 2，val type map[string]any
	r, _ := mapx.GetItems(deploySpec, []string{"selector", "matchLabels"})
	_, ok := r.(map[string]any)
	assert.Equal(t, true, ok)

	// depth 2, val type string
	ret, _ = mapx.GetItems(deploySpec, []string{"strategy", "type"})
	assert.Equal(t, "RollingUpdate", ret)

	// depth 3, val type nil
	ret, _ = mapx.GetItems(deploySpec, []string{"template", "metadata", "creationTimestamp"})
	assert.Nil(t, ret)

	// depth 3, val type string
	ret, _ = mapx.GetItems(deploySpec, []string{"template", "spec", "restartPolicy"})
	assert.Equal(t, "Always", ret)
}

// paths 为 []string 或 其他，失败的情况
func TestGetItemsFailCase(t *testing.T) {
	// invalid paths type error
	_, err := mapx.GetItems(deploySpec, 0)
	assert.True(t, errors.Is(err, mapx.ErrInvalidPathType))

	// not paths error
	_, err = mapx.GetItems(deploySpec, []string{})
	assert.NotNil(t, err)

	// not map[string]any type error
	_, err = mapx.GetItems(deploySpec, []string{"replicas", "testKey"})
	assert.NotNil(t, err)

	_, err = mapx.GetItems(deploySpec, []string{"template", "spec", "containers", "image"})
	assert.NotNil(t, err)

	// key not exist
	_, err = mapx.GetItems(deploySpec, []string{"templateKey", "spec"})
	assert.NotNil(t, err)

	_, err = mapx.GetItems(deploySpec, []string{"selector", "spec"})
	assert.NotNil(t, err)

	// paths type error
	_, err = mapx.GetItems(deploySpec, []int{123, 456})
	assert.NotNil(t, err)

	_, err = mapx.GetItems(deploySpec, 123)
	assert.NotNil(t, err)
}

func TestGet(t *testing.T) {
	ret := mapx.Get(deploySpec, []string{"replicas"}, 1)
	assert.Equal(t, 3, ret)

	ret = mapx.Get(deploySpec, []string{}, nil)
	assert.Nil(t, ret)

	ret = mapx.Get(deploySpec, "container.name", "defaultName")
	assert.Equal(t, "defaultName", ret)
}

func TestGetBool(t *testing.T) {
	assert.True(t, mapx.GetBool(deploySpec, "template.spec.boolKey4GetBool"))
	assert.False(t, mapx.GetBool(deploySpec, "template.spec.notExistsKey"))
}

func TestGetInt64(t *testing.T) {
	assert.Equal(t, int64(10), mapx.GetInt64(deploySpec, "template.metadata.int64Key4GetInt64"))
	assert.Equal(t, int64(0), mapx.GetInt64(deploySpec, "template.spec.notExistsKey"))
}

func TestGetStr(t *testing.T) {
	assert.Equal(t, "value", mapx.GetStr(deploySpec, "template.metadata.labels.strKey4GetStr"))
	assert.Equal(t, "default-scheduler", mapx.GetStr(deploySpec, "template.spec.schedulerName"))
	assert.Equal(t, "", mapx.GetStr(deploySpec, "template.spec.notExistsKey"))
}

func TestGetList(t *testing.T) {
	assert.Equal(
		t, []any{map[string]any{"key": "value"}, "key-value"},
		mapx.GetList(deploySpec, "template.interfaceList"),
	)
	assert.Equal(t, []any{}, mapx.GetList(deploySpec, "template.spec.notExistsKey"))
}

func TestGetMap(t *testing.T) {
	assert.Equal(t, map[string]any{"app": "nginx"}, mapx.GetMap(deploySpec, "selector.matchLabels"))
	assert.Equal(t, map[string]any{}, mapx.GetMap(deploySpec, "template.spec.notExistsKey"))
}
