package main

import (
	"fmt"

	"ia04/agt/restserveragent"
)

func main() {
	server := restserveragent.NewRestServerAgent(":8080")
	server.Start()
	fmt.Scanln()
}
