package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// DocsHandler handles API documentation
type DocsHandler struct{}

// NewDocsHandler creates a new DocsHandler
func NewDocsHandler() *DocsHandler {
	return &DocsHandler{}
}

// Scalar serves the Scalar API documentation UI
func (h *DocsHandler) Scalar(c *gin.Context) {
	html := `<!DOCTYPE html>
<html>
<head>
    <title>Podoru API Documentation</title>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <style>
        body {
            margin: 0;
            padding: 0;
        }
    </style>
</head>
<body>
    <script
        id="api-reference"
        data-url="/api/v1/docs/openapi.json"
        data-configuration='{
            "theme": "purple",
            "layout": "modern",
            "searchHotKey": "k"
        }'
    ></script>
    <script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
</body>
</html>`
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, html)
}

// OpenAPISpec serves the OpenAPI specification JSON
func (h *DocsHandler) OpenAPISpec(c *gin.Context) {
	c.File("./docs/swagger.json")
}
