package swagger

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

//go:embed static\swagger-ui\*
var swaggerUI embed.FS

//go:embed static\swagger-ui\documentation.yaml
var openAPISpec []byte

// SwaggerUIHandler serves the Swagger UI static assets.
func SwaggerUIHandler(c *gin.Context) {
	subFS, err := fs.Sub(swaggerUI, "static/swagger-ui")
	if err != nil {
		log.Println("the embed path is not correct. \n", err.Error())
		c.String(http.StatusInternalServerError, "Failed to load Swagger UI")
		return
	}

	fileServer := http.FileServer(http.FS(subFS))

	c.Request.URL.Path = strings.TrimPrefix(c.Request.URL.Path, "/swagger")

	// Serve the request.
	fileServer.ServeHTTP(c.Writer, c.Request)
}

// OpenAPISpecHandler serves the raw openapi.yaml file.
func OpenAPISpecHandler(c *gin.Context) {
	c.Data(http.StatusOK, "application/x-yaml", openAPISpec)
}
