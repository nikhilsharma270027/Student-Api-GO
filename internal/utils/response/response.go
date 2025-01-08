package response

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}

const (
	StatusOk    = "OK"
	StatusError = "Error"
)

// response object w. http.ResponseWriter, status int, data
// dont know wht type of data would reveive
// will return error in response
func WriteJson(w http.ResponseWriter, status int, data interface{}) error {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	// the incoming data we had decode the data to put in struct
	// we eill use encode here
	// to convert struct code to json we use encode
	return json.NewEncoder(w).Encode(data)
}

// to get error as json not "EOF"
func GeneralError(err error) Response {
	return Response{
		Status: StatusError,
		Error:  err.Error(),
	}
}

// response for validation errors
func ValidateError(errs validator.ValidationErrors) Response {
	var errMsgs []string

	for _, err := range errs {
		// from package we get err accoding to the tag we set on student struct
		switch err.ActualTag() {
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is required field", err.Field()))
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is invalid", err.Field()))
		}
	}

	return Response{
		Status: StatusError,
		Error:  strings.Join(errMsgs, ", "),
		// join will add slice to string means convert to string
	}
}
