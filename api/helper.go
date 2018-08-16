package main

import (
	"bufio"
	"encoding/json"
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

func getEnv(params ...string) string {
	if value, ok := os.LookupEnv(params[0]); ok {
		return value
	} else if len(params) >= 2 {
		return params[1]
	}

	return ""
}

func errorResp(w http.ResponseWriter, message interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)

	json.NewEncoder(w).Encode(APIResponse{message.(string)})
}

func successResp(w http.ResponseWriter, message interface{}) {
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

func getAPIVersion() string {
	path := getRootDir() + "/Gopkg.toml"
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

func getRootDir() string {
	dir, _ := os.Getwd()
	if strings.HasSuffix(dir, "/api") {
		dir = strings.Replace(dir, "/api", "", 1)
	}

	return dir
}

func loadDotEnv() {
	envPath := getRootDir() + "/.env"
	file, _ := os.Open(envPath)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		lineSplitted := strings.Split(line, "=")

		key := strings.TrimSpace(lineSplitted[0])
		val := strings.TrimSpace(lineSplitted[1])
		os.Setenv(key, val)
	}
}

func loadAWSConfig() aws.Config {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		panic("Unable to load SDK config: " + err.Error())
	}

	cfg.Region = endpoints.UsEast1RegionID

	return cfg
}
