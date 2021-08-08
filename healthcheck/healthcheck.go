package healthcheck

import (
	"net/http"
	"sfr-backend/response"
	"sfr-backend/tid"

	//envs
	log "github.com/sirupsen/logrus"
)

//ResponseHealthcheck - Response struct for healthcheck data
type ResponseHealthcheck struct {
	Status string `json:"status"`
}

// Healthcheck - returns an empty OK response for load balancers
func Healthcheck(w http.ResponseWriter, r *http.Request) {
	tid := tid.GetTid(r)
	requestLogger := log.WithFields(log.Fields{"transaction_id": tid})
	var resp ResponseHealthcheck
	resp.Status = "OK"
	w.Header().Set("Content-Type", "application/json")
	requestLogger.Info("Healthcheck received OK, will response with a 200 OK")
	response.WriteResponse(w, resp)
}
