package student

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/nikhilsharma270027/GOLang-student-api/internal/http/types"
)

func New() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

// we use struct type
		var student types.Student
		// we encode json we use json package
		// In go error is return 
		err := json.NewDecoder(r.Body).Decode(&student)
		// here we check wheather the body is empty or not
		if errors.Is(err, io.EOF) {
			//we can send w like response but we want json response
		}

		slog.Info("Creating a student")

		w.Write([]byte("Welcome to student api"))
	}
}

// we dont get data directly we need to decode and use struct to serilize
