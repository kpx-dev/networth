package main

import (
	"bufio"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/endpoints"
	"github.com/aws/aws-sdk-go-v2/aws/external"
)

func getEnv(params ...string) string {
	if value, ok := os.LookupEnv(params[0]); ok {
		return value
	} else if len(params) >= 2 {
		return params[1]
	}

	return ""
}

func getAPIVersion() string {
	file, _ := os.Open("Gopkg.toml")
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

func loadDotEnv() {
	file, _ := os.Open(".env")
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
