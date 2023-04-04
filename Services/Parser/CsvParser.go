package Parser

import (
	"context"
	"cpool/Configurator"
	"cpool/Helpers"
	"cpool/Models"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"
)

type CsvParser struct {
	cfg  *Configurator.Cfg
	file *os.File
}

func (parser *CsvParser) Initial(cfg *Configurator.Cfg) error {
	parser.cfg = cfg
	return nil
}

func (parser *CsvParser) Run() {
	var err error
	go func() {
		t1 := time.Now()
		for {
			parser.file, err = os.Open(parser.cfg.CsvParser.FilePath)
			if err != nil {
				fmt.Println("Error: " + err.Error())
			} else {
				fmt.Println("Start parsing")
				csvReader := csv.NewReader(parser.file)
				// TODO delete this after completing
				blockOfInserts := make([]Models.Promotion, 0, parser.cfg.CsvParser.InsertsCount)
				for {
					rec, err := csvReader.Read()
					if err == io.EOF {
						break
					}
					if err != nil {
						fmt.Println("Error: " + err.Error())
					} else {
						var promotion Models.Promotion
						promotion.IdHash = rec[0]
						promotion.Price, _ = strconv.ParseFloat(rec[1], 8)
						promotion.ExpirationDate, err = time.Parse("2006-01-02 15:04:05 +0200 CEST", rec[2])

						blockOfInserts = append(blockOfInserts, promotion)
						if len(blockOfInserts) >= parser.cfg.CsvParser.InsertsCount {
							err = parser.InsertPortionOfPromotions(blockOfInserts)
							if err != nil {
								fmt.Printf("Err[!]: During insert: %s", err.Error())
							}
							blockOfInserts = nil
							blockOfInserts = make([]Models.Promotion, 0, parser.cfg.CsvParser.InsertsCount)
						}
					}
				}

				if len(blockOfInserts) > 0 {
					err = parser.InsertPortionOfPromotions(blockOfInserts)
					if err != nil {
						fmt.Printf("Err[!]: During insert: %s", err.Error())
					}
				}

				// remember to close the file at the end of the program
				err := parser.file.Close()
				if err != nil {
					fmt.Println(err.Error())
				}
				parser.file = nil
				os.Remove(parser.cfg.CsvParser.FilePath)
			}
			t2 := time.Now()
			fmt.Println(fmt.Sprintf("Parsing of file taked %f seconds", t2.Sub(t1).Seconds()))
			time.Sleep(time.Second * time.Duration(parser.cfg.CsvParser.ParsingDuration))

		}
	}()
}

func (parser *CsvParser) InsertPortionOfPromotions(blockOfInserts []Models.Promotion) error {
	return Helpers.InsertPromotions(blockOfInserts, parser.cfg.Pgsql.Connection)
}

func (parser *CsvParser) Close(ctx context.Context) error {
	return parser.handleCancel(ctx)
}

func (parser *CsvParser) handleCancel(ctx context.Context) error {
	if parser.file != nil {
		err := parser.file.Close()
		if err != nil {
			return err
		}
	}
	return nil
}
