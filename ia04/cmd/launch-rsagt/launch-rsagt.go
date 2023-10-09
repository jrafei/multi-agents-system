package main

import (
	"fmt"

	"ia04/agt/restballotagent"
)

func main() {
	server := restserveragent.NewRestServerAgent(":8080")
	server.Start()
	fmt.Scanln()
}
