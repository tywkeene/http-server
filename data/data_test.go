package data

import (
	"testing"
)

type dummyStruct struct {
	num int
	str string
}

func DummyGetter(foo interface{}) interface{} {
	return &dummyStruct{1337, "bar"}
}

func TestRegisterGetter(t *testing.T) {
	const docName = "index.html"

	InitGetters()
	if Getters == nil {
		t.Fatal("Failed to initialize data getter function map")
	}

	//Ensure we can register a getter for a document
	if err := RegisterGetter(docName, DummyGetter); err != nil {
		t.Fatalf("Failed to register data getter %p for page %s", DummyGetter, docName)
	}

	//Ensure we recieve an error if we try to register a getter again
	if err := RegisterGetter(docName, DummyGetter); err == nil {
		t.Fatalf("Should have recieved error when attempting to register existent function",
			DummyGetter, docName)
	}
}

func TestGetGetter(t *testing.T) {
	const docName = "index.html"

	//Getters in data.go must be initialized
	InitGetters()
	if Getters == nil {
		t.Fatal("Failed to initialize data getter function map")
	}

	if err := RegisterGetter(docName, DummyGetter); err != nil {
		t.Fatalf("Failed to register data getter %p for page %s", DummyGetter, docName)
	}

	//Ensure we get a function for an existent document
	getter, err := GetGetter(docName)
	if getter == nil || err != nil {
		t.Fatal("Failed to get existent data-getter for document", docName)
	}

	//Ensure we get an error if we try to get a function for a non-existent document
	getter, err = GetGetter("non-existent.html")
	if getter != nil || err == nil {
		t.Fatalf("Recieved getter %p for non-existent document %s", getter, docName)
	}
}
