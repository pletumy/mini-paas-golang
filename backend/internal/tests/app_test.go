package tests

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestAppIntegration(t *testing.T) {
	//	create
	payload := `{"name":"test-app", "git_url":"https://example.com/repo.git", "description":"intergration_test"}`
	resp, err := http.Post(
		testServer.URL+"/api/apps",
		"application/json",
		strings.NewReader(payload))
	if err != nil {
		t.Fatalf("failed to create app: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("create app expected 201 got %d", resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	var created map[string]interface{}
	json.Unmarshal(body, &created)
	appID := created["id"].(string)

	//	get by id
	resp2, _ := http.Get(testServer.URL + "/api/apps/app/" + appID)
	if resp2.StatusCode != http.StatusOK {
		t.Fatalf("get app by id expected 200 got %d", resp2.StatusCode)
	}
	resp2.Body.Close()

	//	list
	resp3, _ := http.Get(testServer.URL + "/api/apps")
	if resp3.StatusCode != http.StatusOK {
		t.Fatalf("list app expected 200 got %d", resp3.StatusCode)
	}
	resp3.Body.Close()

	//	delete
	req, _ := http.NewRequest(
		http.MethodDelete,
		testServer.URL+"/api/apps/app/"+appID, nil)
	resp4, _ := http.DefaultClient.Do(req)
	if resp4.StatusCode != http.StatusOK {
		t.Fatalf("delete app expected 200 got %d", resp4.StatusCode)
	}
	resp4.Body.Close()
}
