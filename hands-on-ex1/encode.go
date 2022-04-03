package main

import (
	"encoding/json"
	"net/http"
)

type person struct {
	Last string
}

func main() {
	http.HandleFunc("/encode", foo)
	http.ListenAndServe(":8081", nil)
}

func foo(w http.ResponseWriter, r *http.Request) {
	p1 := person{
		Last: "Person 1",
	}
	p2 := person{
		Last: "Person 2",
	}
	sls := []person{p1, p2}
	json.NewEncoder(w).Encode(sls)
}
