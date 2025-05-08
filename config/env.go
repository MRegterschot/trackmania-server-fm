package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/MRegterschot/trackmania-server-fm/structs"
	"github.com/joho/godotenv"
)

var AppEnv *structs.Env

func LoadEnv() error {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Failed to load .env file, using default values")
	}

	port, err := strconv.Atoi(os.Getenv("FM_PORT"))
	if err != nil {
		port = 3300
	}

	userDataPath := os.Getenv("FM_USERDATA_PATH")
	if userDataPath == "" {
		userDataPath = "/app/UserData"
	}

	absPath, err := filepath.Abs(userDataPath)
	if err != nil {
		return errors.New("failed to get absolute path for UserData directory")
	}

	AppEnv = &structs.Env{
		Port:         port,
		LogLevel:     os.Getenv("FM_LOG_LEVEL"),
		UserDataPath: absPath,
	}

	return nil
}
