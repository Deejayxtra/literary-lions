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
