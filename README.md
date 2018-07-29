# stopengo

## What is stopengo?

Stopengo is a simple-to-use Steam OpenID login authorizer.
It provides several functions to create and read Steam OpenID requests.

We explicitly do not serve the general Steam Web API because this is only an OpenID authenticator. Stopengo only can gather the SteamID64 from a request.

Although, there are several other go packages that solve the problem of getting more information about the user, like:
* https://github.com/Philipp15b/go-steamapi

## Get Stopengo

```sh
go get github.com/playnet-public/stopengo
```

## Example

```go
package main

import (
    "fmt"
    "net/http"
    "net/url"

    "github.com/playnet-public/stopengo"
)

func main() {
    realm, _ := url.Parse("http://localhost:5100/")
    returnTo, _ := url.Parse("http://localhost:5100/redirect?someStatelessValue=AStatelessValue")

    http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
        u, err := stopengo.RedirectURL(realm, returnTo)
        if err != nil {
            panic(err)
        }
        http.Redirect(w, r, u, 301)
    })
    http.HandleFunc("/redirect", func(w http.ResponseWriter, r *http.Request) {
        if err := stopengo.Validate(r); err != nil {
            panic(err)
        }

        steamid64, err := stopengo.SteamID64(r)
        if err != nil {
            panic(err)
        }

        fmt.Println(steamid64)
    })
    http.ListenAndServe(":5100", nil)
}
```
