package disk

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/jpillora/cloud-torrent/cloudtorrent/fs"
	"github.com/jpillora/cloud-torrent/cloudtorrent/module"
	homedir "github.com/mitchellh/go-homedir"
)

//============================

type factory struct{}

func (f *factory) Type() string {
	return "disk"
}

func (f *factory) New(id string) module.Module {
	m := &diskModule{
		id: id,
	}
	return m
}

func init() {
	module.Register(&factory{})
}

//============================

type file struct {
	fs.BaseNode
}

//============================

type diskModule struct {
	id     string
	root   *file
	config struct {
		Base string
	}
}

func (d *diskModule) TypeID() string {
	return "disk:" + d.id
}

func (d *diskModule) Configure(raw json.RawMessage) (interface{}, error) {
	prevBase := d.config.Base
	if err := json.Unmarshal(raw, &d.config); err != nil {
		return nil, err
	}
	base := d.config.Base
	if base == "" {
		if hd, err := homedir.Dir(); err == nil {
			base = filepath.Join(hd, "downloads")
		} else if wd, err := os.Getwd(); err == nil {
			base = filepath.Join(wd, "downloads")
		} else {
			return nil, errors.New("Cannot find default base directory")
		}
	}
	info, err := os.Stat(base)
	if os.IsNotExist(err) {
		return nil, errors.New("Cannot find directory")
	} else if err != nil {
		return nil, err
	} else if !info.IsDir() {
		return nil, errors.New("Path is not a directory")
	}
	//base dir changed
	if prevBase != base {
		// d.logf("watching %s for file changes", base)
	}
	//ready!
	d.config.Base = base
	return &d.config, nil
}

func (d *diskModule) Sync(updates chan fs.Node) error {
	d.root = &file{
		BaseNode: fs.BaseNode{
			BaseInfo: fs.BaseInfo{
				Name:  "/",
				IsDir: true,
			},
		},
	}
	if info, err := os.Stat(d.config.Base); err == nil {
		n := int64(0)
		if err := walk(&d.root.BaseNode, d.config.Base, info, &n); err != nil {
			return err
		}
	}
	updates <- d.root
	time.Sleep(60 * time.Second)
	return nil
}

func (d *diskModule) logf(format string, args ...interface{}) {
	log.Printf("[Disk:"+d.id+"] "+format, args...)
}
