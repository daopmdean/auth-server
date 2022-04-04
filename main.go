package main

import (
	"encoding/base64"
	"log"
)

func main() {
	log.Println(base64.StdEncoding.EncodeToString([]byte("user:pass")))
}
