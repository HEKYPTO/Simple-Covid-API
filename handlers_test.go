// handlers_test.go

package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestSummaryHandler(t *testing.T) {
	// Create a new Gin router and set the SummaryHandler as the handler function
	router := gin.Default()
	router.GET("/covid/summary", SummaryHandler)

	// Create a mock request to test the handler function
	req, _ := http.NewRequest("GET", "/covid/summary", nil)
	w := httptest.NewRecorder()

	// Serve the request to the handler
	router.ServeHTTP(w, req)

	// Assert the response status code and body
	assert.Equal(t, http.StatusOK, w.Code)

	// Parse the response JSON body
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Error parsing JSON response: %v", err)
	}

	// Assert the structure of the response JSON
	assert.Contains(t, response, "Province")
	assert.Contains(t, response, "AgeGroup")

}