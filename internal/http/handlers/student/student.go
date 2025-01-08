package student

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/nikhilsharma270027/GOLang-student-api/internal/types"
	"github.com/nikhilsharma270027/GOLang-student-api/internal/utils/response"
)

func New() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Creating a student")

		// we use struct type
		var student types.Student
		// we encode json we use json package
		// In go error is return
		err := json.NewDecoder(r.Body).Decode(&student)
		// here we check wheather the body is empty or not
		if errors.Is(err, io.EOF) {
			//we can send w like response but we want json response
			// response is from response package
			// response.WriteJson(w, http.StatusBadRequest, err.Error()) // passed error
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err)) // passed error
			// badResquest means here we receive empty body/ error 400 code
			return
		}

		// if its not empty body err
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		// request validation - we need to import package
		if err := validator.New().Struct(student); err != nil {
			validateErrs := err.(validator.ValidationErrors)
			// here normal err can be passed as validator take its own error
			response.WriteJson(w, http.StatusBadRequest, response.ValidateError(validateErrs))
			return
		}

		// w.Write([]byte("We- welcome to student api"))
		response.WriteJson(w, http.StatusCreated, map[string]string{"success": "ok"})
		// if rest api is created we use 200 , like its working good

	}
}

// we dont get data directly we need to decode and use struct to serilize

// we can use custom error fmt.Errorf("error body")
