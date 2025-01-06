package response

import "net/http"

// response object w. http.ResponseWriter, status int, data
// dont know wht type of data would reveive
func WriteJson(w http.ResponseWriter, status int, data int)
