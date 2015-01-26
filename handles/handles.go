package handles

import (
	"github.com/SaviorPhoenix/http-server/cache"
	"html/template"
	"log"
	"net/http"
)

type PageData struct {
	Title string
}

func Register() {
	//Catch all handler
	http.HandleFunc("/", RootHandle)

	//For static images/stylesheets/files
	http.HandleFunc("/static/", StaticHandle)
}

//Handles static file requests.
func StaticHandle(res http.ResponseWriter, req *http.Request) {
	log.Println("<< GET /static -", req.UserAgent())
	http.ServeFile(res, req, req.URL.Path[1:])
}

//Catch all handler.
// If no document name is in the url, (i.e localhost/) we return index.html
// If there is a document name then we look it up and return it
func RootHandle(res http.ResponseWriter, req *http.Request) {
	var reply *template.Template

	log.Println("<< GET / -", req.UserAgent())
	if req.URL.Path[1:] == "" {
		reply = cache.Docs.GetDoc("index.html")
	} else {
		reply = cache.Docs.GetDoc(req.URL.Path[1:])
	}

	data := &PageData{"http-server"}
	reply.Execute(res, data)
}
