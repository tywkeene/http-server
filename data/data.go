package data

import (
	"fmt"
	"log"
)

type Getter struct {
	Name string
	Get  GetterFunc
}

type GetterFunc func(interface{}) interface{}

//map[document name] data-getter defined in getters/getters.go
var Getters map[string]*Getter

//Initializes the global/exported Getters map
func InitGetters() {
	log.Println("Initializing data-getter function map")
	Getters = make(map[string]*Getter)
}

//Allocates and returns a new Getter struct
func NewGetter(name string, get GetterFunc) *Getter {
	return &Getter{name, get}
}

//Returns true if a getter exists in Getters, false otherwise
func getterExists(name string) bool {
	_, ok := Getters[name]
	return ok
}

//Register a new getter in Getters. Returns error if a getter with 'name' already exists
func RegisterGetter(name string, get GetterFunc) error {
	if getterExists(name) == true {
		return fmt.Errorf("Data-getter already exists for document: %s", name)
	} else {
		log.Println("\t++ Registering data getter ", get, " for", name)
		Getters[name] = NewGetter(name, get)
	}
	return nil
}

//Get a data-getter from Getters if it exists, returns error otherwise
func GetGetter(name string) (*Getter, error) {
	if getterExists(name) == false {
		return nil, fmt.Errorf("No data-getter exists for document: %s", name)
	} else {
		return Getters[name], nil
	}
}
