package loadEnv

import (
	"bufio"
	"errors"
	"os"
	"strings"
)

func Load() (err error) {
	var file *os.File
	if file, err = os.Open("./.env"); err != nil {
		return errors.New(".env file missing")
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		pair := strings.SplitN(scanner.Text(), "=", 2)
		if os.Getenv(pair[0]) != "" {
			break
		}
		os.Setenv(pair[0], pair[1])
	}
	if err = scanner.Err(); err != nil {
		return err
	}
	return nil
}
