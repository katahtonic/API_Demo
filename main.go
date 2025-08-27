package main

import (
	"log"
	"net/http"
	"os"

	"API_Demo/handlers"
)

func main() {
	/// HTTP Handlers / Routes
	// Customer endpoints
	http.HandleFunc("/listCustomers", handlers.ListCustomers)
	http.HandleFunc("/addCustomer", handlers.AddCustomer)
	http.HandleFunc("/editCustomer", handlers.EditCustomer)
	http.HandleFunc("/getCustomer", handlers.GetCustomer)
	http.HandleFunc("/deleteCustomer", handlers.DeleteCustomer)

	// User endpoints
	http.HandleFunc("/login", handlers.Login)
	http.HandleFunc("/loginValidation", handlers.LoginValidation)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on port %s", port)
	if err := http.ListenAndServe("localhost:"+port, nil); err != nil {
		log.Fatalf("Could not start server: %s\n", err.Error())
	}
}
