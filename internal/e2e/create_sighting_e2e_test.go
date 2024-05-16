package e2e

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateSighting(t *testing.T) {
	TestCreateTiger(t)
	// Open the image file to upload
	imageFile, err := os.Open("./image.jpeg")
	if err != nil {
		t.Fatalf("error opening image file: %v", err)
	}
	defer imageFile.Close()

	// Create a new HTTP request
	var requestBody bytes.Buffer
	multipartWriter := multipart.NewWriter(&requestBody)

	// Write GraphQL operations and variables
	operations := `{
        "query": "mutation($file: Upload!, $tigerID: String!) { create { createSighting(input: { tigerID: $tigerID, Coordinate: { latitude: 12.2, longitude: 120 }, image: $file }) { id lastSeenTime lastSeenCoordinate { latitude longitude } image } } }",
        "variables": { "tigerID": "` + TigerID + `", "file": null }
    }`
	_ = multipartWriter.WriteField("operations", operations)

	// Write GraphQL map for file upload
	_ = multipartWriter.WriteField("map", `{"0": ["variables.file"]}`)

	// Add image file to multipart request
	imageWriter, err := multipartWriter.CreateFormFile("0", "image.jpg")
	if err != nil {
		t.Fatalf("error creating form file: %v", err)
	}
	_, err = io.Copy(imageWriter, imageFile)
	if err != nil {
		t.Fatalf("error copying image file: %v", err)
	}

	// Close multipart writer and set Content-Type header
	multipartWriter.Close()

	// Create HTTP request with multipart form data
	req, err := http.NewRequest("POST", "http://localhost:8081/query", &requestBody)
	if err != nil {
		t.Fatalf("error creating HTTP request: %v", err)
	}
	req.Header.Set("Content-Type", multipartWriter.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+Token)

	// Create an HTTP client
	client := &http.Client{}

	// Send the HTTP request
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("HTTP POST request failed: %v", err)
	}
	defer resp.Body.Close()

	// Check the response status code
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Read and print the response body for debugging
	var responseBytes bytes.Buffer
	_, _ = io.Copy(&responseBytes, resp.Body)
	fmt.Println("GraphQL Response:", responseBytes.String())

	// Optionally, parse and validate the GraphQL response as needed
	// Example: Use JSON unmarshaling to extract relevant data from the response
}
