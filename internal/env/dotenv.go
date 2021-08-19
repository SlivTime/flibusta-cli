package env

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"path/filepath"
)

const dotenvFileName = ".env"

func Load() {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal("Getwd failed", err)
	}

	envFile := filepath.Join(wd, dotenvFileName)
	err = godotenv.Load(envFile)
	if err != nil {
		log.Println(fmt.Sprintf("loading %s file, %s", envFile, err))
	}
}
