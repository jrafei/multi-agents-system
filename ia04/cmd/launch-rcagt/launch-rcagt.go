package main

import (
	"fmt"

	"ia04/agt/restvoteragent"
)

func main() {
	ag := restclientagent.NewRestClientAgent("id1", "http://localhost:8080", "+", 11, 1)
	ag.Start()
	fmt.Scanln()
}
