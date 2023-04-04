package Response

import "encoding/json"

type ApiResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func (ar *ApiResponse) ToByte() ([]byte, error) {
	return json.Marshal(ar)
}
