package server_test

import (
	"net/http"
	"net/http/httptest"
	"sfr-backend/server"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func dummyHandler(w http.ResponseWriter, r *http.Request) {}

func TestCorsMiddlewareTest(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/", dummyHandler).Methods("OPTIONS", "GET")
	router.Use(server.CORS)

	testTable := []string{
		"OPTIONS", "GET",
	}

	for _, testCase := range testTable {
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest(testCase, "/", nil)

		router.ServeHTTP(rr, req)

		assert.NotEqual(t, "", req.Header.Clone().Get("X-Request-ID"))
		for _, responseHeader := range rr.HeaderMap {
			assert.NotEqual(t, "", responseHeader)
		}
	}
}
