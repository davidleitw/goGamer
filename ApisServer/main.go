package main

import (
	"fmt"
	"goGamer/server"
)

func main() {
	r := server.NewRouter()
	err := r.Run()
	if err != nil {
		fmt.Println("Some error with server.")
	}
}
