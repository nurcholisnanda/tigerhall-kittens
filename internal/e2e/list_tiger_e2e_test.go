package e2e

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListTigers(t *testing.T) {
	// Prepare the GraphQL query to fetch a list of tigers
	query := `
        query ($limit: Int!, $offset: Int!) {
            list {
                listTigers(limit: $limit, offset: $offset) {
                    id
                    name
                    dateOfBirth
                    lastSeenTime
                    lastSeenCoordinate {
                        latitude
                        longitude
                    }
                }
            }
        }
    `

	// Prepare the query variables (limit and offset)
	variables := map[string]interface{}{
		"limit":  10,
		"offset": 0,
	}

	// Marshal the query and variables into a JSON payload
	payload := map[string]interface{}{
		"query":     query,
		"variables": variables,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("error marshalling JSON payload: %v", err)
	}

	// Create an HTTP client
	client := &http.Client{}

	// Create an HTTP request
	req, err := http.NewRequest("POST", "http://localhost:8081/query", bytes.NewBuffer(payloadBytes))
	if err != nil {
		t.Fatalf("error creating HTTP request: %v", err)
	}

	// Add necessary headers (e.g., Content-Type)
	req.Header.Set("Content-Type", "application/json")

	// Send the HTTP request using the client
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("HTTP POST request failed: %v", err)
	}
	defer resp.Body.Close()

	// Check the response status code
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Decode the response JSON
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("error decoding JSON response: %v", err)
	}

	// Extract the list of tigers from the response
	data, ok := result["data"].(map[string]interface{})
	if !ok {
		t.Fatal("data field not found in GraphQL response")
	}

	list, ok := data["list"].(map[string]interface{})
	if !ok {
		t.Fatal("list field not found in GraphQL response")
	}

	listTigers, ok := list["listTigers"].([]interface{})
	if !ok {
		t.Fatal("listTigers field not found or not a list in GraphQL response")
	}
	assert.NotNil(t, listTigers)
}
