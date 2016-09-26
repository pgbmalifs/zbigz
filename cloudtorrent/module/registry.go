package module

import "log"

//global module registry
var moduleFactories = map[string]ModuleFactory{}

func Register(factory ModuleFactory) {
	moduleType := factory.Type()
	_, ok := moduleFactories[moduleType]
	if ok {
		log.Panicf("module already registered: %s", moduleType)
	}
	moduleFactories[moduleType] = factory
}

func New(moduleType, moduleID string) Module {
	f, ok := moduleFactories[moduleType]
	if !ok {
		return nil
	}
	return f.New(moduleID)
}
