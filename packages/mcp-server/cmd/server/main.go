package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// HealthHandler responds to health check requests with HTTP 200 and a JSON body indicating the service is up.
// The response JSON contains "message": "codeguard is up.." and "status": "OK".
func HealthHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "codeguard is up..",
		"status":  "OK",
	})
}

// ScanHandler parses the request JSON body into a map and returns it as the response.
// If JSON binding fails it writes a 400 JSON error `{"message":"Not a valid JSON"}` and then still writes a 200 JSON response containing the parsed map (which may be empty or nil).
func ScanHandler(c *gin.Context) {
	var requestData map[string]interface{}

	if err := c.BindJSON(&requestData); err != nil {
		c.JSON(400, gin.H{
			"message": "Not a valid JSON",
		})
	}

	c.JSON(200, requestData)
}

// FixPlansHandler responds to requests with a 200 status and a JSON payload containing `{"message":"Working.."}`.
func FixPlansHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Working..",
	})
}

// EventsHandler serves the contents of the local "events.json" file as application/json.
// 
// If "events.json" is present and readable, it responds with HTTP 200 and the raw file bytes
// with Content-Type "application/json". If the file does not exist it responds with HTTP 404
// and JSON {"message": "No file found at that path"}. If reading the file fails for another
// reason it responds with HTTP 500 and JSON {"message": "Error reading from the file"}.
func EventsHandler(c *gin.Context) {
	jsonBytes, err := os.ReadFile("events.json")
	if err != nil {
		if os.IsNotExist(err) {
			c.JSON(http.StatusNotFound, gin.H{"message": "No file found at that path"})
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error reading from the file"})
	}

	c.Data(200, "application/json", jsonBytes)
}

// main starts the HTTP server, registers application routes, and runs the Gin router.
// It registers handlers for /health, /v1/scan, /v1/fixplans, and /v1/events and uses Gin's default server settings.
func main() {
	router := gin.Default()
	router.GET("/health", HealthHandler)
	router.POST("/v1/scan", ScanHandler)
	router.POST("/v1/fixplans", FixPlansHandler)
	router.GET("/v1/events", EventsHandler)
	router.Run()
}