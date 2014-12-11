package httpie

import (
    "net/http"
)

// Authorizer is implemented to allow delayed 
// attachment of auth data to a request
type Authorizer interface {
    Authorize(*http.Request)
}

type BasicAuth struct {
    Username, Password string
}

// Authorize applys BasicAuth to a request
func (b BasicAuth) Authorize(req *http.Request) {
    req.SetBasicAuth(b.Username, b.Password)
}

type HeaderAuth struct {
    Auth string
}

// Authorize applys HeaderAuth to a request
func (h HeaderAuth) Authorize(req *http.Request) {
    req.Header.Set("Authorization", h.Auth)
}
