package config

import (
	"os"
	"strconv"
)

var folderPath string
var port int

func getEnvironmentVariables() {
	// Get value from environment variable called FOLDER_PATH
	folderPath = os.Getenv("FOLDER_PATH")

	// Parse port as a uint16 from the environment variable called PORT
	// If the environment variable is not set, use the default value of 3131
	retrievedPort, error := strconv.Atoi(os.Getenv("PORT"))
	if error != nil {
		port = 3131
	} else {
		port = retrievedPort
	}
}

func init() {
	getEnvironmentVariables()
}

func FolderPath() string {
	return folderPath
}

func Port() int {
	return port
}
