package error

import (
	"net/http"
)

// HandleError - returns bad request response to user
func HandleError(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusBadRequest)
}
