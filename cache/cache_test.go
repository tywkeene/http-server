package cache

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"
)

// Get the count of documents and total size of the documents in docDir
func getCacheInfo(docDir string) (int64, int) {
	var size int64
	var count int

	list, _ := ioutil.ReadDir(docDir)
	for _, file := range list {
		if strings.HasSuffix(file.Name(), ".html") == true {
			stat, _ := os.Stat(docDir + "/" + file.Name())
			size += stat.Size()
			count++
		}
	}
	return size, count
}

func TestInitCache(t *testing.T) {
	const docDir = "../docs"
	const badDocDir = "../ducks"
	const badDoc = "bad.txt"

	if err := InitCache(badDocDir); err == nil {
		t.Fatal("Recieved no error when attempting to build cache from bogus directory")
	}

	//Create a bogus document in the cache directory to test that we ignore bad documents
	CreateDocument(docDir, badDoc)

	if err := InitCache(docDir); err != nil {
		t.Fatal("Failed to initialize document cache in", docDir)
	}
	os.Remove(docDir + "/" + badDoc)

	//Ensure we recieved an initialized Doc structure
	if Docs == nil {
		t.Fatal("Document structure should not be nil")
	}
	//Ensure the cache was actually filled by Docs.BuildCache()
	if Docs.docs == nil {
		t.Fatal("Document cache should not be nil")
	}

	//Ensure we have all of the documents are in the cache and they aren't nil
	var expectDocs = []string{"404.html", "index.html", "link1.html", "link2.html"}

	for _, name := range expectDocs {
		if Docs.Exists(name) != true {
			t.Fatal("Non-existent document on disk or cache", name)
		}
	}

	expectedSize, expectedCount := getCacheInfo(docDir)

	//Ensure we get the corrent size and count of the documents
	if Docs.size != expectedSize {
		t.Fatalf("Inconsistent size of documents in cache %s. Got %d should be %d",
			docDir, Docs.size, expectedSize)
	}

	if Docs.count != expectedCount {
		t.Fatalf("Inconsistent count of documents in cache %s. Got %d should be %d",
			docDir, Docs.count, expectedCount)
	}
}

func TestCacheDoc(t *testing.T) {
	const badDoc = "bad.html"
	if Docs.CacheDoc(badDoc) == nil {
		t.Fatal("Did not recieve error when attempting to cache non-existent document")
	}
}

//Simply create a dummy document
func CreateDocument(docDir string, docName string) {
	file, _ := os.Create(docDir + "/" + docName)
	defer file.Close()
	file.WriteString(string("<html>test</html>"))
	log.Println("Created dummy document", docDir+"/"+docName)
}

func TestGetDoc(t *testing.T) {
	const docDir = "../docs"
	const docName = "index.html"
	const testName = "test.html"

	// Test getting document from the cache
	if Docs.GetDoc(docName) == nil {
		t.Fatalf("Failed to get existing document %s in cache built from %s", docName, docDir)
	}

	// Test reading a new document off disk
	CreateDocument(docDir, testName)
	if Docs.GetDoc("test.html") == nil {
		os.Remove(docDir + "/" + testName)
		t.Fatalf("Failed to get document test.html from cache %s: %s", docDir)
	}
	os.Remove(docDir + "/" + testName)

	// Test trying to get a non-existant document
	if Docs.GetDoc("not_there.html") != nil {
		t.Fatalf("GetDoc() returned non-nil when looking for non-existant document")
	}
}

func TestRefreshDoc(t *testing.T) {
	const docName = "index.html"
	const badDoc = "notthere.html"

	//Ensure we can refresh an existing document
	if err := Docs.RefreshDoc(docName); err != nil {
		t.Fatalf("Failed to refresh document %s in %s", docName, Docs.Path)
	}

	//Ensure we get an error if we try to refresh a phantom document
	if err := Docs.RefreshDoc(badDoc); err == nil {
		t.Fatal("Did not recieve error when attempting to refresh non-existent document")
	}
}
