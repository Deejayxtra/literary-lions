// package config

// const BaseApi = "http://localhost:8080/api/v1.0"

package config

import (
	"os"
)

var (
	BaseApi = os.Getenv("API_URL")
)

func init() {
	if BaseApi == "" {
		BaseApi = "http://localhost:8080/api/v1.0"
	}
}
