package module

import "encoding/json"

type Module interface {
	//Module type and ID joined with a colon
	TypeID() string
	Configure(json.RawMessage) (interface{}, error)
}
