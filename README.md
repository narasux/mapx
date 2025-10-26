# mapx

[中文](./README_ZH.md)

> mapx is a Golang map utility package that includes some common map usage shortcuts.

## Usage

You can use the mapx package in the following way:

```go
import "github.com/narasux/mapx"
```

## Spec

### Exists

Check if the key exists in the map.

```go
m := map[string]any{"k1": "v1"}

// true
mapx.Exists(m, "k1")

// false
mapx.Exists(m, "k2")
```

### Differ

Compare two maps and output the differences, supporting nested item comparison.

**NOTE**
- Nested item comparison is only supported for `[]any` and `map[string]any`.
- If a key contains `.`, parentheses will be added in the output (e.g. `[]string{"k1", "k2.2", "k3"} → "k1.(k2.2).k3"`)
- For complete test cases, see [differ_test.go](./differ_test.go)

```go
o := map[string]any{
    "k1": "v1", 
    "k2": "v2",
    "k3": map[string]any{
        "k4": "v4",
    },
}

n := map[string]any{
    "k1": "v1.1", 
    "k3": map[string]any{
        "k4": "v4.1",
    }, 
    "k5": "v5",
}

/* 
[
    {"Add",     "k5",     <nil>,  "v5"}
    {"Change",  "k1",     "v1",   "v1.1"}
    {"Change",  "k3.k4",  "v4",   "v4.1"}
    {"Remove",  "k2",     "v2",   <nil>}
]
*/
diffRets := mapx.NewDiffer(o, n).Do()

for _, r := range diffRets {
    /*
        Add k5: v5
        Change k1: v1 -> v1.1
        Change k3.k4: v4 -> v4.1
        Remove k2: v2
     */
    s := r.String()
}
```

### GetItems

A method to get the value from the nested `map[string]any` according to the specified path.

```go
m := map[string]any{
    "a1": map[string]any{
        "b1": map[string]any{
            "c1": map[string]any{
                "d1": "v1", 
                "d2": "v2", 
                "d.3": 3,
            },
        },
    },
}

// d1val: v1
d1Val, _ := mapx.GetItems(m, "a1.b1.c1.d1")

// if any key contains `.` in path, you can use []string as parameter
// dDot3Val: 3
dDot3Val, _ := mapx.GetItems(m, []string{"a1", "b1", "c1", "d.3"})

// any key not exist or intermediate value not map[string]any type, return error
// err: key c2 not exist
_, err := mapx.GetItems(m, "a1.b1.c2")
```

### Get

A shortcut method for `GetItems` that supports setting a default value. 

When the original `GetItems` returns an error (`err != nil`), this shortcut returns the default value.

```go
m := ...

// d1val: v1
d1Val := mapx.Get(m, "a1.b1.c1.d1", "default")

// c2Val: default
c2Val := mapx.GetItems(m, "a1.b1.c2", "default")
```

### GetBool

A shortcut method for `Get`, with a default return value of `false`.

### GetInt64

A shortcut method for `Get`, with a default return value of `int64(0)`.

### GetFloat64

A shortcut method for `Get`, with a default return value of `float64(0)`.

### GetStr

A shortcut method for `Get`, with a default return value of `""` (empty string).

### GetList

A shortcut method for `Get`, with a default return value of `[]any{}` (empty list).

### GetMap

A shortcut method for `Get`, with a default return value of `map[string]any{}` (empty map).

### SetItems

A method to set values for nested `map[string]any` structures by specifying a path.

```go
m := map[string]any{
    "a1": map[string]any{
        "b1": map[string]any{
            "c1": []any{
                "d1", "d2", "d3",
            },
        },
    },
}

/*
   m = map[string]any{
        "a1": map[string]any{
            "b1": map[string]any{
                "c1": []any{
                    "d1", "d2", "d3",
                },
                "c2": "d4",
            },
        },
    }
 */
_ = mapx.SetItems(m, "a1.b1.c2", "d4")

// error will be returned when an intermediate value does not exist 
// or its value is not map[string]any
// err: key c1 not exists or obj[key] not map[string]any type
err := mapx.SetItems(m, "a1.b1.c1.d1", "d5")
```
