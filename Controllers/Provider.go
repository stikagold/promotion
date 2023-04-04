package Controllers

import (
	"cpool/Configurator"
	"cpool/Helpers"
	"errors"
	"github.com/gorilla/mux"
)

type Runnable interface {
	Run(router *mux.Router) error
	Initial(cfg *Configurator.Cfg) error
}
type Provider struct {
}

func (prv *Provider) GetController(entity string) (Runnable, error) {
	switch entity {
	case Helpers.EXAMPLE_ENTITY:
		var c ExampleController
		return &c, nil
	case Helpers.PROMOTIONS_ENTITY:
		var p PromotionsController
		return &p, nil
	}

	return nil, errors.New("err[!] Controller not found: " + entity)
}
