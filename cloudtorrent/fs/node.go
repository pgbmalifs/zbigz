package fs

import (
	"encoding/json"
	"os"
)

//requirements:
//1. edit node tree. insert. delete
//2. jsonify node tree
//3. extensible

type Node interface {
	Name() string
	Get(path string) Node
	GetChildren() []Node
	Upsert(path string, child Node) bool
	Delete(path string) bool
	json.Marshaler
}

func childmap(n Node) map[string]Node {
	mapper, ok := n.(interface {
		childmap() map[string]Node
	})
	if !ok {
		panic("node should implement childmap()")
	}
	return mapper.childmap()
}

func filename(path string) string {
	for i := len(path) - 2; i >= 0; i-- {
		if path[i] == os.PathSeparator {
			return string(path[i+1:])
		}
	}
	return path
}
