package cache

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
)

type DocCache struct {
	docs  map[string]string //map of documents in 'path', indexed by filename
	size  int64             //Size of the document in bytes
	count int               //Amount of documents in the cache
	path  string            //Path to the document directory
}

var Docs *DocCache // Represents the global document cache

//Allocated and builds a new cache from the directory 'path'
func InitCache(path string) error {
	if Docs = NewDocCache(path); Docs == nil {
		return errors.New("Failed to get document cache")
	} else {
		return Docs.BuildCache()
	}
}

//Allocated a new DocCache with path 'path'
func NewDocCache(path string) *DocCache {
	return &DocCache{make(map[string]string), 0, 0, path}
}

//Builds a cache of documents from the files in the 'docDir' path
func (cache *DocCache) BuildCache() error {
	log.Println("Building document cache from", cache.path)

	list, err := ioutil.ReadDir(cache.path)
	if err != nil {
		return err
	}

	for _, file := range list {
		log.Println("\t++", file.Name())
		cache.CacheDoc(file.Name())
	}

	log.Println("\t!!", cache.count, "file(s) (", cache.size, " bytes) in", cache.path)
	return nil
}

//Adds the document 'name' to the cache
func (cache *DocCache) CacheDoc(name string) error {
	data, err := ioutil.ReadFile(cache.path + name)
	if err != nil {
		return err
	}
	cache.docs[name] = string(data)
	stat, _ := os.Stat(cache.path + name)
	cache.size += stat.Size()
	cache.count++
	return nil
}

//Returns true if 'name' is in cache, false otherwise
func (cache *DocCache) IsCached(name string) bool {
	_, ok := cache.docs[name]
	return ok
}

//Returns true if doc 'name' is on disk, false otherwise
func (cache *DocCache) IsOnDisk(name string) bool {
	_, err := os.Stat(name)
	return err != nil
}

//Looks for the document 'name' in the cache, then on the disk, then gives up and returns 404
// If the document is found in the cache it returns the document immediately
// If the document isn't in the cache, but on the disk, the document is read and
//	added to the cache and returned
// If the document isn't in the cache or on the disk, return 404
func (cache *DocCache) GetDoc(name string) string {
	const fourOhfourStr string = `<html><center><h1>404</h1><body>File not in
	cache or on disk<body></center></html>`

	log.Println("\t", name, "?? Querying cache")
	//File is in cache, simply return it
	if cache.IsCached(name) == true {
		log.Println("\t", name, "++ Cached/Found")
		return cache.docs[name]
		//Doc isn't in cache, see if it's on disk
		//cache and return it if it is
	} else if cache.IsOnDisk(name) == true {
		log.Println("\t", name, "!! Not in cache")
		data, err := ioutil.ReadFile(cache.path + name)
		if err != nil {
			log.Println("\t", name, "!! Not on disk")
		} else {
			log.Println("\t", name, "++ On disk/Found")
			log.Println("\t", name, "!! Caching")

			cache.CacheDoc(name)
			return string(data)
		}
	} else {
		//No luck in cache or on disk, return 404
		log.Println("\t", name, "-- Not cached/Not on disk")
		return fourOhfourStr
	}
	return fourOhfourStr
}