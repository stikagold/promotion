package Api

import (
	"context"
	"cpool/Configurator"
	"cpool/Controllers"
	"cpool/Helpers"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/gorilla/mux"
	"log"
	"net/http"
)

type Provider struct {
	cfg    *Configurator.Cfg
	server *http.Server
	router *mux.Router
}

func (prs *Provider) Initial(cfg *Configurator.Cfg) error {
	var err error
	prs.cfg = cfg
	prs.router = mux.NewRouter()
	prs.server = &http.Server{
		Addr:    prs.cfg.Api.Host + ":" + prs.cfg.Api.Port,
		Handler: prs.router,
	}
	return err
}

func (prs *Provider) Run() error {
	var prv Controllers.Provider
	exampleC, err := prv.GetController(Helpers.EXAMPLE_ENTITY)
	if err != nil {
		fmt.Printf("Err[!]: Can not find controller %s", Helpers.EXAMPLE_ENTITY)
	} else {
		_ = exampleC.Initial(prs.cfg)
		if err = exampleC.Run(prs.router); err != nil {
			fmt.Printf("Err[!]: Can not run controller %s", Helpers.EXAMPLE_ENTITY)
		}
	}

	promotionsC, err := prv.GetController(Helpers.PROMOTIONS_ENTITY)
	if err != nil {
		fmt.Printf("Err[!]: Can not find controller %s", Helpers.PROMOTIONS_ENTITY)
	} else {
		_ = promotionsC.Initial(prs.cfg)
		if err = promotionsC.Run(prs.router); err != nil {
			fmt.Printf("Err[!]: Can not run controller %s", Helpers.PROMOTIONS_ENTITY)
		}
	}

	// TODO - declare and serve started from here
	log.Print(prs.cfg.Api.Host + ":" + prs.cfg.Api.Port + " - Listening...")
	go func() {
		if err = prs.server.ListenAndServe(); err != nil {
			fmt.Printf("Err[!]: Can not http server %s", err.Error())
			return
		}

	}()

	return nil
}

func (prs *Provider) IsAutoConnect() bool {
	return prs.cfg.Api.AutoConnect
}

func (prs *Provider) Close(ctx context.Context) error {
	return prs.server.Shutdown(ctx)
}
