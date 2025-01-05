package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/tushar0305/students-api/internal/storage"
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

// New creates a new student
func New(storage storage.Storage) http.HandlerFunc {
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

		//request validation

		if err := validator.New().Struct(student); err != nil {

			validateErrs := err.(validator.ValidationErrors)
			response.WriteJSON(w, response.ValidationError(validateErrs), http.StatusBadRequest)
			return 
		}

		lastId, err := storage.CreateStudent(
			student.Name, 
			student.Email, 
			student.Age,
		)

		slog.Info("Student Created", slog.String("userId", fmt.Sprint(lastId)))

		if err != nil {
			response.WriteJSON(w, response.GeneralError(err), http.StatusInternalServerError)
			return 
		}

		response.WriteJSON(w, map[string] int64{"id": lastId}, http.StatusCreated)
	}
}


// GetById gets a student by id
func GetById(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		slog.Info("getting a student", slog.String("id", id))

		IntId, err := strconv.ParseInt(id, 10, 64)

		if err != nil {
			response.WriteJSON(w, response.GeneralError(fmt.Errorf("invalid id")), http.StatusBadRequest)
			return
		}

		student, err := storage.GetStudentById(IntId)
		if err != nil {
			slog.Error("Student not found", slog.String("id", id))
			response.WriteJSON(w, response.GeneralError(err), http.StatusInternalServerError)
			return
		}

		response.WriteJSON(w, student, http.StatusOK)
	}
}

// GetList gets a list of students
func GetList(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("getting all students")

		students, err := storage.GetStudents()
		if err != nil {
			slog.Error("Failed to get students", slog.String("error", err.Error()))
			response.WriteJSON(w, response.GeneralError(err), http.StatusInternalServerError)
			return
		}

		response.WriteJSON(w, students, http.StatusOK)
	}
}