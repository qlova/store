package store

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"
)

//Reference is for keeping track of children.
type Reference struct {

	//The package this Reference has been set from.
	Package string

	//The internal package location of the Reference.
	Internal string

	//This should be set to the absolute path of the node in case the implementation changes.
	Fallback string
}

//Children is a list of children for a given node.
type Children interface {

	//Get a reference to this child.
	Reference() Reference

	//Skip to the child with name 'name'.
	SkipTo(reference Reference) error

	//Return the name of the current child.
	Name() string

	//Return the current child as a node, nil if the current child is not a node.
	Node() Node

	//Return the current child as data, nil if the current child is not data.
	Data() Data

	//Goto the next child.
	Next() error
}

//Tree is the lowlevel interface required for Objects.
type Tree interface {
	Node
}

//Node is the lowlevel interface required for Trees.
//A node is a location in a tree, it doesn't have to exist.
type Node interface {

	//Return the name of this node.
	Name() string

	//Create this node.
	Create() error

	//Delete this node.
	Delete() error

	//Return the node at the relative path 'location'.
	Goto(location string) Node

	//Return the data child of this node with the specified 'name'.
	Data(name string) Data

	//Return the absolute path of this node.
	Path() string

	//Returns true if the node is currently available to make modifications.
	Available() bool

	Children(max ...int) Children

	//Returns the parent node.
	Parent() Node
}

//Data is the lowlevel interface required for Values.
type Data interface {

	//Sets data from reader.
	From(io.Reader) error

	//Copy data to writer.
	CopyTo(io.Writer) error

	//Deletes data.
	Delete() error

	//Returns the size of the data, -1 if unknown size.
	Size() int64

	//Return the absolute path of this data.
	Path() string

	//Returns true if data is available to read.
	Available() bool
}

//Root can contain objects, or values.
type Root struct {
	Object
}

//Object can contain other objects, or values.
type Object struct {
	Node
}

//Value returns the value of name 'name'.
func (object Object) Value(name string) Value {
	return Value{object.Data(name)}
}

//Goto returns a new node at the location relative to the current node.
func (object Object) Goto(location string) Object {
	return Object{object.Node.Goto(location)}
}

//List returns a slice of values and objects for this object.
func (object Object) List(amount ...int) (result []string, err error) {

	var children = object.Children(amount...)

	for err := children.Next(); err == nil; err = children.Next() {
		result = append(result, children.Name())
	}
	if err == io.EOF {
		err = nil
	}

	return
}

//Value is an abstract representation of a datatype.
type Value struct {
	Data
}

//SetString sets value to s
func (value Value) SetString(s string) error {
	return value.From(strings.NewReader(s))
}

func (value Value) String() (string, error) {
	var buffer bytes.Buffer
	err := value.CopyTo(&buffer)
	return buffer.String(), err
}

//EncodeIndentedJSON encodes indented JSON onto the value.
func (value Value) EncodeIndentedJSON(i interface{}) (err error) {
	var reader, writer = io.Pipe()
	go func() {
		var encoder = json.NewEncoder(writer)
		encoder.SetIndent("", "\t")
		err = encoder.Encode(i)
		writer.Close()
	}()
	if err2 := value.From(reader); err2 != nil {
		return err2
	}
	return err
}
