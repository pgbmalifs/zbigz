package torrent

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"time"

	"github.com/jpillora/cloud-torrent/cloudtorrent/fs"
	"github.com/jpillora/cloud-torrent/cloudtorrent/module"
	"github.com/spf13/afero"
)

func New() module.Module {
	return &torrentModule{}
}

type torrentModule struct {
	config struct {
		PeerID            string
		DownloadDirectory string
		EnableUpload      bool
		EnableSeeding     bool
		EnableEncryption  bool
		AutoStart         bool
		IncomingPort      int
	}
}

func (t *torrentModule) TypeID() string {
	return "torrent" //singleton
}

func (t *torrentModule) Mode() fs.FSMode {
	return fs.RW
}

func (t *torrentModule) Configure(raw json.RawMessage) (interface{}, error) {
	if err := json.Unmarshal(raw, &t.config); err != nil {
		return nil, err
	}
	unset := t.config.PeerID == "" && t.config.IncomingPort == 0
	if t.config.PeerID == "" {
		t.config.PeerID = "Cloud Torrent"
	}
	if t.config.IncomingPort == 0 {
		t.config.IncomingPort = 4479
	}
	if unset {
		t.config.EnableEncryption = true
		t.config.EnableSeeding = true
		t.config.EnableUpload = true
	}
	return &t.config, nil
}

func (t *torrentModule) Update(chan fs.Node) error {
	return nil
}

func (t *torrentModule) Create(name string) (afero.File, error) {
	return &file{}, nil
}

func (t *torrentModule) Open(name string) (afero.File, error) {
	return &file{}, nil
}

func (t *torrentModule) OpenFile(name string, flag int, perm os.FileMode) (afero.File, error) {
	return t.Open(name)
}

func (t *torrentModule) Mkdir(name string, perm os.FileMode) error {
	return errors.New("not supported yet")
}

func (t *torrentModule) MkdirAll(path string, perm os.FileMode) error {
	return errors.New("not supported yet")
}

func (t *torrentModule) Remove(name string) error {
	return errors.New("not supported yet")
}

func (t *torrentModule) RemoveAll(path string) error {
	return errors.New("not supported yet")
}

func (t *torrentModule) Rename(oldname, newname string) error {
	return errors.New("not supported yet")
}

func (t *torrentModule) Stat(name string) (os.FileInfo, error) {
	return nil, errors.New("not supported yet")
}

func (t *torrentModule) Chmod(name string, mode os.FileMode) error {
	return errors.New("not supported yet")
}

func (t *torrentModule) Chtimes(name string, atime time.Time, mtime time.Time) error {
	return errors.New("not supported yet")
}

func logf(format string, args ...interface{}) {
	log.Printf("[Torrents] "+format, args...)
}
