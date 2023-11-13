**This library allows you to quickly and easily use the Tehnomir Web API via Go.**

## Install Package

`go get github.com/NuclearLouse/tehnomir`


## Quick Start

```go
package main

import (
	"fmt"
	"log"

	"github.com/NuclearLouse/tehnomir"
)

var TOKEN = "<TOKEN>"

func main() {
	cfg := tehnomir.DefaultConfig()
	cfg.Token = TOKEN
	tm := tehnomir.New(cfg)
	if err := tm.TestConnect("my custom test phrease"); err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Test: OK!")

	partnum := "<PART_NUM>"
	res, err := tm.SearchWithAnalogs(partnum)
	if err != nil {
		log.Fatalln(err)
	}
	for _, d := range res.Details {
		for _, s := range d.Stocks {
			fmt.Printf("Brand: %s | Code: %s | Supplier: %s | Price: %.2f\n",
				d.Brand,
				d.Code,
				s.PriceLogo,
				s.Price,
			)
		}
	}
}
```