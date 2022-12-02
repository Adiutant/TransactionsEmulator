package main

import "TestBank/http_server"

func main() {
	server, err := http_server.NewBankServer()
	if err != nil {
		panic(err)
	}
	server.SetRoutes()
	err = server.Run()
	if err != nil {
		panic(err)
	}

}
