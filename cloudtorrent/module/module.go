package module

import "encoding/json"

type ModuleFactory interface {
	Type() string
	New(id string) Module
}

type Module interface {
	//Module type and ID joined with a colon
	TypeID() string
	Configure(json.RawMessage) (interface{}, error)
}
