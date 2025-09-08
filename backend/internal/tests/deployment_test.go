package tests

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"testing"
)

func TestDeploymentIntegration(t *testing.T) {
	// create app prerequisite
	appPayload := `{"name":"dep-app", "git_url":"https://example.com/repo.git"}`
	appResp, err := http.Post(
		testServer.URL+"/api/apps",
		"application/json",
		strings.NewReader(appPayload),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer appResp.Body.Close()
	if appResp.StatusCode != http.StatusCreated {
		t.Fatalf("create app expected 201 got %d", appResp.StatusCode)
	}
	var appCreated map[string]interface{}
	json.NewDecoder(appResp.Body).Decode(&appCreated)
	appID := appCreated["id"].(string)

	//	create deployment
	payload := `{"app_id":"` + appID + `", "version":"v1"}`
	resp, err := http.Post(
		testServer.URL+"/api/deployments",
		"application/json",
		strings.NewReader(payload),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("create deployment expected 201 got %d", resp.StatusCode)
	}

	var created map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&created)
	deployID := created["id"].(string)

	// get by id
	resp2, _ := http.Get(
		testServer.URL + "/api/deployments/" + deployID,
	)
	if resp2.StatusCode != http.StatusOK {
		t.Fatalf("get deployment by id expected 200 got %d", resp2.StatusCode)
	}
	resp2.Body.Close()

	// list
	resp3, _ := http.Get(testServer.URL + "/api/deployments")
	if resp3.StatusCode != http.StatusOK {
		t.Fatalf("list deployments expected 200 got %d", resp3.StatusCode)
	}
	resp3.Body.Close()

}
