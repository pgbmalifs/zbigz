package torrent

import (
	"path/filepath"

	"github.com/anacrolix/torrent"
	"github.com/jpillora/cloud-torrent/cloudtorrent/fs"
)

type torrentNode struct {
	fs.BaseNode
	rx, tx int64
}

func torrentToNode(t *torrent.Torrent) *torrentNode {
	n := &torrentNode{
		BaseNode: fs.BaseNode{
			BaseInfo: fs.BaseInfo{
				Name:  t.Name(),
				IsDir: true,
			},
		},
	}
	for _, f := range t.Files() {
		fnode := fileToNode(&f)
		n.Upsert(f.Path(), fnode)
	}
	return n
}

type fileNode struct {
	fs.BaseNode
	peers int
}

func fileToNode(f *torrent.File) fs.Node {
	n := &fileNode{
		BaseNode: fs.BaseNode{
			BaseInfo: fs.BaseInfo{
				Name: filepath.Base(f.Path()),
				Size: f.Length(),
			},
		},
	}
	return n
}
