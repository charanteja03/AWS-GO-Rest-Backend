package response

import (
	"encoding/json"
	"net/http"

	"sfr-backend/error"
)

// WriteResponse - Writes response encoded as json
func WriteResponse(w http.ResponseWriter, result interface{}) {

	w.Header().Set("Content-Type", "application/json")

	js, err := json.Marshal(result)
	if err != nil {
		error.HandleError(w, err)
		return
	}
	w.Write(js)
}
