package main

import (
	"flag"
	"fmt"

	"github.com/hqhs/url-shortener/service"
	"github.com/hqhs/url-shortener/redis"
)

var routes = flag.Bool("routes", false, "Generate router documentation")

// NewRedisInstance wraps redis.NewRedisConnectionPool for service usage
func NewRedisInstance(addr string) (service.Database, error) {
	ins, err := redis.NewRedisConnectionPool("localhost:6379")
	if err != nil {
		panic(err)
	}
	return ins, nil
}

func main() {
	flag.Parse()
	// TODO: bind env variables, such as host, port etc, ssl.
	// service := service.NewService("localhost:3333", service.NewMockDatabase)
	service := service.NewService("localhost:3333", "localhost:6379", NewRedisInstance)
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
