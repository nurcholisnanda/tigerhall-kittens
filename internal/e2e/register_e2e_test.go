package e2e

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/nurcholisnanda/tigerhall-kittens/internal/api/graph/model"
	"github.com/stretchr/testify/assert"
)

func TestRegisterUser(t *testing.T) {
	// Prepare the input data for user registration
	newUser := model.NewUser{
		Name:     "nucholis",
		Email:    "nurcholis.nanda@gmail.com",
		Password: "123456",
	}

	// Create the GraphQL mutation payload
	mutation := `
        mutation ($input: NewUser!) {
            auth {
                register(input: $input)
            }
        }
    `

	// Prepare the variables object with the input data
	variables := map[string]interface{}{
		"input": newUser,
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

	// Send the HTTP POST request to the GraphQL server
	resp, err := http.Post("http://localhost:8081/query", "application/json", bytes.NewBuffer(payloadBytes))
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

	// Assert that the registration was successful
	data := result["data"].(map[string]interface{})
	auth := data["auth"].(map[string]interface{})
	res := auth["register"].(map[string]interface{})

	assert.Equal(t, newUser.Name, res["name"])
	assert.Equal(t, newUser.Email, res["email"])
}
