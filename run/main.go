package main

import (
	"RestGoTest/httpserver"
)

func main() {
	a := &httpserver.App{Port: ":9004"}
	a.Init()
	a.Run()

}
