package disk

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"

	"github.com/jpillora/cloud-torrent/cloudtorrent/fs"
)

const fileNumberLimit = 100

//custom directory walk
func walk(node *fs.BaseNode, path string, info os.FileInfo, n *int64) error {
	if (!info.IsDir() && !info.Mode().IsRegular()) || strings.HasPrefix(info.Name(), ".") {
		return errors.New("Non-regular file")
	}
	if atomic.AddInt64(n, 1) > fileNumberLimit {
		return errors.New("Over file limit") //limit number of files walked
	}
	node.BaseInfo.Name = info.Name()
	node.BaseInfo.Size = info.Size()
	node.BaseInfo.MTime = info.ModTime()
	if info.IsDir() {
		children, err := ioutil.ReadDir(path)
		if err != nil {
			return fmt.Errorf("Failed to list files")
		}
		node.BaseInfo.Size = 0
		node.Children = map[string]fs.Node{}
		for _, i := range children {
			c := &fs.BaseNode{}
			p := filepath.Join(path, i.Name())
			if err := walk(c, p, i, n); err != nil {
				continue
			}
			node.BaseInfo.Size += c.BaseInfo.Size
			node.Children[c.BaseInfo.Name] = c
		}
	}
	return nil
}
