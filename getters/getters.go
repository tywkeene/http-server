package getters

import (
	"log"
)

type PageData struct {
	Title string
}

type LinkData struct {
	Title string
	Body  string
}

//index.html
func RootGetter(input interface{}) interface{} {
	data := &PageData{"http-server"}
	log.Println("\t** RootGetter() data:", data)
	return data
}

//link1.html
func LinkGetterOne(input interface{}) interface{} {
	data := &LinkData{"Link 1", "This is one body"}
	log.Println("\t** LinkGetterOne() data:", data)
	return data
}

//link2.html
func LinkGetterTwo(input interface{}) interface{} {
	data := &LinkData{"Link 2", "This is another body"}
	log.Println("\t** LinkGetterTwo() data:", data)
	return data
}
