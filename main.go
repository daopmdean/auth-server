package main

import (
	"encoding/json"
	"fmt"
	"log"
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
}
