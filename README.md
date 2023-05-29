# enigma
Enigma is a utility for handling query and insert operations into data structures of unknown types and depths.

## Install
```
go get https://github.com/arrno/enigma
```

## Paths and Nodes
Paths are strings that represent a sequence of nodes to a target. Each node in the path is separated by a dot. A node can be a map key, a slice index, or a struct field. **At this time, only slices and maps with string keys are supported node types.**
```Go
path := "foo.bar.5.buz"
```
## Query by value
Query for all paths to some target value of any type.

## Query by path
Query for the value located at a target path.

## Insert by path
Inert a value at the target path.

## Insert by value
Replace all instances of one value with another value.