package httpie

import (
    "net/url"
    "net/http"
    "io/ioutil"
    "bytes"
)

// Endpoint is implemented to provide delayed
// attachment of the URL/Method/Body of a request
type Endpoint interface {
    ApplyTo(*http.Request)
}

// Represents an HTTP GET
type Get struct {
    *url.URL
}

// ApplyTo sets the requests Method to GET and URL
func (g Get) ApplyTo(req *http.Request) {
    req.Method = "GET"
    req.URL    = g.URL
}

// Represents an HTTP POST
type Post struct {
    *url.URL
    Body []byte
    ContentType string
}

// ApplyTo sets the requests Method to POST, URL and Body
func (p Post) ApplyTo(req *http.Request) {
    req.Method        = "POST"
    req.URL           = p.URL
    req.Body          = ioutil.NopCloser(bytes.NewBuffer(p.Body))
    req.ContentLength = int64(len(p.Body))
    req.Header.Set("Content-Type", p.ContentType)
}

// Represents an HTTP PUT
type Put struct {
    *url.URL
    Body []byte
    ContentType string
}

// ApplyTo sets the requests Method to PUT, URL and Body
func (p Put) ApplyTo(req *http.Request) {
    req.Method        = "PUT"
    req.URL           = p.URL
    req.Body          = ioutil.NopCloser(bytes.NewBuffer(p.Body))
    req.ContentLength = int64(len(p.Body))
    req.Header.Set("Content-Type", p.ContentType)
}

// Represents an HTTP DELETE
type Delete struct {
    *url.URL
}

// ApplyTo sets the requests Method to DELETE and URL
func (d Delete) ApplyTo(req *http.Request) {
    req.Method = "DELETE"
    req.URL    = d.URL
}
