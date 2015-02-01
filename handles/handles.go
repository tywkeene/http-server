package handles

import (
	"github.com/SaviorPhoenix/http-server/cache"
	"github.com/SaviorPhoenix/http-server/data"
	"github.com/SaviorPhoenix/http-server/getters"
	"html/template"
	"log"
	"net/http"
)

func Register() {

	//Register data-getters for the pages from getters/getters.go
	data.RegisterGetter("index.html", getters.RootGetter)
	data.RegisterGetter("link1.html", getters.LinkGetterOne)
	data.RegisterGetter("link2.html", getters.LinkGetterTwo)
	data.RegisterGetter("404.html", getters.FourOhFour)

	//Catch all handler
	http.HandleFunc("/", RootHandle)

	//For static images/stylesheets/files
	http.HandleFunc("/static/", StaticHandle)
}

//Handles static file requests.
func StaticHandle(res http.ResponseWriter, req *http.Request) {
	log.Printf("<< GET /%s - %s", req.URL.Path[1:], req.UserAgent())
	http.ServeFile(res, req, req.URL.Path[1:])
}

//Catch all handler.
// If no document name is in the url, (i.e localhost/) we return index.html
// If there is a document name then we look it up and return it
func RootHandle(res http.ResponseWriter, req *http.Request) {
	var reply *template.Template
	var docName string

	log.Printf("<< GET /%s - %s\n", req.URL.Path[1:], req.UserAgent())

	if req.URL.Path[1:] == "" {
		docName = "index.html"
	} else {
		docName = req.URL.Path[1:]
	}
	reply = cache.Docs.GetDoc(docName)

	getter, err := data.GetGetter(docName)
	if err != nil {
		log.Println(err)
		return
	}
	data := getter.Get(req.UserAgent())
	reply.Execute(res, data)
}
