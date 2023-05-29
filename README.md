# Enigma
Enigma is a utility for handling query and insert operations into data structures of unknown types and depths. Do strange things like extend the data structure by inserting new arrays/maps.

## Install
```
go get https://github.com/arrno/enigma
```

## Initialize
Load some data into a new Enigma to start using the API operations.
```Go
data := map[string]any{
    "foo": "bar",
    "fiz": []int{0, 1, 2},
}

enigma, err := NewEnigma(data)
```

## Paths and nodes
Paths are strings that represent a sequence of nodes to a target. Each node in the path is separated by a dot. A node can be a map key, a slice index, or a struct field. **At this time, only slices and maps with string keys are supported node types.**
```Go
path := "foo.bar.5.buz"
```

## Query by value
Query for all paths to some target value of any type.
```Go
paths, err := enigma.QueryValue("biz")
```

## Query by path
Query for the value located at a target path.
```Go
val, err := enigma.QueryPath("foo.3.bar")
```

## Insert by path
Insert a value of any type at the target path.
```Go
err := enigma.InsertByPath("fiz.buzz", 7)
```

## Insert by value
Replace all instances of one value with another value.
```Go
err := enigma.InsertByValue("original", "changed")
```

## Get data
Get the underlying data back from Enigma.
```Go
data := enigma.Get()
```