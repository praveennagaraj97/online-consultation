package env

import (
	"os"

	"github.com/joho/godotenv"
	logger "github.com/praveennagaraj97/online-consultation/pkg/log"
)

func init() {

	err := godotenv.Load(".env")
	if err != nil {
		logger.PrintLog("Failed to load env variables")
	}

	logger.PrintLog("Environment Variables Loaded ⚙️ ")
}

func GetEnvVariable(key string) string {
	return os.Getenv(key)
}
