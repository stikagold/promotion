package main

import (
	"context"
	"cpool/Configurator"
	"cpool/Helpers"
	"cpool/Services/Api"
	"cpool/Services/Parser"
	"fmt"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	var (
		wGroup    sync.WaitGroup
		wservice  Api.Provider
		csvParser Parser.CsvParser
	)

	ctxCancel, fnCancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	defer fnCancel()

	var cfg, err = Configurator.GetConfigurator()
	if err != nil {
		panic(err)
	}

	// This Part will run only in mode of parser to seed database from file
	if cfg.Mode == Helpers.PARSING_MODE {
		err = csvParser.Initial(cfg)
		if err != nil {
			fmt.Printf("Err[!]: Unable to initial CsvParser: %s", err.Error())
		} else {
			wGroup.Add(1)
			csvParser.Run()
			wGroup.Done()
		}
	}

	// This part only for api gateway
	if cfg.Mode == Helpers.API_MODE {
		err = wservice.Initial(cfg)
		if err != nil {
			fmt.Printf("Err[!]: Unable to initial ApiProvider: %s", err.Error())
		} else {
			wGroup.Add(1)
			go func() {
				err := wservice.Run()
				if err != nil {
					fmt.Printf("Err[!]: Unable to run ApiProvider: %s", err.Error())
				}
				wGroup.Done()
			}()
		}
	}

	wGroup.Wait()

	// Start resources unloading process
	Helpers.ShowMemoryUsage()
	<-ctxCancel.Done()
	err = cfg.Close(ctxCancel)
	if err != nil {
		fmt.Println("Err[!]", err.Error())
	}
	if cfg.Mode == Helpers.API_MODE {
		err = wservice.Close(ctxCancel)
		if err != nil {
			fmt.Println("Err[!]", err.Error())
		}
	}
	if cfg.Mode == Helpers.PARSING_MODE {
		err = csvParser.Close(ctxCancel)
		if err != nil {
			fmt.Println("Err[!]", err.Error())
		}
	}
}
