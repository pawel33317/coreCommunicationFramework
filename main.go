package main

import (
	"fmt"

	"github.com/pawel33317/coreCommunicationFramework/dbHandler"
)

func main() {
	fmt.Println("Hello, World!")
	dbHandler.RunHandler("test")
}
