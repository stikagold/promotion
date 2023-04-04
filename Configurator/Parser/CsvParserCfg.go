package Parser

type CsvParserCfg struct {
	AutoStart       bool   `json:"auto_start"`
	FilePath        string `json:"file_path"`
	InsertsCount    int    `json:"inserts_count"`
	ParsingDuration int    `json:"parsing_duration"`
}

func (cfp *CsvParserCfg) IsAutoRun() bool {
	return cfp.FilePath != "" && cfp.AutoStart
}
