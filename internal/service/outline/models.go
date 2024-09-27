package outline

type createUser struct {
	Method   string `json:"method"`
	Name     string `json:"name"`
	Password string `json:"password"`
}
