package mapx_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/narasux/mapx"
)

var oldMap = map[string]any{
	"a1": map[string]any{
		"b1": map[string]any{
			"c1": map[string]any{
				"d1": "v1",
				"d2": "v2",
				"d3": 3,
				"d4": []any{4, 5},
				"d5": nil,
				"d6": []any{
					6.1, 6.2, 6.3, 6.4, 6.5,
				},
			},
		},
	},
	"a2": []any{
		map[string]any{
			"b2": map[string]any{
				"c2": []any{
					"d1",
					map[string]any{
						"e1": "v1",
						"e2": "v2",
					},
					map[string]string{
						"e3": "v3",
						"e4": "v4",
					},
					2,
				},
			},
			"b3": []any{
				"c3", "c4", 5,
			},
		},
	},
}

var newMap = map[string]any{
	"a1": map[string]any{
		"b1": map[string]any{
			"c1": map[string]any{
				"d1": "v1",
				// change a1.b1.c1.d2 v2-v1
				"d2": "v1",
				// remove a1.b1.c1.d3 ...
				// add a1.b1.c1.d7 ...
				"d7": 3,
				// remove a1.b1.c1.d4[1] ...
				"d4": []any{4},
				// change a1.b1.c1.d5 nil->"nil"
				"d5": "nil",
				// change a1.b1.c1.d6[2] 6.3->6.4
				// change a1.b1.c1.d6[3] 6.4->6.5
				// change a1.b1.c1.d6[4] 6.5->6.3
				"d6": []any{
					6.1, 6.2, 6.4, 6.5, 6.3,
				},
			},
		},
	},
	"a2": []any{
		map[string]any{
			"b2": map[string]any{
				"c2": []any{
					// change a2[0].b2.c2[0] d1->d2
					"d2",
					map[string]any{
						// change a2[0].b2.c2[1].e1 v1->v2
						"e1": "v2",
						// remove a2[0].b2.c2[1].e2 ...
						// add a2[0].b2.c2[1].e3 ...
						"e3": "v2",
						// add a2[0].b2.c2[1].e4 ...
						"e4": "v4",
						// add a2[0].b2.c2[1].(e5.f1) ...
						"e5.f1": "v5",
					},
					// change a2[0].b2.c2[2] ...
					map[string]string{
						"e3": "v4", // only v3->v4, but map[string]string will not expand for compare
						"e4": "v4",
					},
					// change a2[0].b2.c2[3] 2->1
					1,
					// add a2[0].b2.c2[4] 2
					2,
				},
			},
			// change a2[0].b3[0] "c3"->"c4"
			// change a2[0].b3[2] 5->6
			// add a2[0].b3[3] 7
			"b3": []any{
				"c4", "c4", 6, 7,
			},
		},
	},
	// add a3 ...
	"a3": map[string]any{
		"b4": "v1",
	},
}

var exceptedDiffRets = mapx.DiffRetList{
	{mapx.ActionAdd, "a1.b1.c1.d7", nil, 3},
	{mapx.ActionAdd, "a2[0].b2.c2[1].(e5.f1)", nil, "v5"},
	{mapx.ActionAdd, "a2[0].b2.c2[1].e3", nil, "v2"},
	{mapx.ActionAdd, "a2[0].b2.c2[1].e4", nil, "v4"},
	{mapx.ActionAdd, "a2[0].b2.c2[4]", nil, 2},
	{mapx.ActionAdd, "a2[0].b3[3]", nil, 7},
	{mapx.ActionAdd, "a3", nil, map[string]any{"b4": "v1"}},
	{mapx.ActionChange, "a1.b1.c1.d2", "v2", "v1"},
	{mapx.ActionChange, "a1.b1.c1.d5", nil, "nil"},
	{mapx.ActionChange, "a1.b1.c1.d6[2]", 6.3, 6.4},
	{mapx.ActionChange, "a1.b1.c1.d6[3]", 6.4, 6.5},
	{mapx.ActionChange, "a1.b1.c1.d6[4]", 6.5, 6.3},
	{mapx.ActionChange, "a2[0].b2.c2[0]", "d1", "d2"},
	{mapx.ActionChange, "a2[0].b2.c2[1].e1", "v1", "v2"},
	{
		mapx.ActionChange,
		"a2[0].b2.c2[2]",
		map[string]string{"e3": "v3", "e4": "v4"},
		map[string]string{"e3": "v4", "e4": "v4"},
	},
	{mapx.ActionChange, "a2[0].b2.c2[3]", 2, 1},
	{mapx.ActionChange, "a2[0].b3[0]", "c3", "c4"},
	{mapx.ActionChange, "a2[0].b3[2]", 5, 6},
	{mapx.ActionRemove, "a1.b1.c1.d3", 3, nil},
	{mapx.ActionRemove, "a1.b1.c1.d4[1]", 5, nil},
	{mapx.ActionRemove, "a2[0].b2.c2[1].e2", "v2", nil},
}

func TestDiffer(t *testing.T) {
	diffRets := mapx.NewDiffer(oldMap, newMap).Do()
	assert.Equal(t, exceptedDiffRets, diffRets)
}

func TestDiffRetString(t *testing.T) {
	addDiffRet := mapx.DiffRet{Action: mapx.ActionAdd, Dotted: "a1.b1.c1.d7", NewVal: 3}
	assert.Equal(t, "Add a1.b1.c1.d7: 3", addDiffRet.String())

	changeDiffRet := mapx.DiffRet{Action: mapx.ActionChange, Dotted: "a1.b1.c1.d5", NewVal: "nil"}
	assert.Equal(t, "Change a1.b1.c1.d5: <nil> -> nil", changeDiffRet.String())

	changeDiffRet = mapx.DiffRet{Action: mapx.ActionChange, Dotted: "a1.b1.c1.d2", OldVal: "v2", NewVal: "v1"}
	assert.Equal(t, "Change a1.b1.c1.d2: v2 -> v1", changeDiffRet.String())

	removeDiffRet := mapx.DiffRet{Action: mapx.ActionRemove, Dotted: "a1.b1.c1.d4[1]", OldVal: 5}
	assert.Equal(t, "Remove a1.b1.c1.d4[1]: 5", removeDiffRet.String())
}
