package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/go-chi/docgen"
	"github.com/gomodule/redigo/redis"
)

var routes = flag.Bool("routes", false, "Generate router documentation")

func main() {
	flag.Parse()

	p := NewPool("redis:6379")

	// Simple health check: make ping request, get or create url id
	conn := p.Get()
	defer conn.Close()
	_, err := redis.String(conn.Do("PING", ""))
	if err != nil {
		panic(err)
	}
	id, err := GetOrCreateID(conn)
	fmt.Printf("url id starts with: %v\n", id)

	r := NewRouter(p)
	fmt.Println("Initialized router...")

	// Passing -routes to the program will generate docs for the above
	// router definition. See the `routes.json` file in this folder for
	// the output.
	if *routes {
		// fmt.Println(docgen.JSONRoutesDoc(r))
		fmt.Println(docgen.MarkdownRoutesDoc(r, docgen.MarkdownOpts{
			ProjectPath: "github.com/go-chi/chi",
			Intro:       "Welcome to the chi/_examples/rest generated docs.",
		}))
		return
	}

	http.ListenAndServe(":3333", r)
}



