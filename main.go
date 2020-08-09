package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
	"net/http"
	//"go-restapi-practice/"
)


func main() {

	r := mux.NewRouter()
	db, err := gorm.Open("postgres", "host=http://localhost port=5432 dbname=postgres password=123qwe123qwe@")
	if err != nil {
		log.Fatal(err)
		fmt.Println("error")
		return
	}
	defer db.Close()
	http.Handle("/", r)
	fmt.Println("hello world")
	http.ListenAndServe(":8080", r)
}

///Users/admin/Desktop/TruongVN/Golang/go-restapi-practice