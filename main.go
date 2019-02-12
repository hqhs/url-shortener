package main

import (
	"flag"
	"fmt"

	"github.com/hqhs/url-shortener/service"
	"github.com/hqhs/url-shortener/redis"
	kitlog "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
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
	options := service.Options{
		"localhost",
		"3333",
		"localhost:6379",
		NewRedisInstance,
	}
	// NOTE here I just pipe all logging from kitlog to stdlib log package,
	// which is seems useless. But default chi-router logging middleware
	// uses stdlib, and re-writing it seems like a lot of work and valid
	// pull-request %)
	// I'll do it later
	// TODO allow user defined default log level
	logger := kitlog.NewLogfmtLogger(kitlog.StdlibWriter{})
	logger = level.NewFilter(logger, level.AllowInfo())
	service, err := service.NewService(logger, options)
	if err != nil {
		panic(err)
	}
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
