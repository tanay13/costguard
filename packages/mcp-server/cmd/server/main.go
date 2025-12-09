package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func HealthHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "codeguard is up..",
		"status":  "OK",
	})
}

func ScanHandler(c *gin.Context) {
	var requestData map[string]interface{}

	if err := c.BindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Not a valid JSON",
		})
		return
	}

	c.JSON(http.StatusOK, requestData)
}

func FixPlansHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Working..",
	})
}

func EventsHandler(c *gin.Context) {
	jsonBytes, err := os.ReadFile("events.json")
	if err != nil {
		if os.IsNotExist(err) {
			c.JSON(http.StatusNotFound, gin.H{"message": "No file found at that path"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error reading from the file"})
		return
	}

	c.Data(200, "application/json", jsonBytes)
}

func main() {
	router := gin.Default()
	router.GET("/health", HealthHandler)
	router.POST("/v1/scan", ScanHandler)
	router.POST("/v1/fixplans", FixPlansHandler)
	router.GET("/v1/events", EventsHandler)
	router.Run()
}
