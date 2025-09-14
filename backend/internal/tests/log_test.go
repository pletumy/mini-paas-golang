package tests

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"
)

func TestLogIntegration(t *testing.T) {
	// create app
	appPayload := `{"name":"log-app", "git_url":"https://example.com/repo.git"}`
	appResp, err := http.Post(
		testServer.URL+"/api/apps",
		"application/json",
		strings.NewReader(appPayload))
	if err != nil {
		t.Fatal(err)
	}
	defer appResp.Body.Close()
	if appResp.StatusCode != http.StatusCreated {
		t.Fatalf("create app expected 201 got %d", appResp.StatusCode)
	}
	var appCreated map[string]interface{}
	json.NewDecoder(appResp.Body).Decode(&appCreated)
	appID := appCreated["id"].(string)

	// create deployment
	depPayload := `{"app_id":"` + appID + `", "version":"v1"}`
	depResp, err := http.Post(
		testServer.URL+"/api/deployments",
		"application/json",
		strings.NewReader(depPayload))
	if err != nil {
		t.Fatal(err)
	}
	defer depResp.Body.Close()
	if depResp.StatusCode != http.StatusCreated {
		t.Fatalf("create deployment expected 201 got %d", depResp.StatusCode)
	}
	var depCreated map[string]interface{}
	json.NewDecoder(depResp.Body).Decode(&depCreated)
	deploymentID := depCreated["id"].(string)

	// create log
	payload := `{"deployment_id":"` + deploymentID + `", "message":"hello world"}`
	resp, err := http.Post(
		testServer.URL+"/api/logs",
		"application/json",
		strings.NewReader(payload))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("create log expected 201 got %d", resp.StatusCode)
	}

	// list
	resp2, _ := http.Get(testServer.URL + "/api/logs?deployment_id=" + deploymentID)
	if resp2.StatusCode != http.StatusOK {
		t.Fatalf("list logs expected 200 got %d", resp2.StatusCode)
	}
	resp2.Body.Close()
}
