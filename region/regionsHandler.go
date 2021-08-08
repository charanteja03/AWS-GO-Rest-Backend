package region

import (
	"net/http"
	"sfr-backend/response"

	"github.com/aws/aws-sdk-go/aws/endpoints"
)

// GetDefaultRegion - gets region from request or defaults to us-east-1
func GetDefaultRegion(r *http.Request) string {
	urlParams := r.URL.Query()

	if len(urlParams.Get("region")) > 0 {
		return urlParams.Get("region")
	}

	return "us-east-1"
}

// GetRegionsHandler - returns list of all aws regions
func GetRegionsHandler(w http.ResponseWriter, r *http.Request) {
	regionList := []string{}

	resolver := endpoints.DefaultResolver()
	partitions := resolver.(endpoints.EnumPartitions).Partitions()

	for _, p := range partitions {
		for id := range p.Regions() {
			regionList = append(regionList, id)
		}
	}
	// fmt.Println("regions", regionList)
	response.WriteResponse(w, regionList)
}
