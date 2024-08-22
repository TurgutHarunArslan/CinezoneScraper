package api

import (
	"log"
	"net/http"
)

func StartFrontEnd() {
	fs := http.FileServer(http.Dir("./frontend"))

	http.Handle("/", fs)

	log.Println("Serving on http://localhost:8000")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal(err)
	}
}