package main

import (
	"encoding/json"
	"net/http"
	"github.com/gin-gonic/gin"
)

type Case struct {
    Province string `json:"province"`
    Age      int    `json:"age"`
}

func SummaryHandler(c *gin.Context) {
	res, err := http.Get("https://static.wongnai.com/devinterview/covid-cases.json")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch data"})
		return
	}
	defer res.Body.Close()

	var casesData struct {
		Data []Case `json:"Data"`
	}

	err = json.NewDecoder(res.Body).Decode(&casesData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse JSON Data"})
		return
	}

	if len(casesData.Data) == 0 {
		// If there is no data, return empty maps for "Province" and "AgeGroup"
		summary := gin.H{
			"Province": make(map[string]int),
			"AgeGroup": make(map[string]int),
		}
		c.JSON(http.StatusOK, summary)
		return
	}

	provinceCount := make(map[string]int)
	ageGroupCount := map[string]int{
		"0-30": 0,
		"31-60": 0,
		"61+": 0,
		"N/A": 0,
	}

	for _, elem := range casesData.Data {
		provinceCount[elem.Province]++
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

	summary := gin.H{
		"Province": provinceCount,
		"AgeGroup": ageGroupCount,
	}

	c.JSON(http.StatusOK, summary)
}