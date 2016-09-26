package diskModule

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/jpillora/cloud-torrent/cloudtorrent/fs"
	"github.com/jpillora/cloud-torrent/cloudtorrent/module"
	"github.com/jpillora/filenotify"
	homedir "github.com/mitchellh/go-homedir"
)

//============================

type factory struct{}

func (f *factory) Type() string {
	return "disk"
}

func (f *factory) New(id string) module.Module {
	return &diskModule{
		id: id,
	}
}

//============================

type diskModule struct {
	id      string
	watcher filenotify.FileWatcher
	config  struct {
		Base string
	}
}

func (d *diskModule) ID() string {
	return "disk:" + d.id
}

func (d *diskModule) Mode() fs.FSMode {
	return fs.RW
}

func (d *diskModule) Configure(raw json.RawMessage) (interface{}, error) {
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
	//ready!
	d.config.Base = base
	return &d.config, nil
}

func (d *diskModule) Sync(chan fs.Node) error {
	d.watcher = filenotify.New()
	defer d.watcher.Close()
	//set poll interval (if polling is being used)
	filenotify.SetPollInterval(d.watcher, time.Second)
	d.watcher.Add(d.config.Base)
	for event := range d.watcher.Events() {
		log.Printf("event %+v", event)
	}
	return nil
}

func logf(format string, args ...interface{}) {
	log.Printf("[Disk] "+format, args...)
}
