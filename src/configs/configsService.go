package configs

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type EnvConfig struct {
	Host                 string
	Port                 int
	GinMode              string
	DbName               string
	DbHost               string
	DbPort               int
	DbUser               string
	DbPass               string
	Auth0ID              string
	Auth0Secret          string
	Auth0Provider        string
	GGProjectID          string
	GGServiceAccountPath string
	GGCloudStorageBucket string
}

func loadAndValidateEnv() EnvConfig {
	env := os.Getenv("GIN_ENV")
	if env == "" {
		env = "local"
	}

	fmt.Println(env)
	err := godotenv.Load(fmt.Sprintf(".%s.env", env))
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		log.Fatal("Error loading port")
	}

	dbPort, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		log.Fatal("Error loading DB port")
	}

	return EnvConfig{
		Host:                 os.Getenv("HOST"),
		Port:                 port,
		GinMode:              os.Getenv("GIN_MODE"),
		DbName:               os.Getenv("DB_NAME"),
		DbHost:               os.Getenv("DB_HOST"),
		DbPort:               dbPort,
		DbUser:               os.Getenv("DB_USER"),
		DbPass:               os.Getenv("DB_PASS"),
		Auth0ID:              os.Getenv("AUTH0_ID"),
		Auth0Secret:          os.Getenv("AUTH0_SECRET"),
		Auth0Provider:        os.Getenv("AUTH0_PROVIDER"),
		GGProjectID:          os.Getenv("GG_PROJECT_ID"),
		GGServiceAccountPath: os.Getenv("GG_SERVICE_ACCOUNT_PATH"),
		GGCloudStorageBucket: os.Getenv("GG_CLOUDSTORAGE_BUCKET"),
	}
}

// ConfigsService : config constant
var ConfigsService EnvConfig

func init() {
	ConfigsService = loadAndValidateEnv()
}
