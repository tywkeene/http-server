package main

import (
	"github.com/SaviorPhoenix/http-server/cache"
	"github.com/SaviorPhoenix/http-server/handles"
	"log"
	"net/http"
)

func init() {
	//Register our path handles
	handles.Register()

	//We can't do much if we don't have a document cache, so panic
	// if we fail to get one
	if err := cache.InitCache("./docs/"); err != nil {
		panic(err)
	}
}

func main() {
	log.Println("Listening on port 3000")
	if err := http.ListenAndServe(":3000", nil); err != nil {
		panic(err)
	}
}
