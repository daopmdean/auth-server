package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type person struct {
	First string
}

func main() {
	p1 := person{
		First: "Dean",
	}
	p2 := person{
		First: "Soul",
	}
	xp := []person{p1, p2}
	byts, err := json.Marshal(xp)
	if err != nil {
		log.Panic(err)
	}
	fmt.Println("MARSHAL", string(byts))

	xp2 := []person{}
	err = json.Unmarshal(byts, &xp2)
	if err != nil {
		log.Panic(err)
	}
	fmt.Println("UNMARSHAL", xp2)

	http.HandleFunc("/encode", encodeFunc)
	http.HandleFunc("/decode", decodeFunc)
	http.ListenAndServe(":8080", nil)
}

func encodeFunc(w http.ResponseWriter, r *http.Request) {
	p1 := person{
		First: "Dean",
	}
	err := json.NewEncoder(w).Encode(p1)
	if err != nil {
		log.Println("Encoded bad data", err)
	}
}

func decodeFunc(w http.ResponseWriter, r *http.Request) {

}
