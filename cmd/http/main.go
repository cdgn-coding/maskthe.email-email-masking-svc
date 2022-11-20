package main

import (
	"email-masking-svc/src/infrastructure/http"
)

func main() {
	app := http.NewServer()
	app.Run()
}
