package dropbox

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"time"

	"github.com/jpillora/cloud-torrent/cloudtorrent/fs"
	"github.com/jpillora/cloud-torrent/cloudtorrent/module"
	dropbox "github.com/tj/go-dropbox"
)

//============================

type factory struct{}

func (f *factory) Type() string {
	return "dropbox"
}

func (f *factory) New(id string) module.Module {
	return &dropboxModule{
		id: id,
		root: &file{
			BaseNode: fs.BaseNode{
				BaseInfo: fs.BaseInfo{
					Name:  "root",
					IsDir: true,
				},
			},
		},
	}
}

func init() {
	module.Register(&factory{})
}

//============================

type dropboxModule struct {
	id     string
	client *dropbox.Client
	config struct {
		Token string
		Base  string
	}
	root *file
}

func (d *dropboxModule) TypeID() string {
	return "dropbox:" + d.id
}

func (d *dropboxModule) Configure(raw json.RawMessage) (interface{}, error) {
	if err := json.Unmarshal(raw, &d.config); err != nil {
		return nil, err
	}
	if d.config.Token == "" {
		d.client = nil
		return nil, errors.New("API token missing")
	}
	if d.config.Base == "" {
		d.config.Base = string(os.PathSeparator)
	}
	d.client = dropbox.New(dropbox.NewConfig(d.config.Token))
	logf("configured")
	return &d.config, nil
}

func (d *dropboxModule) Sync(updates chan fs.Node) error {
	logf("update...")
	c := d.client
	if c == nil {
		return errors.New("API token was removed")
	}
	emit := true
	//list all files in base
	resp, err := c.Files.ListFolder(&dropbox.ListFolderInput{
		Path:      d.config.Base,
		Recursive: true,
	})
	if err != nil {
		return err
	}
	for {
		//copy results into memory
		for _, m := range resp.Entries {
			if d.updateFile(m) {
				emit = true
			}
		}
		//emit updates
		if !resp.HasMore && emit {
			logf("emit root")
			updates <- d.root
			emit = false
		}
		//poll next set
		resp, err = c.Files.ListFolderContinue(&dropbox.ListFolderContinueInput{
			Cursor: resp.Cursor,
		})
		if err != nil {
			return err
		}
		if !resp.HasMore {
			time.Sleep(3 * time.Second)
		}
	}
}

func (d *dropboxModule) updateFile(m *dropbox.Metadata) bool {
	//node path
	path := m.PathDisplay
	//deletion
	if m.Tag == "deleted" {
		return d.root.Delete(path)
	}
	//node
	f := &file{
		BaseNode: fs.BaseNode{
			BaseInfo: fs.BaseInfo{
				Name:  m.Name,
				Size:  int64(m.Size),
				IsDir: m.Tag == "folder",
				MTime: m.ServerModified,
			},
		},
	}
	logf("%+v", f)
	//insert
	return d.root.Upsert(path, f)
}

func logf(format string, args ...interface{}) {
	log.Printf("[Dropbox] "+format, args...)
}
