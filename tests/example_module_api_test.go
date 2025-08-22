package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	example_module "go-template/internal/modules/example-module"
	"go-template/pkg/platform"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func setupTestApp() *fiber.App {
	app := fiber.New()

	// Initialize database
	db, err := platform.InitPostgres()
	if err != nil {
		panic(err)
	}

	// Auto migrate
	err = db.AutoMigrate(&example_module.ExampleDomain{})
	if err != nil {
		panic(err)
	}

	// Initialize module
	module := example_module.NewExampleModule(db)

	// Setup routes
	api := app.Group("/api/v1")
	module.RegisterRoutes(api)

	return app
}

func TestExampleModuleAPI_CreateExample(t *testing.T) {
	app := setupTestApp()

	// Prepare request with unique email
	example := map[string]interface{}{
		"name":  "Test User",
		"email": fmt.Sprintf("test_%d@example.com", time.Now().UnixNano()),
	}

	jsonData, _ := json.Marshal(example)
	req, _ := http.NewRequest("POST", "/api/v1/examples", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	// Perform request
	resp, err := app.Test(req)
	assert.NoError(t, err)

	// Assert response
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var response map[string]interface{}
	body, _ := io.ReadAll(resp.Body)
	err = json.Unmarshal(body, &response)
	assert.NoError(t, err)

	assert.Equal(t, example["name"], response["name"])
	assert.Equal(t, example["email"], response["email"])
	assert.Equal(t, "active", response["status"])
	assert.NotNil(t, response["id"])
}

func TestExampleModuleAPI_GetExample(t *testing.T) {
	app := setupTestApp()

	// First create an example with unique email
	example := map[string]interface{}{
		"name":  "Test User",
		"email": fmt.Sprintf("test_%d@example.com", time.Now().UnixNano()),
	}

	jsonData, _ := json.Marshal(example)
	req, _ := http.NewRequest("POST", "/api/v1/examples", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)

	var createResponse map[string]interface{}
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &createResponse)
	id := createResponse["id"]
	resp.Body.Close() // Close the response body

	// Now get the example
	getReq, _ := http.NewRequest("GET", fmt.Sprintf("/api/v1/examples/%.0f", id), nil)
	getResp, err := app.Test(getReq)
	assert.NoError(t, err)

	// Assert response
	assert.Equal(t, http.StatusOK, getResp.StatusCode)

	var response map[string]interface{}
	getBody, _ := io.ReadAll(getResp.Body)
	err = json.Unmarshal(getBody, &response)
	assert.NoError(t, err)

	assert.Equal(t, example["name"], response["name"])
	assert.Equal(t, example["email"], response["email"])
}
