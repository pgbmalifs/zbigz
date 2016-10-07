package torrent

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/anacrolix/torrent"
	"github.com/jpillora/cloud-torrent/cloudtorrent/fs"
	"github.com/jpillora/cloud-torrent/cloudtorrent/module"
)

func New() module.Module {
	m := &torrentModule{}
	m.setupRoutes()
	return m
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
	client *torrent.Client
	router http.Handler
}

func (t *torrentModule) TypeID() string {
	return "torrent" //singleton
}

func (t *torrentModule) Configure(raw json.RawMessage) (interface{}, error) {
	err := json.Unmarshal(raw, &t.config)
	if err != nil {
		return nil, err
	}
	unset := t.config.PeerID == "" && t.config.IncomingPort == 0
	if t.config.PeerID == "" {
		t.config.PeerID = "Cloud-Torrent"
	}
	//must be 20 chars
	for len(t.config.PeerID) < 20 {
		t.config.PeerID += " "
	}
	if len(t.config.PeerID) > 20 {
		t.config.PeerID = t.config.PeerID[:20]
	}
	if t.config.IncomingPort == 0 {
		t.config.IncomingPort = 50007
	}
	if unset {
		t.config.EnableEncryption = true
		t.config.EnableSeeding = true
		t.config.EnableUpload = true
	}
	//close previous client
	if t.client != nil {
		t.client.Close()
	}
	//convert cloud-torrent config into anacrolix/torrent config
	t.client, err = torrent.NewClient(&torrent.Config{
		PeerID:            t.config.PeerID,
		NoUpload:          !t.config.EnableUpload,
		Seed:              t.config.EnableSeeding,
		DisableEncryption: !t.config.EnableEncryption,
		ListenAddr:        fmt.Sprintf(":%d", t.config.IncomingPort),
	})
	if err != nil {
		return nil, err
	}
	return &t.config, nil
}

func (t *torrentModule) Sync(updates chan fs.Node) error {
	for {
		updates <- t.sync()
		time.Sleep(1 * time.Second)
	}
}

func (t *torrentModule) sync() fs.Node {
	torrents := &fs.BaseNode{BaseInfo: fs.BaseInfo{IsDir: true}}
	for _, t := range t.client.Torrents() {
		tnode := torrentToNode(t)
		torrents.Upsert(tnode.Name(), tnode)
	}
	return torrents
}

func logf(format string, args ...interface{}) {
	log.Printf("[Torrents] "+format, args...)
}
