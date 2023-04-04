package Controllers

import (
	"cpool/Configurator"
	"github.com/gorilla/mux"
)

type ExampleController struct {
}

func (c *ExampleController) Initial(cfg *Configurator.Cfg) error {
	return nil
}

func (c *ExampleController) Run(router *mux.Router) error {
	return nil
}
