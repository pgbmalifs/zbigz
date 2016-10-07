package module

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/jpillora/backoff"
	"github.com/jpillora/cloud-torrent/cloudtorrent/fs"
	"github.com/jpillora/velox"
)

//shared module state
type State struct {
	sync.Locker  `json:"-"`
	velox.Pusher `json:"-"`
	TypeID       string
	Enabled      bool           `json:",omitempty"`
	Syncing      bool           `json:",omitempty"`
	Config       interface{}    `json:",omitempty"`
	Root         json.Marshaler `json:",omitempty"`
	Error        string         `json:",omitempty"`
}

//Sync runs once after the first
//successful configure, then loops the module's Sync()
//forever, with exponential backoff on failures.
func (s *State) Sync(m Module) {
	typeid := m.TypeID()
	if f, ok := m.(fs.FS); ok {
		s.syncFS(typeid, f)
	}
}

func (s *State) syncFS(typeid string, f fs.FS) {
	updates := make(chan fs.Node)
	//monitor sync updates
	go func() {
		for node := range updates {
			s.Lock()
			// log.Printf("[%s] fs root updated", typeid)
			s.Root = node
			s.Unlock()
			s.Push()
		}
	}()
	//sync loop forever
	go func() {
		b := backoff.Backoff{Max: 5 * time.Minute}
		for {
			//retrieve updates
			err := f.Sync(updates)
			e := ""
			d := 1 * time.Second
			if err == nil {
				b.Reset()
			} else {
				log.Printf("[%s] fs sync failed: %s", typeid, err)
				e = err.Error()
				d = b.Duration()
			}
			//show result
			s.Lock()
			s.Error = e
			s.Unlock()
			s.Push()
			//retry after sleep
			time.Sleep(d)
		}
	}()
	log.Printf("[%s] syncing fs", typeid)
}
