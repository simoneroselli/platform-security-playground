package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// os.Getuid() returns the numeric user id of the caller
		uid := os.Getuid()
		fmt.Fprintf(w, "Welcome to the Security Playground!\n")
		fmt.Fprintf(w, "Current Process Running UID:: %d\n", uid)

		if uid == 0 {
			fmt.Fprint(w, "⚠️WARNING: Running as ROOT!\n")
		} else {
			fmt.Fprint(w, "✅ SUCCESS: Running as NON-ROOT user.\n")
		}
	})
	
	fmt.Println("Server starting on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}