package nwlib

import (
	_ "github.com/networth-app/networth/dotenv"

	"encoding/json"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/endpoints"
	"github.com/aws/aws-sdk-go-v2/aws/external"
)

// APIResponse api reponse
type APIResponse struct {
	Data interface{} `json:"data"`
}

// GetEnv get environment variable with default fall back
func GetEnv(params ...string) string {
	if value, ok := os.LookupEnv(params[0]); ok {
		return value
	} else if len(params) >= 2 {
		return params[1]
	}

	return ""
}

// ErrorResp error response helper for api
func ErrorResp(w http.ResponseWriter, message interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)

	msg := message.(string)
	json.NewEncoder(w).Encode(APIResponse{msg})
}

// SuccessResp success response helper for api
func SuccessResp(w http.ResponseWriter, message interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	switch message.(type) {
	case string:
		json.NewEncoder(w).Encode(APIResponse{message.(string)})
		break

	default:
		json.NewEncoder(w).Encode(APIResponse{message})
	}

}

// LoadAWSConfig set default aws config
func LoadAWSConfig() aws.Config {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		panic("Unable to load SDK config: " + err.Error())
	}

	cfg.Region = endpoints.UsEast1RegionID

	return cfg
}
