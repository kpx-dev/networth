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

	if strings.HasSuffix(dir, "/lib") {
		dir = strings.Replace(dir, "/lib", "", 1)
	}

	if strings.HasSuffix(dir, "/dotenv") {
		dir = strings.Replace(dir, "/dotenv", "", 1)
	}

	envPath := fmt.Sprintf("%s/.env", dir)
	fmt.Println("The .env path is ", envPath)
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
