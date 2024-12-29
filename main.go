package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Pranavkharche24/mail/handlers"
)

func main() {
	http.HandleFunc("/", handlers.HandleHome)
	http.HandleFunc("/send", handlers.HandleSend)
	http.HandleFunc("/admin", handlers.HandleAdmin)
	http.HandleFunc("/admin/save", handlers.HandleAdminSave)

	fmt.Println("Server starting on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
