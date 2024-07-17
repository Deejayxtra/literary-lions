// package docs

// // @title Literary Lions API
// // @version 1.0
// // @description API documentation for Literary Lions backend

// // @host localhost:8080
// // @BasePath /
// func SwaggerInfo() {
// 	swagger.SwaggerInfo.Title = "Literary Lions API"
// 	swagger.SwaggerInfo.Version = "1.0"
// 	swagger.SwaggerInfo.Description = "API documentation for Literary Lions backend"
// 	swagger.SwaggerInfo.Host = "localhost:8080"
// 	swagger.SwaggerInfo.BasePath = "/"
// }

// package docs

// import "github.com/swaggo/swag"

// // @title Literary Lions API
// // @version 1.0
// // @description API documentation for Literary Lions backend

// // @host localhost:8080
// // @BasePath /
// func init() {
// 	swag.SwaggerInfo.Title = "Literary Lions API"
// 	swag.SwaggerInfo.Version = "1.0"
// 	swag.SwaggerInfo.Description = "API documentation for Literary Lions backend"
// 	swag.SwaggerInfo.Host = "localhost:8080"
// 	swag.SwaggerInfo.BasePath = "/"
// }

// package docs

// import (
// 	"github.com/swaggo/swag"
// )

// func init() {
// 	swag.Register(swag.Name, &swaggerDocs{})
// }

// type swaggerDocs struct{}

// func (s *swaggerDocs) ReadDoc() string {
// 	// Your Swagger JSON content here
// 	return ``
// }

package docs

import (
	"github.com/swaggo/swag"
)

func init() {
	swag.Register(swag.Name, &swaggerDocs{})
}

type swaggerDocs struct{}

func (s *swaggerDocs) ReadDoc() string {
	// Your Swagger JSON content here
	return ``
}
