package refresh

import (
	"github.com/SaviorPhoenix/http-server/cache"
	"log"
	"os"
	"testing"
)

func TestInitWatch(t *testing.T) {
	const watchDir = "../docs"

	//Initialize the cache so we have something to watch
	if err := cache.InitCache(watchDir); err != nil {
		t.Fatalf("Failed to initialize cache %s: %s\n", watchDir, err)
	}

	if err := InitCacheWatch(cache.Docs); err != nil {
		t.Fatalf("Failed to initialize watch cache %s: %s\n", cache.Docs.Path, err)
	}
}

//Simply create a dummy document
func CreateDocument(docDir string, docName string) {
	file, _ := os.Create(docDir + "/" + docName)
	defer file.Close()
	file.WriteString(string("<html>test</html>"))
	log.Println("Created dummy document", docDir+"/"+docName)
}

func TestWatchCache(t *testing.T) {
	const watchDir = "../docs"

	//Initialize the cache so we have something to watch
	if err := cache.InitCache(watchDir); err != nil {
		t.Fatalf("Failed to initialize cache %s: %s\n", watchDir, err)
	}

	if err := InitCacheWatch(cache.Docs); err != nil {
		t.Fatalf("Failed to initialize watch cache %s: %s\n", cache.Docs.Path, err)
	}

	if err := Watch.WatchCache(); err != nil {
		t.Fatalf("Failed to watch cache %s: %s\n", watchDir, err)
	}

	CreateDocument(watchDir, "testing.html")
	os.Remove(watchDir + "/" + "testing.html")

	CreateDocument(watchDir, "testing.txt")
	os.Remove(watchDir + "/" + "testing.txt")
}
