package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/tushar0305/students-api/internal/types"
	"github.com/tushar0305/students-api/internal/utils/response"
)

type Response struct {
	Status string `json:"status"`
	Error string `json:"error"`
}

const (
	StatusOk = "OK"
	StatusError = "ERROR"
)

func New() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		slog.Info("New Student Request")

		var student types.Student

		err := json.NewDecoder(r.Body).Decode(&student)

		if errors.Is(err, io.EOF) {
			response.WriteJSON(w, response.GeneralError(fmt.Errorf("empty Body")), http.StatusBadRequest)
			return 
		}

		if err != nil {
			response.WriteJSON(w, response.GeneralError(err), http.StatusBadRequest)
			return
		}

		response.WriteJSON(w, map[string] string{"success": "OK"}, http.StatusCreated)
	}
}
