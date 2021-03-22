package main

import (
	"fmt"
	httphandlers "identity/handlers/http"
	timerhandlers "identity/handlers/timer"
	"log"
	"net/http"
	"os"
)

func init() {
	fmt.Println("package: main - initialized")
}

func main() {

	listenAddr := ":8080"

	if val, ok := os.LookupEnv("FUNCTIONS_CUSTOMHANDLER_PORT"); ok {
		listenAddr = ":" + val
	}

	// httpTriggers
	http.HandleFunc("/api/authuser", httphandlers.AuthUser)
	http.HandleFunc("/api/createuser", httphandlers.CreateUser)
	http.HandleFunc("/api/getuser", httphandlers.GetUser)

	// timerTriggers
	http.HandleFunc("/authtoken", timerhandlers.GetAuthToken)
	http.HandleFunc("/publickeys", timerhandlers.GetPublicKeys)

	log.Printf("About to listen on %s. Go to https://127.0.0.1%s/", listenAddr, listenAddr)
	log.Fatal(http.ListenAndServe(listenAddr, nil))
}
