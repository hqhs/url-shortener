package main

import (
	"flag"
	"fmt"

	"github.com/hqhs/url-shortener/service"
)

var routes = flag.Bool("routes", false, "Generate router documentation")

func main() {
	flag.Parse()
	// TODO: bind env variables, such as host, port etc, ssl.
	service := service.NewService("localhost:3333", service.NewMockDatabase)
	// Passing -routes to the program will generate docs for the above
	// router definition. See the `routes.json` file in this folder for
	// the output.
	if *routes {
		fmt.Println(service.RoutesDoc())
		return
	}
	fmt.Println("Bootstrapped service...")
	service.Serve()
}
