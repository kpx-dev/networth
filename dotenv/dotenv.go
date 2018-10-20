package dotenv

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func init() {
	dir, _ := os.Getwd()
	if strings.HasSuffix(dir, "/dbstream") {
		dir = strings.Replace(dir, "/dbstream", "", 1)
	}

	if strings.HasSuffix(dir, "/lib") {
		dir = strings.Replace(dir, "/lib", "", 1)
	}

	if strings.HasSuffix(dir, "/api") {
		dir = strings.Replace(dir, "/api", "", 1)
	}

	if strings.HasSuffix(dir, "/sync") {
		dir = strings.Replace(dir, "/sync", "", 1)
	}

	envPath := fmt.Sprintf("%s/.env", dir)
	file, _ := os.Open(envPath)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") || len(line) == 0 {
			continue
		}
		lineSplitted := strings.Split(line, "=")
		if len(lineSplitted) != 2 {
			continue
		}

		key := strings.TrimSpace(lineSplitted[0])
		val := strings.TrimSpace(lineSplitted[1])
		os.Setenv(key, val)
	}
}
