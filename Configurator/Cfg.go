package Configurator

import (
	"context"
	"cpool/Configurator/Api"
	"cpool/Configurator/Broker"
	"cpool/Configurator/Databases"
	"cpool/Configurator/Parser"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"
)

const path = "env.json"

var lock = &sync.Mutex{}

func newCfg(path string) (*Cfg, error) {
	var cfg Cfg
	err := cfg.Initial(path)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

type Cfg struct {
	isInitialized bool

	Mode         string              `json:"mode"`
	Token        string              `json:"token"`
	Prefix       string              `json:"prefix"`
	InternalMode string              `json:"internal_mode"`
	Mongo        Databases.MongoDB   `json:"mongo"`
	Pgsql        Databases.Postgres  `json:"postgres"`
	Redis        Databases.Redis     `json:"redis"`
	RabbitMQ     Broker.RabbitMQ     `json:"rabbitmq"`
	Api          Api.Configuration   `json:"api"`
	CsvParser    Parser.CsvParserCfg `json:"csv_parser"`
}

func (cfg *Cfg) IsInitialized() bool {
	return cfg.isInitialized
}

func (cfg *Cfg) preInitial() {
}

func (cfg *Cfg) handleCancel(ctx context.Context) error {
	var err error
	fmt.Println("Lets disconnect database connection")
	if !cfg.Mongo.IsEmpty() {
		fmt.Println("Disconnecting MongoDB")
		err := cfg.Mongo.Close(ctx)
		if err != nil {
			log.Fatalf("Error[!]: %s", err.Error())
		}
	}
	if !cfg.Pgsql.IsEmpty() {
		fmt.Println("Disconnecting Postgres")
		err := cfg.Pgsql.Close()
		if err != nil {
			log.Fatalf("Error[!]: %s", err.Error())
		}
	}
	if !cfg.Redis.IsEmpty() {
		fmt.Println("Disconnecting Redis")
		err := cfg.Redis.Close()
		if err != nil {
			log.Fatalf("Error[!]: %s", err.Error())
		}
	}
	if !cfg.RabbitMQ.IsEmpty() {
		fmt.Println("Disconnecting RabbitMQ")
		err := cfg.RabbitMQ.Close()
		if err != nil {
			log.Fatalf("Error[!]: %s", err.Error())
		}
	}
	return err
}

func (cfg *Cfg) Close(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			err := cfg.handleCancel(ctx)
			return err
		default:
			time.Sleep(5 * time.Second)
		}
	}
}

func (cfg *Cfg) Dispatch(b []byte, qName string) error {
	if cfg.RabbitMQ.IsEmpty() != true {
		return cfg.RabbitMQ.Dispatch(b, qName)
	}
	return errors.New("no broker defined to dispatch")
}

func (cfg *Cfg) Initial(path string) error {
	if cfg.isInitialized == true {
		return nil
	}
	cfg.preInitial()

	configFile, err := os.Open(path)
	defer func(configFile *os.File) {
		err := configFile.Close()
		if err != nil {
			fmt.Println(err.Error())
		}
	}(configFile)
	if err != nil {
		return err
	}

	jsonContent, err := io.ReadAll(configFile)
	if err != nil {
		return err
	}

	err = json.Unmarshal(jsonContent, &cfg)
	if err != nil {
		return err
	}

	// Now we should check database objects and if there are no empty - initial them
	if !cfg.Mongo.IsEmpty() {
		fmt.Println("Connecting MongoDB")
		err = cfg.Mongo.Initial()
		if err != nil {
			log.Fatalf("Error[!] %s", err.Error())
		}
	}
	if !cfg.Pgsql.IsEmpty() {
		fmt.Println("Connecting Postgres")
		err = cfg.Pgsql.Initial()
		if err != nil {
			log.Fatalf("Error[!] %s", err.Error())
		}
	}
	if !cfg.Redis.IsEmpty() {
		fmt.Println("Connecting Redis")
		err = cfg.Redis.Initial()
		if err != nil {
			log.Fatalf("Error[!] %s", err.Error())
		}
	}
	if !cfg.RabbitMQ.IsEmpty() {
		fmt.Println("Connecting RabbitMQ")
		err = cfg.RabbitMQ.Initial()
		if err != nil {
			log.Fatalf("Error[!] %s", err.Error())
		}
	}
	cfg.isInitialized = true
	return nil
}

// Configurator processing
var cfgInstance *Cfg

func GetConfigurator() (*Cfg, error) {
	var err error
	if cfgInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if cfgInstance == nil {
			cfgInstance, err = newCfg(path)
		}
	}
	return cfgInstance, err
}
