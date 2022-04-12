package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", index)
	http.ListenAndServe(":8080", nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	html := `<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<title>OAuth2 Example</title>
	</head>
	<body>
		<form action="/oauth" method="post">
			<input type="email" name="email" />
			<input type="submit" value="Github Login"/>
		</form>
	</body>
	</html>`
	fmt.Fprint(w, html)
}
