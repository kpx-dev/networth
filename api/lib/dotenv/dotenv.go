package dotenv

import (
	"bufio"
	"os"
	"strings"
)

func getRootDir() string {
	dir, _ := os.Getwd()

	if strings.HasSuffix(dir, "/api/lib/dotenv") {
		dir = strings.Replace(dir, "/api/lib/dotenv", "", 1)
	}

	return dir
}

// LoadDotEnv load .env file. TODO: remove in favor of Lambda ENV
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

func init() {
	loadDotEnv()
}
