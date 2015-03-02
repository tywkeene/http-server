package refresh

import (
	"github.com/SaviorPhoenix/http-server/cache"
	"golang.org/x/exp/inotify"
	"log"
	"path/filepath"
	"strings"
)

type CacheWatch struct {
	target  *cache.DocCache
	mask    uint32
	watcher *inotify.Watcher
}

var Watch *CacheWatch

func NewCacheWatch(target *cache.DocCache, mask uint32, watcher *inotify.Watcher) *CacheWatch {
	return &CacheWatch{target, mask, watcher}
}

func (watchCache *CacheWatch) watchFiles() error {
	if err := watchCache.watcher.AddWatch(watchCache.target.Path, watchCache.mask); err != nil {
		return err
	}
	return nil
}

func (watchCache *CacheWatch) WatchCache() error {
	log.Println("\t!! Starting watch on document cache:", filepath.Clean(watchCache.target.Path))

	err := watchCache.watchFiles()
	if err != nil {
		return err
	}
	go func(watchCache *CacheWatch) {
		for {
			select {
			case event := <-watchCache.watcher.Event:
				_, name := filepath.Split(event.Name)
				if strings.HasSuffix(name, ".html") != true {
					continue
				}
				if event.Mask&watchCache.mask != 0 {
					if err := cache.Docs.RefreshDoc(name); err != nil {
						log.Println(err)
					}
				}
			case err := <-watchCache.watcher.Error:
				log.Println(err)
			}
		}
	}(watchCache)
	return nil
}

func InitCacheWatch(target *cache.DocCache) error {
	const defaultMask = inotify.IN_CLOSE_WRITE

	log.Println("Initializing inotify watch on document cache:", target.Path)

	tmp, err := inotify.NewWatcher()
	if err != nil {
		return err
	}
	Watch = NewCacheWatch(target, defaultMask, tmp)
	return nil
}
