package main

import (
	"encoding/json"
	"fmt"
)

type person struct {
	First string
}

func main() {
	p1 := person{First: "Dean"}
	p2 := person{First: "Soul"}
	xp := []person{p1, p2}
	byts, err := json.Marshal(xp)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(byts))
}
