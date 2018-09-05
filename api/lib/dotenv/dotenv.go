package dotenv

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func init() {
	dir, _ := os.Getwd()
	fmt.Println("the .env root dir is", dir)

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
		fmt.Println("Settin env ", key, val)
		os.Setenv(key, val)
	}
}
