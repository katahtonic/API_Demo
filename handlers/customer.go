package handlers

import (
	"encoding/json"
	"net/http"

	models "API_Demo/models"
)

func EditCustomer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var customer models.Customer
	if err := json.NewDecoder(r.Body).Decode(&customer); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if customer.ID == "" {
		http.Error(w, "Customer ID is required", http.StatusBadRequest)
		return
	}

	if err := models.UpdateCustomer(customer.ID, customer); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Customer updated successfully"))
}

func ListCustomers(w http.ResponseWriter, r *http.Request) {
	customers := models.GetCustomers()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(customers)
}

func GetCustomer(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "ID query parameter is required", http.StatusBadRequest)
		return
	}

	customer, err := models.GetCustomerByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(customer)
}

func AddCustomer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var customer models.Customer
	if err := json.NewDecoder(r.Body).Decode(&customer); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if customer.ID == "" || customer.Name == "" {
		http.Error(w, "Customer ID and Name are required", http.StatusBadRequest)
		return
	}

	models.AddCustomer(customer)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Customer added successfully"))
}

func DeleteCustomer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "ID query parameter is required", http.StatusBadRequest)
		return
	}

	if err := models.DeleteCustomer(id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Customer deleted successfully"))
}
