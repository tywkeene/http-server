# http-server
http-server is a simple http-server in golang. It tries to be more efficient and
disk friendly by reading documents into memory and serving them from there. It will
also search for a file on the disk if it doesn't exist in the cache, and caches it if
it has appeared. It doesn't return 404 unless the file really _really_ doesn't exist.

## Cache Implementation
It's rather simple, really.

First we need an appropriate data structure. I chose golang's map type since it's well
suited for 'looking things up' by name. We wrap this map in a simple struct like so:
```Go
type DocCache struct {
        docs  map[string]string //map of documents in 'path', indexed by filename
        size  int64             //Size of the document in bytes
        count int               //Amount of documents in the cache
        path  string            //Path to the document directory
}
```

First we have the ```docs``` variable which is the map we talked about. This is a simple
name->data map. It's indexed by filename for convienence and stores a string representation of
the document we can write to a ```http.request``` writer. 

When we want to fill this cache at startup, we step through the document directory
```path``` and enumerate each of the documents, storing them in this cache.

The ```size``` and ```count``` variables are used mainly for logging, they don't really affect much else.
However, these could be used to limit the size and count of files read into the cache, but that's beyond the scope
of what I was trying to accomplish in this project

## Automatic Document Refreshing
One problem arises when you store all your documents in the cache. Since we read the documents into
memory at runtime, we can't update them later without restarting the server since we don't read the files
every time they're served, which means we'd have to restart the server every time we make a change. This is not
ideal and almost makes the cache useless.

This is where inotify comes in. [Inotify] (http://en.wikipedia.org/wiki/Inotify) is a great feature of the Linux
kernel that allows us to watch files for events. We can use this to make our server a little more nimble by watching
the document directory for changes and refreshing the cache when they occur. Golang already has a [package](https://godoc.org/golang.org/x/exp/inotify)
that makes this very simple

First we need a small data structure to hold some needed information:
```Go
type CacheWatch struct {
	target *cache.DocCache // The document cache from cache/cache.go
	mask uint32	       // The mask we pass to inotify
	watcher *inotify.Watcher // The inotify watcher
}
```

We can then pass all of this to the inotify package, and get ourselves a watch:
```Go
if err := watchCache.watcher.AddWatch(watchCache.target.Path, watchCache.mask); err != nil {
	return err
}
```

Finally, we run a simple loop wrapped in a go routine and act on the events:
```Go
go func(watchCache *CacheWatch) {
	for {
		select {
		case event := <-watchCache.watcher.Event:
			_, name := filepath.Split(event.Name)
			if strings.HasSuffix(name, ".html") {
				log.Println("Ignoring document:", name)
				continue
			}
			if event.Mask&watchCache.mask != 0 && watchCache.target.IsCached(name) == true {
				log.Println("\t~~ Document modified:", name)
				if err := cache.Docs.RefreshDoc(name); err != nil {
					log.Println(err)
				}
			}
		case err := <-watchCache.watcher.Error:
			log.Println(err)
		}
	}
}(watchCache)
```

Now we can run our server, make changes to the files in the document cache and have our changes show up on refresh
This will make it much easier to quickly and efficiently prototype pages.

## Handling Requests
Another dead simple solution.
```Go
func RootHandle(res http.ResponseWriter, req *http.Request) {
	var reply *template.Template
	var docName string

	log.Println("<< GET / -", req.UserAgent())

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
```

This simple handler can be used as a 'catch-all' route. First we check if there's even a document name in the url, if there
isn't, we simply return the ```index.html`` document. If there does happen to be a document name, we get it from the cache (or disk)
and return it instead.

html/document allows us to actually pass data to the documents for rendering, so we use it instead of a raw html file.
This is great, but not complete. Each page may need different data, and trying to figure out what data goes with what document requires a little more work which should (and will) be in its own package.

The way I see it is we can create "data getters" and register them with each document, instead of hardcoding them in the actual route handlers. This will keep things organized and will allow us to get different kinds of data in different ways (e.g redis, mongodb, sql or the server itself) making the server a little more adaptable.
(Implemented, see below)

## Static Files
No website is complete without a good stylesheet, maybe some images and javascript for good measure.
This is super easy:
```Go
func StaticHandle(res http.ResponseWriter, req *http.Request) {
	log.Println("<< GET /static -", req.UserAgent())
	http.ServeFile(res, req, req.URL.Path[1:])
}
```
See? All we have to do is use ```http.ServeFile()``` to return the contents of the requested file or directory.
Short and sweet.

## Getting data with data-getters
Webpages are very data-driven. What's a webpage without data? Sure you technically *could* hardcode things into the 
raw html file, but this is 2015, not the 90's.I've come up with a solution. The first thing that occured to me was 
that different pages need different data, and that data might need to be retrieved in different ways. This is another great use for the all-powerful map data structure (no language should be without it).

We can solve this little problem with a simple struct and a global variable that holds the structs:

```Go
type Getter struct {
	Name string
	Get  GetterFunc
}

var Getters map[string]*Getter
```

To register a data-getter, we can simple make sure it doesn't already exist, then add it to the map:
```Go
func RegisterGetter(name string, get GetterFunc) error {
	if getterExists(name) == true {
		return fmt.Errorf("Data-getter already exists for document: %s", name)
	} else {
		log.Println("\t++ Registering data getter ", get, " for", name)
		Getters[name] = NewGetter(name, get)
	}
	return nil
}
```

and to grab one of these data-getters:
```Go
func GetGetter(name string) (*Getter, error) {
	if getterExists(name) == false {
		return nil, fmt.Errorf("No data-getter exists for document: %s", name)
	} else {
		return Getters[name], nil
	}
}
```

This makes it a lot easier to get data in different ways (e.g from SQL, mongodb or redis). These getters are defined in getters/getters.go and are registered in handles/handles.go when we register our route handlers.

# SSL/TLS

In the year 2015, encryption is very necessary, even standard. This is very easy in Golang's case since it provides
```http.ListenAndServeTLS()``` which is just like ```http.ListenAndServe()``` but takes two additional arguments:

from the Golang wiki:

```func (srv *Server) ListenAndServeTLS(certFile, keyFile string) error```
> Filenames containing a certificate and matching private key for the server must be provided. If the certificate is signed by a certificate authority, the certFile should be the concatenation of the server's certificate followed by the CA's certificate.

So it's as simple as:

```Go
if err := http.ListenAndServeTLS(":"+*Opts.port, *Opts.certPath, *Opts.certKey, nil); err != nil {
	panic(err)
}
```

Just get the required paths and arguments by using the ```flag``` package, pass them to this function
and we're up and serving over a secure connection.

# Conclusion
This was actually a spur of the moment idea, and something I've toyed with before but never really felt I had accomplished
anything with any other language or framework. In Golang this was very easy and simple, even with my limited time in go I was
able to slap this together in just a few hours. I'm very excited about this, I definitely look forward to seeing what I can come up
with in webapp land using Go.


# Contact
Pull request? Questions? Criticism? You can hit me up on twitter [@tywkeene](https://twitter.com/tywkeene) or over email <tyrell.wkeene@gmail.com>

All feedback is greatly appreciated :)
