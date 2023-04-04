package Api

type Configuration struct {
	Host        string            `json:"host"`
	Port        string            `json:"port"`
	PortMapping map[string]string `json:"port_mapping"`
	AutoConnect bool              `json:"auto_connect"`
}

func (acfg *Configuration) IsAutoConnect() bool {
	return acfg.AutoConnect
}
