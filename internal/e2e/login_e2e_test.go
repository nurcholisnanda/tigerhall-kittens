package e2e

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

var Token string

func TestLoginUser(t *testing.T) {
	// Prepare the input data for user login
	email := "nurcholis.nanda@gmail.com"
	password := "123456"

	// Create the GraphQL mutation payload
	mutation := `
        mutation ($email: String!,$password: String!){
            auth {
                login(email: $email, password: $password)
            }
        }
    `

	// Prepare the variables object with the input data
	variables := map[string]interface{}{
		"email":    email,
		"password": password,
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

	// Extract the login result from the response
	data, ok := result["data"].(map[string]interface{})
	if !ok {
		t.Fatal("data field not found in GraphQL response")
	}

	auth, ok := data["auth"].(map[string]interface{})
	if !ok {
		t.Fatal("auth field not found in GraphQL response")
	}

	login, ok := auth["login"].(map[string]interface{})
	if !ok {
		t.Fatal("login field not found in GraphQL response")
	}

	token, ok := login["token"].(string)
	if !ok {
		t.Fatal("token not generated in GraphQL response")
	}

	Token = token
	// Assert that the login was successful
	assert.NotNil(t, token)
}
