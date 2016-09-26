package fs

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"
	"time"
)

type Foo struct {
	BaseNode
}

type BaseInfo struct {
	Name  string
	Size  int64
	IsDir bool
	MTime time.Time
}

func NewBaseNode(name string) *BaseNode {
	return &BaseNode{
		Children: map[string]Node{},
		BaseInfo: BaseInfo{Name: name},
	}
}

type BaseNode struct {
	Children map[string]Node
	BaseInfo
}

func (b *BaseNode) childmap() map[string]Node {
	return b.Children
}

func (n *BaseNode) get(path string, mkdirp bool) (node Node, parent Node) {
	parts := strings.Split(strings.TrimPrefix(path, "/"), "/")
	if len(parts) == 0 {
		return nil, nil
	}
	//find/initialise all parent nodes
	parent = n
	name := parts[len(parts)-1]
	parents := parts[:len(parts)-1]
	for _, pname := range parents {
		m := childmap(parent)
		p, ok := m[pname]
		if ok {
			parent = p
		} else if mkdirp {
			//create missing parent
			b := NewBaseNode(pname)
			b.BaseInfo.IsDir = true
			parent.Upsert(pname, b)
			parent = b
		} else {
			return nil, nil
		}
	}
	m := childmap(parent)
	//get child and parent node
	c := m[name]
	if c != nil {
	} else {
	}
	return c, parent
}

func (n *BaseNode) Get(path string) Node {
	node, _ := n.get(path, false)
	return node
}

func (n *BaseNode) Upsert(path string, child Node) bool {
	existing, parent := n.get(path, true)
	if parent == child {
		panic("cannot insert node into itself")
	}
	m := childmap(parent)
	m[child.Name()] = child
	return reflect.DeepEqual(existing, child)
}

func (n *BaseNode) Delete(path string) bool {
	existing, parent := n.get(path, false)
	if existing == nil {
		return false
	}
	m := childmap(parent)
	delete(m, existing.Name())
	return true
}

func (n *BaseNode) GetChildren() []Node {
	nodes := make([]Node, len(n.Children))
	i := 0
	for _, node := range n.Children {
		nodes[i] = node
		i++
	}
	return nodes
}

func (n *BaseNode) ToMap() map[string]interface{} {
	m := map[string]interface{}{
		"Name": n.BaseInfo.Name,
	}
	if n.BaseInfo.IsDir {
		m["IsDir"] = true
	}
	if n.BaseInfo.Size > 0 {
		m["Size"] = n.BaseInfo.Size
	}
	if !n.MTime.IsZero() {
		m["MTime"] = n.BaseInfo.MTime
	}
	if c := n.GetChildren(); len(c) > 0 {
		m["Children"] = c
	}
	return m
}

//MarshalJSON is split into ToMap to allow classes
//to embed BaseNode, then augment the ToMap() output
//with class specifics
func (n *BaseNode) MarshalJSON() ([]byte, error) {
	b, err := json.Marshal(n.ToMap())
	if err != nil {
		return nil, err
	}
	return b, err
}

// Name of the file.
func (b *BaseNode) Name() string {
	return b.BaseInfo.Name
}

// Size of the file.
func (b *BaseNode) Size() int64 {
	return b.BaseInfo.Size
}

// IsDir returns true if the file is a directory.
func (b *BaseNode) IsDir() bool {
	return b.BaseInfo.IsDir
}

// TODO Sys is not implemented.
func (b *BaseNode) Sys() interface{} {
	return nil
}

// ModTime returns the modification time.
func (b *BaseNode) ModTime() time.Time {
	return b.BaseInfo.MTime
}

// Mode returns the file mode flags.
func (b *BaseNode) Mode() os.FileMode {
	var m os.FileMode
	if b.BaseInfo.IsDir {
		m |= os.ModeDir
	}
	return m
}

func (n *BaseNode) print() {
	cs := n.GetChildren()
	fmt.Printf("%s (#%d)\n", n.Name(), len(cs))
	for _, c := range cs {
		fmt.Printf("%s-", n.Name())
		c.(*BaseNode).print()
	}
}
