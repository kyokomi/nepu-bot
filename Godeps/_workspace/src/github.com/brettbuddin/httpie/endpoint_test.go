package httpie

import (
    "testing"
    "net/url"
    "net/http"
)

func testUrl(raw string) *url.URL {
    u, _ := url.Parse(raw)
    return u
}

func TestGet(t *testing.T) {
    end := Get{testUrl("http://google.com")}
    req := &http.Request{}

    end.ApplyTo(req)

    if req.Method != "GET" {
        t.Error("Get did not properly set the Method on the Request")
    }

    if req.URL.Host != "google.com" {
        t.Error("Get did not properly set the Host on the URL")
    }
}

func TestPost(t *testing.T) {
    data        := []byte("thedata")
    contentType := "application/text"
    end         := Post{testUrl("http://google.com"), data, contentType}
    req         := &http.Request{Header: http.Header{}}

    end.ApplyTo(req)

    if req.Method != "POST" {
        t.Error("Post did not properly set the Method on the Request")
    }

    if req.URL.Host != "google.com" {
        t.Error("Post did not properly set the Host on the URL")
    }
}

func TestPut(t *testing.T) {
    data        := []byte("thedata")
    contentType := "application/text"
    end         := Put{testUrl("http://google.com"), data, contentType}
    req         := &http.Request{Header: http.Header{}}

    end.ApplyTo(req)

    if req.Method != "PUT" {
        t.Error("Post did not properly set the Method on the Request")
    }

    if req.URL.Host != "google.com" {
        t.Error("Post did not properly set the Host on the URL")
    }
}

func TestDelete(t *testing.T) {
    end := Delete{testUrl("http://google.com")}
    req := &http.Request{}

    end.ApplyTo(req)

    if req.Method != "DELETE" {
        t.Error("Post did not properly set the Method on the Request")
    }

    if req.URL.Host != "google.com" {
        t.Error("Post did not properly set the Host on the URL")
    }
}
