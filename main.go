package main

import (
	"fmt"
	"log"

	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

func index(c *fasthttp.RequestCtx) {
	fmt.Fprint(c, "Welcome!\n")
}

func hello(c *fasthttp.RequestCtx) {
	fmt.Fprintf(c, "Hello, %s\n", ctx.UserValue("name"))
}

func main() {
	router := fasthttprouter.New()
	router.GET("/", Index)
	router.GET("/hello/:name", Hello)
	router.GET("/static/*file", fasthttp.FSHandler("./static", 0))

	log.Fatal(fasthttp.ListenAndServe(":8080", router.Handler))
}
