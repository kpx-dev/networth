package dotenv

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func init() {
	dir, _ := os.Getwd()
	if strings.HasSuffix(dir, "/api/dbstream") {
		dir = strings.Replace(dir, "/api/dbstream", "", 1)
	}

	if strings.HasSuffix(dir, "/api/lib") {
		dir = strings.Replace(dir, "/api/lib", "", 1)
	}

	if strings.HasSuffix(dir, "/api") {
		dir = strings.Replace(dir, "/api", "", 1)
	}

	envPath := fmt.Sprintf("%s/.env", dir)
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
