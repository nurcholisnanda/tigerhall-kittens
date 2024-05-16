package e2e

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListSightings(t *testing.T) {
	TestCreateSighting(t)
	// Prepare GraphQL query with variables
	query := `
        query ListSightings($tigerID: String!) {
            list {
                listSightings(tigerID: $tigerID, limit: 10, offset: 0) {
                    id
                    tigerID
                    lastSeenTime
                    lastSeenCoordinate {
                        latitude
                        longitude
                    }
                }
            }
        }
    `
	variables := map[string]interface{}{
		"tigerID": TigerID, // Replace with actual tigerID
	}

	// Encode the GraphQL query and variables
	requestBody, err := json.Marshal(map[string]interface{}{
		"query":     query,
		"variables": variables,
	})
	if err != nil {
		t.Fatalf("error encoding GraphQL request: %v", err)
	}

	// Create HTTP request to GraphQL endpoint
	req, err := http.NewRequest("POST", "http://localhost:8081/query", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatalf("error creating HTTP request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Send HTTP request and handle response
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("HTTP POST request failed: %v", err)
	}
	defer resp.Body.Close()

	// Check response status code
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Decode GraphQL response
	var responseBody map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	if err != nil {
		t.Fatalf("error decoding JSON response: %v", err)
	}

	// Extract and validate listSightings data from response
	data, ok := responseBody["data"].(map[string]interface{})
	if !ok {
		t.Fatal("data field not found in GraphQL response")
	}

	list, ok := data["list"].(map[string]interface{})
	if !ok {
		t.Fatal("list field not found in GraphQL response")
	}

	listSightings, ok := list["listSightings"].([]interface{})
	if !ok {
		t.Fatal("listSightings field not found in GraphQL response")
	}

	assert.NotEmpty(t, listSightings, "expected non-empty list of sightings")
}
