package main

import (
	"flag"
	"fmt"

	"example.com/url-shortener/service"
)

var routes = flag.Bool("routes", false, "Generate router documentation")

func main() {
	service := service.NewService()
	// Passing -routes to the program will generate docs for the above
	// router definition. See the `routes.json` file in this folder for
	// the output.
	if *routes {
		fmt.Println(service.RoutesDoc())
	}
	service.Serve()
}
