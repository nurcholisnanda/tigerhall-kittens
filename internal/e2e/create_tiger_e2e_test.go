package e2e

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/nurcholisnanda/tigerhall-kittens/internal/api/graph/model"
	"github.com/stretchr/testify/assert"
)

var TigerID string

func TestCreateTiger(t *testing.T) {
	TestLoginUser(t)
	// Prepare the input data for creating a tiger
	tigerInput := model.TigerInput{
		Name:         "tiger j",
		DateOfBirth:  time.Date(2018, 06, 25, 15, 0, 0, 0, time.UTC),
		LastSeenTime: time.Now().Add(time.Hour * -10),
		LastSeenCoordinate: &model.CoordinateInput{
			Latitude:  -6.175,
			Longitude: 95.55,
		},
	}

	// Create the GraphQL mutation payload
	mutation := `
        mutation ($input: TigerInput!) {
            create {
                createTiger(input: $input) {
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

	// Prepare the variables object with the input data
	variables := map[string]interface{}{
		"input": tigerInput,
	}

	// Marshal the mutation and variables into a JSON payload
	payload := map[string]interface{}{
		"query":     mutation,
		"variables": variables,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("error marshalling JSON payload: %v", err)
	}

	// Create an HTTP client
	client := &http.Client{}

	// Send the HTTP POST request to the GraphQL server
	req, err := http.NewRequest("POST", "http://localhost:8081/query", bytes.NewBuffer(payloadBytes))
	if err != nil {
		t.Fatalf("error creating HTTP request: %v", err)
	}

	// Add Authorization header with bearer token
	req.Header.Set("Authorization", "Bearer "+Token)
	req.Header.Set("Content-Type", "application/json")

	// Send the HTTP request using the client
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("HTTP POST request failed: %v", err)
	}
	defer resp.Body.Close()

	// Decode the response JSON
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("error decoding JSON response: %v", err)
	}

	// Extract the created tiger details from the response
	data, ok := result["data"].(map[string]interface{})
	if !ok {
		t.Fatalf("data field not found in GraphQL response : %v", result)
	}

	create, ok := data["create"].(map[string]interface{})
	if !ok {
		t.Fatal("create field not found in GraphQL response")
	}

	createdTiger, ok := create["createTiger"].(map[string]interface{})
	if !ok {
		t.Fatal("createTiger field not found in GraphQL response")
	}
	TigerID = createdTiger["id"].(string)
	dateOfBirth := tigerInput.DateOfBirth.Format("2006-01-02T15:04:05Z")

	// Validate the created tiger's details
	assert.NotNil(t, createdTiger["id"], "id field not found in created tiger")
	assert.Equal(t, tigerInput.Name, createdTiger["name"], "name mismatch in created tiger")
	assert.Equal(t, dateOfBirth, createdTiger["dateOfBirth"], "dateOfBirth mismatch in created tiger")

	// Validate lastSeenCoordinate
	lastSeenCoord, ok := createdTiger["lastSeenCoordinate"].(map[string]interface{})
	if !ok {
		t.Fatal("lastSeenCoordinate field not found or not a map in created tiger")
	}
	assert.Equal(t, tigerInput.LastSeenCoordinate.Latitude, lastSeenCoord["latitude"], "latitude mismatch in LastSeenCoordinate")
	assert.Equal(t, tigerInput.LastSeenCoordinate.Longitude, lastSeenCoord["longitude"], "longitude mismatch in lastSeenCoordinate")
}
