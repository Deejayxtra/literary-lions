// Package docs contains auto-generated Swagger API documentation.
// To generate or update the documentation, run `swag init` in the project root .
// cd ~/literary-lions/backend... "swag init -g src/cmd/main.go"

package docs

import "github.com/swaggo/swag"

// SwaggerInfo holds exported Swagger Info so we can programmatically access it from Go code.
var SwaggerInfo = &swag.Spec{
	Version:     "1.0",
	Host:        "localhost:8080",
	BasePath:    "/",
	Schemes:     []string{"http"},
	Title:       "API Documentation",
	Description: "This is the API documentation for the Gin API.",
}

// ReadDoc reads the Swagger document.
func ReadDoc() string {
	return `{
		"swagger": "2.0",
		"info": {
			"description": "This is the API documentation for the Gin API.",
			"version": "1.0",
			"title": "API Documentation",
			"contact": {
				"email": "support@example.com"
			},
			"license": {
				"name": "Apache 2.0",
				"url": "http://www.apache.org/licenses/LICENSE-2.0.html"
			}
		},
		"host": "localhost:8080",
		"basePath": "/",
		"schemes": [
			"http"
		],
		"paths": {}
	}`
}
