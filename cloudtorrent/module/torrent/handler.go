package torrent

import (
	"io"
	"io/ioutil"
	"net/http"

	"golang.org/x/net/context"

	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/metainfo"
	"github.com/jpillora/cloud-torrent/cloudtorrent/util"

	goji "goji.io"
	"goji.io/pat"
)

const (
	MaxTorrentSize = 2e6  //2MB
	MaxURLSize     = 2048 //2KB
)

// Bound to: /modules/torrent/<path>
func (t *torrentModule) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.router.ServeHTTP(w, r)
}

func (t *torrentModule) setupRoutes() {
	mux := goji.NewMux()
	mux.Handle(pat.Post("/file"), util.ErrHandler(t.addTorrentFile))
	mux.Handle(pat.Post("/url"), util.ErrHandler(t.addTorrentURL))
	mux.Handle(pat.Post("/magnet"), util.ErrHandler(t.addTorrentMagnet))
	t.router = mux
}

func (t *torrentModule) addTorrentMagnet(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	uri, err := ioutil.ReadAll(io.LimitReader(r.Body, MaxURLSize))
	if err != nil {
		return err
	}
	mag, err := metainfo.ParseMagnetURI(string(uri))
	if err != nil {
		return err
	}
	torrent, _, err := t.client.AddTorrentSpec(&torrent.TorrentSpec{
		Trackers:    [][]string{mag.Trackers},
		DisplayName: mag.DisplayName,
		InfoHash:    mag.InfoHash,
	})
	if err != nil {
		return err
	}
	w.Write([]byte(torrent.Name()))
	return nil
}

func (t *torrentModule) addTorrentFile(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	torrent, err := t.addTorrentFromReader(r.Body)
	if err != nil {
		return err
	}
	w.Write([]byte(torrent.Name()))
	return nil
}

func (t *torrentModule) addTorrentURL(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	url, err := ioutil.ReadAll(io.LimitReader(r.Body, MaxURLSize))
	if err != nil {
		return err
	}
	resp, err := http.Get(string(url))
	if err != nil {
		return err
	}
	torrent, err := t.addTorrentFromReader(resp.Body)
	if err != nil {
		return err
	}
	w.Write([]byte(torrent.Name()))
	return nil
}

//helper methods

func (t *torrentModule) addTorrentFromReader(r io.Reader) (*torrent.Torrent, error) {
	lr := io.LimitReader(r, MaxTorrentSize)
	mi, err := metainfo.Load(lr)
	if err != nil {
		return nil, err
	}
	return t.client.AddTorrent(mi)
}
