package server

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"sfr-backend/authentication"
	awsprovider "sfr-backend/awsProvider"
	"sfr-backend/execution"
	"sfr-backend/healthcheck"
	"sfr-backend/machine"
	"sfr-backend/region"
	"sfr-backend/tid"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

// CORS Middleware
func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tid := tid.GetTid(r)
		requestLogger := log.WithFields(log.Fields{"transaction_id": tid})
		// Set headers
		w.Header().Set("Access-Control-Allow-Origin", os.Getenv("BASE_URL"))
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		r.Header.Add("X-Request-ID", tid)

		if r.Method == "OPTIONS" {
			requestLogger.Info("OPTIONS operation received, will send response and end...")
			w.WriteHeader(http.StatusOK)
			return
		}
		requestLogger.Info("Operation received, will process operation")
		requestLogger.Info(formatRequest(r))
		next.ServeHTTP(w, r)
		return
	})
}

func formatRequest(r *http.Request) string {
	var request []string
	url := fmt.Sprintf("%v %v %v", r.Method, r.URL, r.Proto)
	request = append(request, url)
	request = append(request, fmt.Sprintf("Host: %v", r.Host))
	for name, headers := range r.Header {
		name = strings.ToLower(name)
		for _, h := range headers {
			request = append(request, fmt.Sprintf("%v: %v", name, h))
		}
	}
	if r.Method == "POST" {
		r.ParseForm()
		request = append(request, "\n")
		request = append(request, r.Form.Encode())
	}
	// Return the request as a string
	return strings.Join(request, "\n")
}

// StartServer - starts server and setups possible routes for server
func StartServer() http.Handler {
	fmt.Println("Starting server.")
	router := mux.NewRouter()

	// We use our custom CORS Middleware
	router.Use(CORS)

	router.Handle("/aws/machines", authentication.CheckAuthentication(
		func(w http.ResponseWriter, r *http.Request) {
			machine.GetMachinesHandler(w, r, &awsprovider.AwsStepFunctionsRealProvider{})
		})).Methods("GET")

	router.Handle("/aws/regions", authentication.CheckAuthentication(region.GetRegionsHandler)).Methods("GET")

	router.Handle("/aws/executions", authentication.CheckAuthentication(
		func(w http.ResponseWriter, r *http.Request) {
			execution.GetExecutionsHandler(w, r, &awsprovider.AwsStepFunctionsRealProvider{})
		})).Methods("GET").Queries("machine", "{machine}")

	router.Handle("/aws/execution/{execution}", authentication.CheckAuthentication(
		func(w http.ResponseWriter, r *http.Request) {
			execution.GetExecutionHandler(w, r, &awsprovider.AwsStepFunctionsRealProvider{})
		})).Methods("GET")

	router.Handle("/aws/execution", authentication.CheckAuthentication(
		func(w http.ResponseWriter, r *http.Request) {
			execution.PostStartExecution(w, r, &awsprovider.AwsStepFunctionsRealProvider{})
		})).Methods("POST")

	router.Handle("/aws/execution/restart", authentication.CheckAuthentication(
		func(w http.ResponseWriter, r *http.Request) {
			execution.PostRestartExecution(w, r, &awsprovider.AwsStepFunctionsRealProvider{})
		})).Methods("POST")

	router.Handle("/aws/execution/batch", authentication.CheckAuthentication(
		func(w http.ResponseWriter, r *http.Request) {
			execution.PostRestartBatch(w, r, &awsprovider.AwsStepFunctionsRealProvider{})
		})).Methods("POST")

	router.HandleFunc("/logout", authentication.Logout).Methods("GET", "OPTIONS")

	router.HandleFunc("/login", authentication.LoginHandler).Methods("POST", "OPTIONS")

	router.HandleFunc("/createuser", authentication.CreateUser).Methods("POST", "OPTIONS")

	router.HandleFunc("/refreshtoken", authentication.RefreshTokenCheck).Methods("POST", "OPTIONS")

	router.HandleFunc("/healthcheck", healthcheck.Healthcheck).Methods("GET")
	http.Handle("/", router)
	return router

}
