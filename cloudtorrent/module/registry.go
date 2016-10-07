package module

import (
	"errors"
	"log"
	"regexp"
)

type ModuleFactory interface {
	Type() string
	New(id string) Module
}

//global module registry
var moduleFactories = map[string]ModuleFactory{}

func Register(factory ModuleFactory) {
	moduleType := factory.Type()
	_, ok := moduleFactories[moduleType]
	if ok {
		log.Panicf("module already registered: %s", moduleType)
	}
	moduleFactories[moduleType] = factory
	log.Printf("registered module: %s", moduleType)
}

var isHex = regexp.MustCompile(`[a-f0-9]`)

func New(moduleType, moduleID string) (Module, error) {
	f, ok := moduleFactories[moduleType]
	if !ok {
		return nil, errors.New("missing module type: " + moduleType)
	}
	if moduleID == "" {
		return nil, errors.New("empty module id")
	}
	if !isHex.MatchString(moduleID) {
		return nil, errors.New("non-hex module id: " + moduleID)
	}
	return f.New(moduleID), nil
}
