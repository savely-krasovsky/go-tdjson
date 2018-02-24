# Golang bindings for TDLib (Telegram MTProto library)
[![GoDoc](https://godoc.org/github.com/L11R/go-tdjson?status.svg)](https://godoc.org/github.com/L11R/go-tdjson)

In addition to the built-in methods:
- Destroy()
- Execute()
- Receive()
- Send()
- SetFilePath()
- SetLogVerbosityLevel()

It also has two interesting methods:
- Auth()
- SendAndCatch()

# Linking statically against TDLib
I recommend you to link it statically if you don't want compile TDLib on production (don't forget that it requires at least 8GB of RAM). 
<br/>To do that, just build your source with tag `tdjson_static`: `go build -tags tdjson_static`
<br />For more details read this issue: https://github.com/tdlib/td/issues/8

# Example
```golang
package main

import (
	"github.com/L11R/go-tdjson"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	tdjson.SetLogVerbosityLevel(1)
	tdjson.SetFilePath("./errors.txt")

	// Get API_ID and API_HASH from env vars
	apiId := os.Getenv("API_ID")
	if apiId == "" {
		log.Fatal("API_ID env variable not specified")
	}
	apiHash := os.Getenv("API_HASH")
	if apiHash == "" {
		log.Fatal("API_HASH env variable not specified")
	}

	// Create new instance of client
	client := tdjson.NewClient()

	// Handle Ctrl+C
	ch := make(chan os.Signal, 2)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-ch
		client.Destroy()
		os.Exit(1)
	}()

	// Main loop
	for update := range client.Updates {
		// Show all updates
		fmt.Println(update)

		// Authorization block
		if update["@type"].(string) == "updateAuthorizationState" {
			if authorizationState, ok := update["authorization_state"].(tdjson.Update)["@type"].(string); ok {
				res, err := client.Auth(authorizationState, apiId, apiHash)
				if err != nil {
					log.Println(err)
				}
				log.Println(res)
			}
		}
	}
}

```
