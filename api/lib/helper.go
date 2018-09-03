package nwlib

import (
	"bufio"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

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
	log.Println("Response error: " + msg)
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

// GetAPIVersion api version
func GetAPIVersion() string {
	path := getRootDir() + "/api/Gopkg.toml"
	file, _ := os.Open(path)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	metadataFound := false
	for scanner.Scan() {
		line := scanner.Text()
		if metadataFound {
			lineSplitted := strings.Split(line, "=")
			version := strings.Replace(lineSplitted[1], "\"", "", 2)
			return strings.TrimSpace(version)
		}

		if line == "[metadata]" {
			metadataFound = true
		}
	}

	return ""
}

// GetRootDir root dir
func GetRootDir() string {
	dir, _ := os.Getwd()
	if strings.HasSuffix(dir, "/api") {
		dir = strings.Replace(dir, "/api", "", 1)
	}

	return dir
}

// TODO: remove in favor of Lambda ENV
// func LoadDotEnv() {
// 	envPath := getRootDir() + "/.env"
// 	file, _ := os.Open(envPath)
// 	defer file.Close()

// 	scanner := bufio.NewScanner(file)
// 	for scanner.Scan() {
// 		line := scanner.Text()
// 		lineSplitted := strings.Split(line, "=")

// 		key := strings.TrimSpace(lineSplitted[0])
// 		val := strings.TrimSpace(lineSplitted[1])
// 		os.Setenv(key, val)
// 	}
// }

// LoadAWSConfig set default aws config
func LoadAWSConfig() aws.Config {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		panic("Unable to load SDK config: " + err.Error())
	}

	cfg.Region = endpoints.UsEast1RegionID

	return cfg
}
