package main

import (
	"flag"
	"github.com/SaviorPhoenix/http-server/cache"
	"github.com/SaviorPhoenix/http-server/data"
	"github.com/SaviorPhoenix/http-server/handles"
	"log"
	"net/http"
)

type Options struct {
	port     *string
	useTls   *bool
	certPath *string
	certKey  *string
}

var Opts *Options

func init() {
	Opts = &Options{
		flag.String("port", "", "Bind port (Default 80 for HTTP, 443 for HTTPS)"),
		flag.Bool("use-tls", false, "Use SSL/TLS (https)"),
		flag.String("cert", "./cert.pem", "Path to a signed SSL/TLS certificate"),
		flag.String("cert-key", "./key.pem", "Path to the key associated with the signed SSL/TLS certificate"),
	}
	flag.Parse()

	//Set the default port depending on if we're serving HTTP or HTTPS
	if *Opts.useTls == true && *Opts.port == "" {
		*Opts.port = "443"
	} else if *Opts.useTls == false && *Opts.port == "" {
		*Opts.port = "80"
	}

	//Initialize the data-getter map for the documents
	data.InitGetters()

	//Register our path handles
	handles.Register()

	//We can't do much if we don't have a document cache, so panic
	// if we fail to get one
	if err := cache.InitCache("./docs/"); err != nil {
		panic(err)
	}

}

func main() {
	if *Opts.useTls == true {
		log.Println("Accepting HTTPS on port:", *Opts.port)
		if err := http.ListenAndServeTLS(":"+*Opts.port, *Opts.certPath, *Opts.certKey, nil); err != nil {
			panic(err)
		}
	} else {
		log.Println("Accepting HTTP on port:", *Opts.port)
		if err := http.ListenAndServe(":"+*Opts.port, nil); err != nil {
			panic(err)
		}
	}
}
