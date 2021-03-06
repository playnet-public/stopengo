# stopengo

## What is stopengo?

Stopengo is a simple-to-use Steam OpenID login authorizer.
It provides several functions to create and read Steam OpenID requests.

We explicitly do not serve the general Steam Web API because this is only an OpenID authenticator. Stopengo only can gather the SteamID64 from a request.

Although, there are several other Go packages that solve the problem of getting more information about the user, like:
* https://github.com/Philipp15b/go-steamapi

## Get Stopengo

```sh
go get github.com/playnet-public/stopengo
```

## Example

### Single Domain Usage

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
    returnTo, _ := url.Parse("http://localhost:5100/callback?someStatelessValue=AStatelessValue")

    http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
        u, err := stopengo.RedirectURL(realm, returnTo)
        if err != nil {
            panic(err)
        }
        http.Redirect(w, r, u, 301)
    })
    http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
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

### Front- and Backend Usage

The server listening on port 8081 is the example frontend whereas the server listening on port 8088 is the example backend validation server.

```go
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/playnet-public/stopengo"

	"github.com/gorilla/mux"
)

func main() {
	go func() {
		router := mux.NewRouter()

		router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			if err := stopengo.Validate(r); err != nil {
				log.Fatal(err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			steamID64str, err := stopengo.SteamID64(r)
			if err != nil {
				log.Fatal(err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			fmt.Println(steamID64str)

			// Do whatever you want with steam id
			w.Write([]byte(steamID64str))
		})

		http.ListenAndServe(":8088", router)
	}()

	realm, _ := url.Parse("http://localhost")
	returnTo, _ := url.Parse("http://localhost:8081/callback")

	router := mux.NewRouter()
	router.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		redirectURL, _ := stopengo.RedirectURL(realm, returnTo)
		http.Redirect(w, r, redirectURL, 307)
	})
	router.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		nURL, _ := url.Parse("http://localhost:8088")
		nURL.RawQuery = r.URL.Query().Encode()

		resp, err := http.Get(nURL.String())
		if err != nil {
			panic(err)
		}

		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		w.Write(b)
	})
	http.ListenAndServe(":8081", router)
}
```
