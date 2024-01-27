package main

import (
	"fmt"
	"log"
	"net/http"
)

const webPort = 80

func main() {

	fmt.Printf("Starting web server on http://localhost:%d\n", webPort)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", webPort),
		Handler: routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Fatalf("couldn't start the web Server %v\n", err)
	}

}
