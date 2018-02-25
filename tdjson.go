package tdjson

//#cgo linux CFLAGS: -I/usr/local/include
//#cgo darwin CFLAGS: -I/usr/local/include
//#cgo windows CFLAGS: -IC:/src/td -IC:/src/td/build
//#cgo linux,!tdjson_static LDFLAGS: -L/usr/local/lib -ltdjson
//#cgo linux,tdjson_static LDFLAGS: -L/usr/local/lib -ltdjson_static -ltdjson_private -ltdclient -ltdcore -ltdactor -ltddb -ltdsqlite -ltdnet -ltdutils -lstdc++ -lssl -lcrypto -ldl -lz -lm
//#cgo darwin LDFLAGS: -L/usr/local/lib -L/usr/local/opt/openssl/lib -ltdjson_static -ltdjson_private -ltdclient -ltdcore -ltdactor -ltddb -ltdsqlite -ltdnet -ltdutils -lstdc++ -lssl -lcrypto -ldl -lz -lm
//#cgo windows LDFLAGS: -LC:/src/td/build/Debug -ltdjson
//#include <stdlib.h>
//#include <td/telegram/td_json_client.h>
//#include <td/telegram/td_log.h>
import "C"

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"
	"unsafe"
)

type Update = map[string]interface{}

type Client struct {
	Client     unsafe.Pointer
	Updates    chan Update
	waiters    sync.Map
	parameters *options
}

// Creates a new instance of TDLib.
// Has two public fields:
// Client itself and Updates channel
func NewClient(params ...Option) *Client {
	// Seed rand with time
	rand.Seed(time.Now().UnixNano())

	client := Client{Client: C.td_json_client_create()}
	client.Updates = make(chan Update, 100)
	client.parameters = &options{
		systemLanguageCode: "en",
		systemVersion:      "Unknown",
		applicationVersion: "1.0",
	}

	for _, option := range params {
		option(client.parameters)
	}

	go func() {
		for {
			// get update
			update := client.Receive(10)

			// does new update has @extra field?
			if extra, hasExtra := update["@extra"].(string); hasExtra {
				// trying to load update with this salt
				if waiter, found := client.waiters.Load(extra); found {
					// found? send it to waiter channel
					waiter.(chan Update) <- update

					// trying to prevent memory leak
					close(waiter.(chan Update))
				}
			} else {
				// does new updates has @type field?
				if _, hasType := update["@type"]; hasType {
					// if yes, send it in main channel
					client.Updates <- update
				}
			}
		}
	}()

	return &client
}

// Destroys the TDLib client instance.
// After this is called the client instance shouldn't be used anymore.
func (c *Client) Destroy() {
	C.td_json_client_destroy(c.Client)
}

// Sends request to the TDLib client.
// You can provide string or Update.
func (c *Client) Send(jsonQuery interface{}) {
	var query *C.char

	switch jsonQuery.(type) {
	case string:
		query = C.CString(jsonQuery.(string))
	case Update:
		jsonBytes, _ := json.Marshal(jsonQuery.(Update))
		query = C.CString(string(jsonBytes))
	}

	defer C.free(unsafe.Pointer(query))
	C.td_json_client_send(c.Client, query)
}

// Receives incoming updates and request responses from the TDLib client.
// You can provide string or Update.
func (c *Client) Receive(timeout float64) Update {
	result := C.td_json_client_receive(c.Client, C.double(timeout))

	var update Update
	json.Unmarshal([]byte(C.GoString(result)), &update)
	return update
}

// Synchronously executes TDLib request.
// Only a few requests can be executed synchronously.
func (c *Client) Execute(jsonQuery interface{}) Update {
	var query *C.char

	switch jsonQuery.(type) {
	case string:
		query = C.CString(jsonQuery.(string))
	case Update:
		jsonBytes, _ := json.Marshal(jsonQuery.(Update))
		query = C.CString(string(jsonBytes))
	}

	defer C.free(unsafe.Pointer(query))
	result := C.td_json_client_execute(c.Client, query)

	var update Update
	json.Unmarshal([]byte(C.GoString(result)), &update)
	return update
}

// Sets the path to the file to where the internal TDLib log will be written.
// By default TDLib writes logs to stderr or an OS specific log.
// Use this method to write the log to a file instead.
func SetFilePath(path string) {
	query := C.CString(path)
	defer C.free(unsafe.Pointer(query))

	C.td_set_log_file_path(query)
}

// Sets the verbosity level of the internal logging of TDLib.
// By default the TDLib uses a verbosity level of 5 for logging.
func SetLogVerbosityLevel(level int) {
	C.td_set_log_verbosity_level(C.int(level))
}

// Sends request to the TDLib client and catches the result in updates channel.
// You can provide string or Update.
func (c *Client) SendAndCatch(jsonQuery interface{}) (Update, error) {
	var update Update

	switch jsonQuery.(type) {
	case string:
		// unmarshal JSON into map, we don't have @extra field, if user don't set it
		json.Unmarshal([]byte(jsonQuery.(string)), &update)
	case Update:
		update = jsonQuery.(Update)
	}

	// letters for generating random string
	letterBytes := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	// generate random string for @extra field
	b := make([]byte, 32)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	randomString := string(b)

	// set @extra field
	update["@extra"] = randomString

	// create waiter chan and save it in Waiters
	waiter := make(chan Update, 1)
	c.waiters.Store(randomString, waiter)

	// send it through already implemented method
	c.Send(update)

	select {
	// wait response from main loop in NewClient()
	case response := <-waiter:
		return response, nil
		// or timeout
	case <-time.After(10 * time.Second):
		c.waiters.Delete(randomString)
		return Update{}, errors.New("timeout")
	}
}

// Method for interactive authorizations process, just provide it authorization state from updates and api credentials.
func (c *Client) Auth(authorizationState string) (Update, error) {
	switch authorizationState {
	case "authorizationStateWaitTdlibParameters":
		res, err := c.SendAndCatch(Update{
			"@type":      "setTdlibParameters",
			"parameters": c.parameters.toTdlibParameters(),
		})
		if err != nil {
			return nil, err
		}
		return res, nil
	case "authorizationStateWaitEncryptionKey":
		res, err := c.SendAndCatch(Update{
			"@type": "checkDatabaseEncryptionKey",
		})
		if err != nil {
			return nil, err
		}
		return res, nil
	case "authorizationStateWaitPhoneNumber":
		fmt.Print("Enter phone: ")
		var number string
		fmt.Scanln(&number)

		res, err := c.SendAndCatch(Update{
			"@type":        "setAuthenticationPhoneNumber",
			"phone_number": number,
		})
		if err != nil {
			return nil, err
		}
		return res, nil
	case "authorizationStateWaitCode":
		fmt.Print("Enter code: ")
		var code string
		fmt.Scanln(&code)

		res, err := c.SendAndCatch(Update{
			"@type": "checkAuthenticationCode",
			"code":  code,
		})
		if err != nil {
			return nil, err
		}
		return res, nil
	case "authorizationStateWaitPassword":
		fmt.Print("Enter password: ")
		var passwd string
		fmt.Scanln(&passwd)

		res, err := c.SendAndCatch(Update{
			"@type":    "checkAuthenticationPassword",
			"password": passwd,
		})
		if err != nil {
			return nil, err
		}
		return res, nil
	case "authorizationStateReady":
		fmt.Println("Authorized!")
		return nil, nil
	default:
		return nil, errors.New(fmt.Sprintf("unexpected authorization state: %s", authorizationState))
	}
}
