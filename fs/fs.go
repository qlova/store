package fs

import (
	"io"
	"net/http"
	"os"
	"strings"
)

//MetaData returned by File.Stat
type MetaData os.FileInfo

//Index is an index into a list returned by Node.List
type Index struct {
	Path   Path
	Int    int
	String string
}

//Child is a node's child.
type Child interface {
	//Get the current child's index.
	Index() Index

	//Return the current child as a node, nil if the current child is not a node.
	Node() Node

	//Return the current child as data, nil if the current child is not data.
	Data() Data
}

//Data is the lowlevel interface required for Files.
type Data interface {
	io.ReaderFrom
	io.WriterTo
	http.Handler

	//Return the path of this data.
	Path() Path

	//Stat returns info about the data.
	Stat() (os.FileInfo, error)

	//Deletes the data.
	Delete() error
}

//File contains Data
type File struct {
	Data
}

//SetString sets the contents of the file to the given string.
func (f File) SetString(contents string) error {
	_, err := f.ReadFrom(strings.NewReader(contents))
	return err
}

//String returns the contents of the file as a string.
func (f File) String() string {
	var builder strings.Builder
	f.WriteTo(&builder)
	return builder.String()
}

//Node is the lowlevel interface required for Directories.
//A node is a location in a tree, it doesn't have to exist.
type Node interface {
	//Create this node and all parent nodes required for this node to exist.
	Create() error

	//Delete this node and all children nodes and data.
	Delete() error

	//Return the node at the relative path.
	Goto(relative Path) Node

	//Return the data of this node at the relative path.
	Data(name string) Data

	//Return the path of this node.
	Path() Path

	//Stat returns info about the directory.
	Stat() (os.FileInfo, error)

	//Slice returns a slice of children.
	Slice(offset Index, length int) ([]Child, error)
}

//Directory can contain files or other directories.
type Directory struct {
	Node
}

//File returns the file inside the given directory with the given name.
func (d Directory) File(name string) File {
	return File{d.Data(name)}
}

//Goto goes to the given path.
func (d Directory) Goto(path Path) Directory {
	return Directory{d.Node.Goto(path)}
}

//Root is a file-system root.
type Root struct {
	Directory
}

//NewRoot returns a new Root from a given Node.
func NewRoot(node Node) Root {
	var root Root
	root.Node = node
	return root
}
