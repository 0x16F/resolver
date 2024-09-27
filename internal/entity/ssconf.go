package entity

type SSConf struct {
	Server   string `json:"server"`
	Port     uint16 `json:"server_port"`
	Password string `json:"password"`
	Method   string `json:"method"`
}
