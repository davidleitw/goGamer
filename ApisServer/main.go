package main

import (
	"fmt"
	"goGamer/server"
)

func main() {
	r := server.NewRouter()
	err := r.Run(":3000")
	if err != nil {
		fmt.Println("Some error with server.")
	}
}
