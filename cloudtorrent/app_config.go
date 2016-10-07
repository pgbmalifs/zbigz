package cloudtorrent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sort"
	"strings"

	"github.com/jpillora/cloud-torrent/cloudtorrent/module"
)

var EmptyConfig = json.RawMessage("{}")

type AppConfig struct {
	User, Pass string
	Title      string
}

//App itself is also Configurable
func (a *App) Configure(raw json.RawMessage) (interface{}, error) {
	if err := json.Unmarshal(raw, &a.appConfig); err != nil {
		return nil, err
	}
	if a.appConfig.Title == "" {
		a.appConfig.Title = "Cloud Torrent"
	}
	a.auth.SetUserPass(a.appConfig.User, a.appConfig.Pass)
	return &a.appConfig, nil
}

func (a *App) configureAllRaw(b []byte) error {
	cfgs := rawMessages{}
	if err := json.Unmarshal(b, &cfgs); err != nil {
		return fmt.Errorf("initial configure failed: %s", err)
	}
	return a.configureAll(cfgs)
}

func (a *App) configureAll(cfgs rawMessages) error {
	for id, raw := range cfgs {
		if err := a.configureModule(id, raw); err != nil {
			return err
		}
	}
	return nil
}

func (a *App) configureModule(typeid string, rawConfig json.RawMessage) error {
	//normalize raw json
	indented := bytes.Buffer{}
	err := json.Indent(&indented, rawConfig, "", "  ")
	if err != nil {
		return err
	}
	config := indented.Bytes()
	//check for existing module
	m, ok := a.modules[typeid]
	if !ok {
		//doesn't exist, init using type:id
		pair := strings.SplitN(typeid, ":", 2)
		if len(pair) != 2 {
			return fmt.Errorf("Invalid typeid  ('%s')", typeid)
		}
		typ := pair[0]
		id := pair[1]
		//lookup module registry
		m, err = module.New(typ, id)
		if err != nil {
			return err
		}
		if t := m.TypeID(); typeid != t {
			return fmt.Errorf("New module typeid mismatch ('%s' vs '%s')", typeid, t)
		}
		//initialise module
		a.initModule(m)
	}
	//compare to last update
	prev := a.configs[typeid]
	if bytes.Equal(prev, config) {
		return nil //skip, already have this config
	}
	//apply!
	newConfig, err := m.Configure(config)
	if err != nil {
		if bytes.Equal(config, EmptyConfig) {
			return nil //skip empty config errors
		}
		logf("[%s] configuration error: %s", typeid, err)
		return err
	}
	//convert result
	config, err = json.MarshalIndent(newConfig, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal: %s (%s)", typeid, err)
	}
	//successful configure!
	a.configs[typeid] = config
	mstate := a.state.Modules[typeid]
	if !mstate.Enabled {
		mstate.Sync(m)
		mstate.Lock()
		mstate.Enabled = true
		mstate.Unlock()
	}
	mstate.Lock()
	mstate.Config = newConfig
	mstate.Unlock()
	mstate.Push()
	//write back to disk if changed
	b, _ := json.MarshalIndent(&a.configs, "", "  ")
	ioutil.WriteFile(a.ConfigPath, b, 0600)
	logf("reconfigured %s", typeid)
	return nil
}

//rawMessages allows json marshalling of string->raw
type rawMessages map[string]json.RawMessage

func (m rawMessages) MarshalJSON() ([]byte, error) {
	buf := bytes.Buffer{}
	keys := make([]string, len(m))
	i := 0
	for k, _ := range m {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	//manually write object
	buf.WriteString("{")
	for i, k := range keys {
		buf.WriteString(`"`)
		buf.WriteString(k)
		buf.WriteString(`":`)
		buf.Write(m[k])
		if i < len(keys)-1 {
			buf.WriteRune(',')
		}
	}
	buf.WriteString("}")
	return buf.Bytes(), nil
}
