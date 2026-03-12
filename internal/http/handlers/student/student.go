package student

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/Vinay-Madarkhandi/go-rest-practice/internal/types"
	"github.com/Vinay-Madarkhandi/go-rest-practice/internal/utils/response"
	"github.com/go-playground/validator/v10"
)

func CreateStudent() http.HandlerFunc {
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
		e := response.WriteJSON(
			writer,
			http.StatusCreated,
			types.StudentCreatedDTO{
				Success: true,
				Data:    student,
			})
		if e != nil {
			slog.Error("Error while creating the student", slog.String("Error", e.Error()))
			err := response.WriteJSON(writer, http.StatusInternalServerError, response.GeneralError(e))
			if err != nil {
				slog.Error("error while writing response", slog.String("Error", e.Error()))
			}
			return
		}

	}
}
