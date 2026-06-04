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
		fmt.Fprintf(w, "Welcome to Docker Non-Root User Example!\n")
		fmt.Fprintf(w, "Your UID is: %d\n", uid)

		if uid == 0 {
			fmt.Fprint(w, "⚠️WARNING: running as ROOT!\n")
		} else {
			fmt.Fprint(w, "✅ SUCCESS:Running as NON-ROOT user.\n")
		}
	})
	
	fmt.Println("Server starting on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}