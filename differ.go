package mapx

import (
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

const (
	// ActionAdd means has added
	ActionAdd = "Add"
	// ActionChange means has changed
	ActionChange = "Change"
	// ActionRemove means has removed
	ActionRemove = "Remove"
)

const (
	// NodeTypeIndex is index type node
	NodeTypeIndex = "index"
	// NodeTypeKey is key type node
	NodeTypeKey = "key"
)

// Node Map tree node
type Node struct {
	Type  string
	Index string
	Key   string
}

// NewIdxNode ...
func NewIdxNode(index int) Node {
	return Node{Type: NodeTypeIndex, Index: strconv.Itoa(index)}
}

// NewKeyNode ...
func NewKeyNode(key string) Node {
	return Node{Type: NodeTypeKey, Key: key}
}

// DiffRet Diff Result
type DiffRet struct {
	Action string
	Dotted string
	OldVal any
	NewVal any
}

// NewDiffRet ...
func NewDiffRet(action string, nodes []Node, old, new any) DiffRet {
	var b strings.Builder
	for _, n := range nodes {
		switch n.Type {
		case NodeTypeKey:
			// if key contains ".", add "(" and ")" for distinguish
			if strings.Contains(n.Key, ".") {
				b.WriteString(".(" + n.Key + ")")
			} else {
				b.WriteString("." + n.Key)
			}
		case NodeTypeIndex:
			b.WriteString("[" + n.Index + "]")
		}
	}
	dotted := strings.Trim(b.String(), ".")
	return DiffRet{Action: action, Dotted: dotted, OldVal: old, NewVal: new}
}

// String change DiffRet to string
func (r *DiffRet) String() string {
	ret := fmt.Sprintf("%s %s: ", r.Action, r.Dotted)
	switch r.Action {
	case ActionAdd:
		ret += fmt.Sprintf("%v", r.NewVal)
	case ActionChange:
		ret += fmt.Sprintf("%v -> %v", r.OldVal, r.NewVal)
	case ActionRemove:
		ret += fmt.Sprintf("%v", r.OldVal)
	}
	return ret
}

// DiffRetList Diff Result List
type DiffRetList []DiffRet

// Len ...
func (l DiffRetList) Len() int {
	return len(l)
}

// Less sort by type and dotted
func (l DiffRetList) Less(i, j int) bool {
	cmpRet := strings.Compare(l[i].Action, l[j].Action)
	if cmpRet == 0 {
		return strings.Compare(l[i].Dotted, l[j].Dotted) < 0
	}
	return cmpRet < 0
}

// Swap ...
func (l DiffRetList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

// Differ is Map differ, used to compare the keys and values of two maps.
type Differ struct {
	old  map[string]any
	new  map[string]any
	rets DiffRetList
}

// NewDiffer ...
func NewDiffer(old, new map[string]any) *Differ {
	return &Differ{old: old, new: new, rets: DiffRetList{}}
}

// Do execute compare
func (d *Differ) Do() DiffRetList {
	d.handleMap(d.old, d.new, []Node{})
	sort.Sort(d.rets)
	return d.rets
}

func (d *Differ) handleMap(old, new map[string]any, nodes []Node) {
	intersection, addition, deletion := []string{}, []string{}, []string{}
	for key := range old {
		if Exists(new, key) {
			intersection = append(intersection, key)
		}
	}
	for key := range new {
		if !Exists(old, key) {
			addition = append(addition, key)
		}
	}
	for key := range old {
		if !Exists(new, key) {
			deletion = append(deletion, key)
		}
	}

	// intersection
	for _, key := range intersection {
		curNodes := append(nodes, NewKeyNode(key))
		d.handle(old[key], new[key], curNodes)
	}

	// addition
	for _, key := range addition {
		ret := NewDiffRet(ActionAdd, append(nodes, NewKeyNode(key)), nil, new[key])
		d.rets = append(d.rets, ret)
	}

	// deletion
	for _, key := range deletion {
		ret := NewDiffRet(ActionRemove, append(nodes, NewKeyNode(key)), old[key], nil)
		d.rets = append(d.rets, ret)
	}
}

func (d *Differ) handleList(old, new []any, nodes []Node) {
	oldLen, newLen, minLen := len(old), len(new), len(old)
	if newLen < oldLen {
		minLen = newLen
	}

	// intersection
	for idx := 0; idx < minLen; idx++ {
		d.handle(old[idx], new[idx], append(nodes, NewIdxNode(idx)))
	}

	// addition
	for idx := minLen; idx < newLen; idx++ {
		ret := NewDiffRet(ActionAdd, append(nodes, NewIdxNode(idx)), nil, new[idx])
		d.rets = append(d.rets, ret)
	}

	// deletion
	for idx := minLen; idx < oldLen; idx++ {
		ret := NewDiffRet(ActionRemove, append(nodes, NewIdxNode(idx)), old[idx], nil)
		d.rets = append(d.rets, ret)
	}
}

func (d *Differ) handle(old, new any, nodes []Node) {
	oldMap, oldIsMap := old.(map[string]any)
	newMap, newIsMap := new.(map[string]any)
	if oldIsMap && newIsMap {
		d.handleMap(oldMap, newMap, nodes)
		return
	}

	oldList, oldIsList := old.([]any)
	newList, newIsList := new.([]any)
	if oldIsList && newIsList {
		d.handleList(oldList, newList, nodes)
		return
	}

	if !reflect.DeepEqual(old, new) {
		ret := NewDiffRet(ActionChange, nodes, old, new)
		d.rets = append(d.rets, ret)
	}
}
