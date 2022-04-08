package handsonlv1

import (
	"encoding/json"
	"log"
	"net/http"
)

// type person struct {
// 	First string
// }

// func main() {
// 	http.HandleFunc("/decode", bar)
// 	http.ListenAndServe(":8082", nil)
// }

func bar(w http.ResponseWriter, r *http.Request) {
	slices := []person{}
	json.NewDecoder(r.Body).Decode(&slices)
	log.Println(slices)
}
