package Models

import "time"

type Promotion struct {
	Id             int       `db:"id"`
	IdHash         string    `db:"id_hash"`
	Price          float64   `db:"price"`
	ExpirationDate time.Time `db:"expiration_date"`
}

func (prObj *Promotion) GetForResponse() interface{} {
	tmp := make(map[string]interface{})
	tmp["id"] = prObj.IdHash
	tmp["price"] = prObj.Price
	tmp["expiration_date"] = prObj.ExpirationDate.Format("2006-01-02 15:04:05")
	return tmp
}
