# Shex API Client

Client library for the Shex API. 

Supported Endpoints:

- Mods
- Profiles


## Installation

> $ go get github.com/danielkrainas/shexd/api/client


## Usage

How to instantiate a new client:

```go
package main

import (
	"net/http"
	"log"
	
	"github.com/danielkrainas/shexd/api/client"
)

// http/https url of the shexd service
const ENDPOINT = "http://localhost:9366"

func main() {
	// Create a new client
	c := client.New(ENDPOINT, http.DefaultClient)

	if err := c.Ping(); err != nil {
		log.Fatal(err)
		return
	}
}
```
	
## Example

A more detailed example can be found [here.](https://github.com/danielkrainas/shexd/tree/master/api/client/example)

