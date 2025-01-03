package response

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Status string `json:"status"`
	Error string `json:"error"`
}

const (
	StatusOk = "OK"
	StatusError = "ERROR"
)

func WriteJSON(w http.ResponseWriter, data interface{}, status int) error{

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	
	return json.NewEncoder(w).Encode(data)
}

func GeneralError(err error) Response {
	return Response{
		Status: StatusError,
		Error: err.Error(),
	}
}