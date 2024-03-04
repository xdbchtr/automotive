package controllers

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
	"io"
	"log"
	"net/http"
	"strconv"
)

type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func Upload(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.String(http.StatusBadRequest, "Bad Request")
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		c.String(http.StatusInternalServerError, "Error opening file")
		return
	}

	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("Error closing file: %v", err)
		}
	}()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error Reading File")
		return
	}

	f, err := excelize.OpenReader(bytes.NewReader(fileBytes))
	if err != nil {
		c.String(http.StatusInternalServerError, "Error opening file content")
		return
	}

	rows, err := f.GetRows("Sheet1")
	if err != nil {
		c.String(http.StatusInternalServerError, "Error reading file content")
		return
	}

	var people []Person

	for _, row := range rows {
		if len(row) >= 2 {
			age, err := strconv.Atoi(row[1])
			if err != nil {
				log.Printf("Error converting age: %v", err)
				continue
			}
			person := Person{
				Name: row[0],
				Age:  age,
			}
			people = append(people, person)
		}
	}

	jsonData, err := json.Marshal(people)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error marshalling JSON")
		return
	}

	c.Data(http.StatusOK, "application/json", jsonData)
}
