package Databases

import (
	"context"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

type Redis struct {
	Driver       string `json:"driver"`
	Host         string `json:"host"`
	Port         string `json:"port"`
	DatabaseName int    `json:"name"`
	User         string `json:"user"`
	Password     string `json:"password"`
	AutoConnect  bool   `json:"auto_connect"`
	Expiration   int    `json:"expiration"`
	Connection   *redis.Client
}

func (rds *Redis) IsEmpty() bool {
	return !rds.AutoConnect
}

func (rds *Redis) Initial() error {
	if !rds.IsEmpty() {
		var ctx = context.TODO()
		var err error

		rds.Connection = redis.NewClient(&redis.Options{
			Addr:     rds.Host + ":" + rds.Port,
			Password: rds.Password, // no password set
			DB:       0,            // use default DB
		})

		err = rds.Connection.Set(ctx, "check", true, 0).Err()
		if err != nil {
			return err
		}

	}
	return nil
}

func (rds *Redis) GetConnection() (*redis.Client, error) {
	return rds.Connection, nil
}

func (rds *Redis) GetDatabase() (*redis.Client, error) {
	return rds.Connection, nil
}

func (rds *Redis) Close() error {
	return rds.Connection.Close()
}
