## HTTPie

A small HTTP library for Go. It helps with normal HTTP requests and HTTP streaming.

### Example Usage

```go
package main

import (
    "fmt"
    "net/url"
    "github.com/brettbuddin/httpie"
)

func main() {
    // Create a GET endpoint
    endpoint := httpie.Get{&url.URL{
        Scheme: "http",
        Host: "google.com",
        Path: "/robots.txt",
    }}

    // We could alternatively, if we don't need auth, set this to nil
    authorizer := httpie.BasicAuth{"username", "password"}

    // Do the business
    client := httpie.NewClient(authorizer)
    resp, err := client.Request(endpoint)

    if err != nil {
        fmt.Printf("We have an error: %s\n", err)
        return
    }

    fmt.Printf("Our response: %s\n", resp)
}
```
