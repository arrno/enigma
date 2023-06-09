# Enigma
[![Go Reference](https://pkg.go.dev/badge/github.com/arrno/enigma.svg)](https://pkg.go.dev/github.com/arrno/enigma)
[![Go Build](https://github.com/arrno/enigma/actions/workflows/go.yml/badge.svg)](https://github.com/arrno/enigma/actions/workflows/go.yml)


Enigma provides a very simple API for doing recursive queries and insertions into data of unknown type, structure, and depth. This package is currently Alpha/Experimental.

## Install
```
go get github.com/arrno/enigma
```

## Initialize
Load some data into a new Enigma to start using the API. The data can contain any combination of maps, slices, and structs. The root must be a pointer.
```Go
data := map[string]any{
    "foo": "bar",
    "fiz": []int{0, 1, 2},
}

enigma := NewEnigma(&data)
```

## Paths and nodes
Paths are strings that represent a sequence of nodes to a target value. Each node in the path is separated by a dot. A node can represent a map key, a slice index, or a struct field depending on the runtime context.
```Go
path := "foo.bar.5.buz"
```

## Query by value
Query for all paths that end at some arbitrary value.
```Go
paths, err := enigma.QueryValue("biz")

for _, path := range paths {
    fmt.Println(path)
}
```

## Query by key
Query for all paths ending at some key.
```Go
paths, err := enigma.QueryKey("foo")

for _, path := range paths {
    fmt.Println(path)
}
```

## Query by path
Query for the value located at a target path.
```Go
val, err := enigma.QueryPath("foo.3.bar")

fmt.Println(val)
```

## Insert by path
Insert a value into a target path. **Insertions on struct fields will only work if the target is a public field on a struct pointer. The inserted value must also match the type definition of the parent node per usual.**
```Go
err := enigma.InsertByPath("fiz.buzz", 7)
```

## Find and replace by value
Replace all instances of one value with another value.
```Go
err := enigma.InsertByValue("original", "changed")
```

## Find and update by key
All instances of the target key will have their corresponding values updated.
```Go
err := enigma.InsertByKey("daysSinceToday", 0)
```

## Drop by path
The last node in the provided path will be pruned from the data.
```Go
err := enigma.DropByPath("col.items.5")
```
## Drop by value
Every parent node that holds the target value will be pruned from the data.
```Go
// Drop anything that holds a nil value.
err := enigma.DropByValue(nil)
```
## Get the data back
Display the underlying data in pretty format.
```Go
enigma.Display()
```
Return the underlying data back from Enigma.
```Go
data := enigma.Get()
```