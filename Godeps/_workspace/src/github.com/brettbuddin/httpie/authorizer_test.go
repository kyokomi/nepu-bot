package httpie

import (
	"net/http"
	"testing"
)

func TestBasicAuth(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://localhost", nil)
	auth := BasicAuth{"user", "pass"}
	auth.Authorize(req)
	if _, ok := req.Header["Authorization"]; !ok {
		t.Error("BasicAuth did not set Authorization header")
	}
}

func TestHeaderAuth(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://localhost", nil)
	auth := HeaderAuth{"authtokenkeywhatever"}
	auth.Authorize(req)
	if _, ok := req.Header["Authorization"]; !ok {
		t.Error("HeaderAuth did not set Authorization header")
	}
}
