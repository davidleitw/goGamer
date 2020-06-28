package main

import "goGamer/server"

func main() {
	r := server.NewRouter()
	r.Run()
}
