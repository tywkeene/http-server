package main

import (
	"flag"
	"github.com/tywkeene/http-server/cache"
	"github.com/tywkeene/http-server/config"
	"github.com/tywkeene/http-server/data"
	"github.com/tywkeene/http-server/handles"
	"github.com/tywkeene/http-server/refresh"
	"log"
	"net/http"
)

var Conf *config.Conf

func init() {
	var confFile string
	var err error

	Conf = &config.Conf{}
	flag.StringVar(&confFile, "config", "", "Path to a .toml config file (overrides command line arguments)")
	flag.StringVar(&Conf.Options.BindPort, "port", "", "Bind port (Default 80 for HTTP, 443 for HTTPS)")
	flag.StringVar(&Conf.Options.Cert, "cert", "./cert.pem", "Path to a signed SSL/TLS certificate")
	flag.StringVar(&Conf.Options.CertKey, "cert-key", "./key.pem", "Path to the key associated with the signed SSL/TLS certificate")
	flag.StringVar(&Conf.Options.DocDir, "doc-dir", "./docs/", "Path to the directory containing the documents to serve")
	flag.BoolVar(&Conf.Options.Refresh, "refresh", true, "Automatically Refresh modified documents in ./docs")
	flag.BoolVar(&Conf.Options.UseTls, "use-tls", false, "Use SSL/TLS (https)")
	flag.Parse()

	//Settings from a configuration file override the command line arguments
	if confFile != "" {
		log.Println("Configuration file overriding command line arguments")
		Conf, err = config.ParseConfig(confFile)
		if err != nil {
			panic(err)
		}
	}
	//Set the default port depending on if we're serving HTTP or HTTPS
	if Conf.Options.UseTls == true && Conf.Options.BindPort == "" {
		Conf.Options.BindPort = "443"
	} else if Conf.Options.UseTls == false && Conf.Options.BindPort == "" {
		Conf.Options.BindPort = "80"
	}

	//Initialize the data-getter map for the documents
	data.InitGetters()

	//Register our path handles
	handles.Register()

	//We can't do much if we don't have a document cache, so panic
	// if we fail to get one
	if err := cache.InitCache(Conf.Options.DocDir); err != nil {
		panic(err)
	}

	if Conf.Options.Refresh == true {
		if err := refresh.InitCacheWatch(cache.Docs); err != nil {
			panic(err)
		}
		if err := refresh.Watch.WatchCache(); err != nil {
			panic(err)
		}
	}
}

func main() {
	if Conf.Options.UseTls == true {
		log.Println("Accepting HTTPS on port:", Conf.Options.BindPort)
		if err := http.ListenAndServeTLS(":"+Conf.Options.BindPort, Conf.Options.Cert, Conf.Options.CertKey, nil); err != nil {
			panic(err)
		}
	} else {
		log.Println("Accepting HTTP on port:", Conf.Options.BindPort)
		if err := http.ListenAndServe(":"+Conf.Options.BindPort, nil); err != nil {
			panic(err)
		}
	}
}
