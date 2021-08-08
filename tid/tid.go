package tid

import (
	"net/http"
	"strconv"
	"time"
)

//generateTid - returns Transaction ID from headers or generates a new one
func generateTid() string {
	return strconv.FormatInt(time.Now().UnixNano(), 10)
}

//GetTid - returns Transaction ID from headers or generates a new one
func GetTid(r *http.Request) string {
	tid := r.Header.Get("X-Request-ID")
	if tid == "" {
		tid = generateTid()
	}
	return tid
}
