package tests

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func TestUserIntegration(t *testing.T) {
	// generate unique email per run
	suffix := rand.Int63()
	email := fmt.Sprintf("test+%d@gmail.com", suffix)

	//	create
	payload := fmt.Sprintf(`{"email":"%s", "name":"john doe"}`, email)
	resp, err := http.Post(
		testServer.URL+"/api/users",
		"application/json",
		strings.NewReader(payload),
	)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("create user expected 201 got: %d", resp.StatusCode)
	}
	var created map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&created)
	userID := created["id"].(string)

	//	get by id
	resp2, _ := http.Get(
		testServer.URL + "/api/users/user/" + userID,
	)
	if resp2.StatusCode != http.StatusOK {
		t.Fatalf("get user by id expected 200 got %d", resp2.StatusCode)
	}
	resp2.Body.Close()

	// get by email
	u := testServer.URL + "/api/users/user/email?email=" + url.QueryEscape(email)
	resp3, _ := http.Get(u)
	if resp3.StatusCode != http.StatusOK {
		t.Fatalf("get user by email expected 200 got %d", resp3.StatusCode)
	}
	resp3.Body.Close()
}
