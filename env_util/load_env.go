package envutil

import (
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/joho/godotenv"
)

// GetGoDotEnvVariable : .env 환경변수 로드
func GetGoDotEnvVariable(key string) string {
	var (
		_, b, _, _     = runtime.Caller(0)
		workingDirPath = filepath.Dir(filepath.Dir(b))
	)

	err := godotenv.Load(workingDirPath + "/.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return os.Getenv(key)
}
