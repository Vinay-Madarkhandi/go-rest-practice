package student

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Vinay-Madarkhandi/go-rest-practice/internal/storage"
	"github.com/Vinay-Madarkhandi/go-rest-practice/internal/types"
	"github.com/Vinay-Madarkhandi/go-rest-practice/internal/utils/response"
	"github.com/go-playground/validator/v10"
)

func NewStudent(storage storage.Storage) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		var student types.Student

		// Error Handelling while incomming request
		err := json.NewDecoder(request.Body).Decode(&student)
		if errors.Is(err, io.EOF) {
			slog.Error("")
			err := response.WriteJSON(writer, http.StatusBadRequest, response.GeneralError(err))
			if err != nil {
				slog.Error("Something went wrong while decoding", slog.String("Error", err.Error()))
				return
			}
			return
		}
		if err != nil {
			slog.Error("Something went wrong", slog.String("Error", err.Error()))
			err := response.WriteJSON(writer, http.StatusBadRequest, response.GeneralError(err))
			if err != nil {
				slog.Error("Internal Server error", slog.String("Error", err.Error()))
				return
			}
			return
		}

		// Input Request Validation
		err = validator.New().Struct(&student)
		if err != nil {
			var validateErrs validator.ValidationErrors
			errors.As(err, &validateErrs)
			slog.Error("Error while validationg", slog.String("Error", err.Error()))
			response.WriteJSON(writer, http.StatusBadRequest, response.ValidationError(validateErrs))
			return
		}

		// Creating Student
		slog.Info("Creating a Student")
		id, err := storage.CreateStudent(student.Name, student.Email, student.Age)
		if err != nil {
			slog.Error("Error while creating the student", slog.String("Error", err.Error()))
			response.WriteJSON(writer, http.StatusBadRequest, response.GeneralError(err))
		}
		student.Id = id
		err = response.WriteJSON(
			writer,
			http.StatusCreated,
			types.StudentCreatedDTO{
				Success: true,
				Data:    student,
			})
		if err != nil {
			slog.Error("Error while creating the student", slog.String("Error", err.Error()))
			err := response.WriteJSON(writer, http.StatusInternalServerError, response.GeneralError(err))
			if err != nil {
				slog.Error("error while writing response", slog.String("Error", err.Error()))
			}
			return
		}

	}
}

func GetStudentByID(storage storage.Storage) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		id := request.PathValue("id")

		slog.Info("Getting student by ID", slog.String("id", id))

		intId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			slog.Error("Error while parsing", slog.String("Error", err.Error()))
			err = response.WriteJSON(writer, http.StatusBadRequest, response.GeneralError(err))
			if err != nil {
				slog.Error("error while writing response", slog.String("Error", err.Error()))
				return
			}
			return
		}

		student, err := storage.GetById(intId)
		if err != nil {
			slog.Error("Error while getting student", slog.String("Error", err.Error()))
			err := response.WriteJSON(writer, http.StatusInternalServerError, response.GeneralError(err))
			if err != nil {
				slog.Error("error while writing response", slog.String("Error", err.Error()))
				return
			}
			return
		}

		err = response.WriteJSON(writer, http.StatusOK, student)
		if err != nil {
			slog.Error("error while writing response", slog.String("Error", err.Error()))
			return
		}

	}
}
