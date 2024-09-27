package entity

type OutlineCreateUserReq struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	OutlineInfo
}

type OutlineCreateUserResp struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Password  string `json:"password"`
	Port      uint16 `json:"port"`
	Method    string `json:"method"`
	AccessUrl string `json:"accessUrl"`
}

type OutlineDeleteUserReq struct {
	UserID string
	OutlineInfo
}

type OutlineInfo struct {
	OutlineURL    string `json:"-"`
	OutlinePort   uint16 `json:"-"`
	OutlineSecret string `json:"-"`
}
