package Helpers

import (
	"cpool/Configurator"
	"cpool/Models"
	"cpool/Response"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"runtime"
)

func GetUrl() (string, error) {
	cfg, err := Configurator.GetConfigurator()
	if err != nil || cfg.IsInitialized() != true {
		return "", errors.New("configurator not built")
	}

	if cfg.InternalMode == DEVELOPER_MODE {
		return cfg.Api.Host + ":" + cfg.Api.PortMapping[cfg.Mode], nil
	}

	return cfg.Api.Host + ":" + cfg.Api.Port, nil

}

func InsertPromotions(promotions []Models.Promotion, db *sql.DB) error {
	sqlStr := "INSERT INTO promotions(id_hash, price, expiration_date) VALUES "
	for index, row := range promotions {
		sqlStr += fmt.Sprintf("('%s', %f, '%s')", row.IdHash, row.Price, row.ExpirationDate.UTC().Format("2006-01-02 15:04:05"))
		if index != (len(promotions) - 1) {
			sqlStr += ","
		}
	}
	res, err := db.Query(sqlStr)
	if res != nil {
		_ = res.Close()
	}
	return err
}

func ShowMemoryUsage() (alloc float64, malloc float64, salloc float64, gcount int) {
	var m runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&m)
	alloc = bToMb(m.Alloc)
	malloc = bToMb(m.TotalAlloc)
	salloc = bToMb(m.Sys)
	gcount = runtime.NumGoroutine()

	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Printf("Alloc = %v Mb", bToMb(m.Alloc))
	fmt.Printf("\tHeap size = %v Mb", bToMb(m.HeapInuse))
	fmt.Printf("\tSys = %v Mb", bToMb(m.Sys))
	fmt.Printf("\tCoroutines = %v\n", gcount)
	return alloc, malloc, salloc, gcount
}

func bToMb(b uint64) float64 {
	return float64(b) / 1024 / 1024
}

func WriteAsResponse(resp Response.ApiResponse, code int, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	encoded, err := resp.ToByte()
	if err != nil {
		var tmp Response.ApiResponse
		tmp.Code = http.StatusInternalServerError
		tmp.Message = "Internal unknown error happened"
		encoded, _ = tmp.ToByte()
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(encoded)
		return
	}
	w.WriteHeader(code)
	w.Write(encoded)
}
