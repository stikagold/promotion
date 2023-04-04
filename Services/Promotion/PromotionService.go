package Promotion

import (
	"context"
	"cpool/Configurator"
	"cpool/Models"
	"encoding/json"
	"errors"
	"strconv"
	"time"
)

type Promotion struct {
	cfg *Configurator.Cfg
}

func (pr *Promotion) Initial(cfg *Configurator.Cfg) {
	pr.cfg = cfg
}

func (pr *Promotion) GetPromotion(id int) (*Models.Promotion, error) {
	idc := strconv.Itoa(id)
	promotion, err := pr.GetFromCache(idc)
	if err != nil {
		var promotionFromPgsql Models.Promotion
		sqlQuery := "SELECT * FROM promotions where id=" + idc
		row := pr.cfg.Pgsql.Connection.QueryRow(sqlQuery)
		err := row.Scan(&promotionFromPgsql.Id, &promotionFromPgsql.IdHash, &promotionFromPgsql.Price, &promotionFromPgsql.ExpirationDate)
		if err != nil {
			return nil, errors.New("not found")
		}
		promotionMarshaled, _ := json.Marshal(promotionFromPgsql)
		_ = pr.SetToCache(promotionMarshaled, promotionFromPgsql.Id)
		return &promotionFromPgsql, nil
	}
	return promotion, nil
}

func (pr *Promotion) GetFromCache(id string) (*Models.Promotion, error) {
	var promotion Models.Promotion
	response, _ := pr.cfg.Redis.Connection.Get(context.Background(), id).Result()
	if response == "" {
		return nil, errors.New("not found")
	}
	err := json.Unmarshal([]byte(response), &promotion)
	if err != nil {
		return nil, errors.New("not found")
	}
	return &promotion, nil
}

func (pr *Promotion) SetToCache(promotion []byte, id int) error {
	return pr.cfg.Redis.Connection.Set(context.TODO(), strconv.Itoa(id), promotion, time.Second*time.Duration(pr.cfg.Redis.Expiration)).Err()
}
