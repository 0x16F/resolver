package errors

import (
	"net/http"
	"os"

	"github.com/goccy/go-json"
)

type Service struct {
	errors map[int]Error
}

func New(errorsPath string) Service {
	service := Service{
		errors: make(map[int]Error),
	}

	data, err := os.ReadFile(errorsPath)
	if err != nil {
		return service
	}

	var errors []Error
	if err := json.Unmarshal(data, &errors); err != nil {
		return service
	}

	for _, err := range errors {
		service.errors[err.Code] = err
	}

	return service
}

type Error struct {
	Code        int    `json:"code"`
	HttpCode    int    `json:"http_code"`
	Message     string `json:"message"`
	Description string `json:"description"`
}

func (ce Error) Error() string {
	encoded, _ := json.Marshal(ce)
	return string(encoded)
}

func (e Service) GetError(code int) error {
	if err, ok := e.errors[code]; ok {
		return &err
	}

	return &Error{
		Code:     0,
		HttpCode: http.StatusInternalServerError,
		Message:  "Unknown error",
	}
}
