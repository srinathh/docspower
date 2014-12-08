//package nodesprovider abstracts access to a data source for docs power
package nodesprovider

import (
	"errors"
	"os"
	"path/filepath"
	"sync"
	"time"
)

//the nodesprovider.Interface defines the methods that can be used to access nodes data
type Interface interface {
	GetNode(string id) (Node, error)   //returns the requested node special id of # returns root
	GetData(string id) ([]byte, error) //reads data and returns the bytes of a file
}

//Node data structure is used to represent data about a node. the Provider implementation sets
//an ID that can be used to access the node hiding information about where the item is stored
//the size of a node can be obtained by the Size() function which calculates size of a
//node hierarchy on the fly rather than by a data element
type Node interface {
	Name() string
	ModTime() time.Time
	Size() int64
	Id() string
	Children() []string
}

//fsnode is an implementation of Node for accessing file system nodes. It embeds os.FileInfo
//and contains a path
type fsNode struct {
	Fi         os.FileInfo
	Path       string
	ChildNodes []string
}

func (node *fsNode) Name() string {
	return node.Fi.Name()
}

func (node *fsNode) ModTime() time.Time {
	return node.Fi.ModTime()
}

func (node *fsNode) Size() int64 {
	size := node.Fi.Size()
	for _, child := range node.Children() {
		size = size + child.Size()
	}
	return size
}

func (n *fsNode) Id() string {
	hasher := fnv.New64()
	hasher.Write([]byte(n.Path))
	return hex.EncodeToString(hasher.Sum(nil))
}

func (node *fsNode) Children() []string {
	return node.ChildNodes
}
