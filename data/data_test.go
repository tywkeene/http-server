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

	if err := RegisterGetter(docName, DummyGetter); err != nil {
		t.Fatalf("Failed to register data getter %p for page %s", DummyGetter, docName)
	}
}
