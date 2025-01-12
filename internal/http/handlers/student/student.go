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
	"github.com/nikhilsharma270027/GOLang-student-api/internal/storage"
	"github.com/nikhilsharma270027/GOLang-student-api/internal/types"
	"github.com/nikhilsharma270027/GOLang-student-api/internal/utils/response"
)

// dependence injection : student like plug and play
func New(storage storage.Storage) http.HandlerFunc {
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

		// as a dependenice we receiving / props
		lastId, err := storage.CreateStudent(
			student.Name,
			student.Email,
			student.Age,
		)

		slog.Info("user created successfully!!", slog.String("userId", fmt.Sprint(lastId)))

		if err != nil { // it is a database err
			response.WriteJson(w, http.StatusInternalServerError, err)
			return
		}

		// w.Write([]byte("We- welcome to student api"))
		// response.WriteJson(w, http.StatusCreated, map[string]string{"success": "ok"}) // changed
		response.WriteJson(w, http.StatusCreated, map[string]int64{"id:": lastId})
		// if rest api is created we use 200 , like its working good

	}
}

// we dont get data directly we need to decode and use struct to serilize

// we can use custom error fmt.Errorf("error body")

func GetById(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		id := r.PathValue("id") // to find by url /{id}

		slog.Info("Getting a student", slog.String("id", id))

		intId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		student, err := storage.GetStudentById(intId)
		if err != nil {
			slog.Error("error getting user", slog.String("id", id))
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		response.WriteJson(w, http.StatusOK, student)
	}
}

func GetList(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		slog.Info("Getting all student")

		students, err := storage.GetStudents()
		if err != nil {
			// internal server error/ database error
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		response.WriteJson(w, http.StatusOK, students)
	}
}
