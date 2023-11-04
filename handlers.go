package main

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Case represents the structure of each case data
type Case struct {
	Province string `json:"province"`
	Age      int    `json:"age"`
}

// SummaryHandler handles the summary endpoint
func SummaryHandler(c *gin.Context) {
	// Fetch data from the API endpoint
	res, err := http.Get("https://static.wongnai.com/devinterview/covid-cases.json")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch data"}) // Handle error if failed to fetch data from the API
		return
	}
	defer res.Body.Close()

	var casesData struct {
		Data []Case `json:"Data"`
	}

	// Parse JSON data from the API response
	err = json.NewDecoder(res.Body).Decode(&casesData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse JSON Data"}) // Handle error if failed to parse JSON data
		return
	}

	// If there is no data, return empty maps for "Province" and "AgeGroup"
	if len(casesData.Data) == 0 {
		summary := gin.H{
			"Province": make(map[string]int),
			"AgeGroup": make(map[string]int),
		}
		c.JSON(http.StatusOK, summary)
		return
	}

	// Count cases by province and age group
	provinceCount := make(map[string]int)
	ageGroupCount := map[string]int{
		"0-30":  0,
		"31-60": 0,
		"61+":   0,
		"N/A":   0,
	}

	for _, elem := range casesData.Data {
		provinceCount[elem.Province]++

		// Categorize cases into age groups
		switch {
		case elem.Age >= 0 && elem.Age <= 30:
			ageGroupCount["0-30"]++
		case elem.Age >= 31 && elem.Age <= 60:
			ageGroupCount["31-60"]++
		case elem.Age >= 61:
			ageGroupCount["61+"]++
		default:
			ageGroupCount["N/A"]++
		}
	}

	// Prepare summary data and respond with a JSON
	summary := gin.H{
		"Province": provinceCount,
		"AgeGroup": ageGroupCount,
	}
	c.JSON(http.StatusOK, summary)
}

// This redudancy should cover some cases handle in fetching + sorting data (N/A is a fall through)