package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/covid/summary", SummaryHandler)
	router.Run(":8080")
}