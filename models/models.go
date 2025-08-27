package models

import (
	"errors"
	"sync"
)

type Customer struct {
	Name         string   `json:"name"`
	Email        string   `json:"email"`
	Phone        string   `json:"phone"`
	ID           string   `json:"id"`
	Restricted   bool     `json:"restricted"`
	Restrictions []string `json:"restrictions"`
}

type GeoIPResponse struct {
	IPAddress string `json:"ip_address"`
	City      struct {
		Names map[string]string `json:"names"`
	} `json:"city"`
	Country struct {
		IsoCode string            `json:"iso_code"`
		Names   map[string]string `json:"names"`
	} `json:"country"`
	Location struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"location"`
}

// Example IPs for testing geo-restrictions
var countryIPs = map[string]string{
	"AU": "101.167.184.6", // not in any of the test customers' allowed countries
	"GB": "104.239.8.7",   // allowed for Dunder Mifflin and Daily Planet
	"SG": "103.140.3.5",   // allowed for Dunder Mifflin
	"US": "102.129.255.4", // allowed for ACME, Dunder Mifflin, and Ghostbusters
}

var customerData = []Customer{
	{
		Name:         "ACME Corp",
		ID:           "089a1e8e-3ae4-483a-876e-27e64ed44388",
		Email:        "info@acme.com",
		Phone:        "555-123-4567",
		Restricted:   true,
		Restrictions: []string{"US"},
	},
	{
		Name:         "Dunder Mifflin Paper Company",
		ID:           "27d8e79f-c3c2-4a3c-a07f-8a6a019725da",
		Email:        "staff@dundermifflin.com",
		Phone:        "555-867-5309",
		Restricted:   true,
		Restrictions: []string{"US", "GB", "SG"},
	},
	{
		Name:         "Jurassic Fork",
		ID:           "030621c5-f493-4d1c-b478-2c64bfcd5dd9",
		Email:        "help@jurassicfork.com",
		Phone:        "",
		Restricted:   false,
		Restrictions: []string{},
	},
	{
		Name:         "Ghostbusters Inc",
		ID:           "bbddea84-e555-4a50-8d97-9a42f69b4184",
		Email:        "iaintfraidofnoghosts@ghostbusters.com",
		Phone:        "555-555-5555",
		Restricted:   true,
		Restrictions: []string{"US"},
	},
	{
		Name:         "Daily Planet",
		ID:           "6a61df14-ca32-415f-86f3-3a4f5e739bbf",
		Email:        "subscriptions@dailyplanet.com",
		Phone:        "555-234-5678",
		Restricted:   true,
		Restrictions: []string{"GB"},
	},
}

var mu sync.Mutex

// GetCustomers returns all customers
func GetCustomers() []Customer {
	mu.Lock()
	defer mu.Unlock()
	return customerData
}

// GetCustomerByID fetches a customer by their ID
func GetCustomerByID(id string) (*Customer, error) {
	mu.Lock()
	defer mu.Unlock()

	for _, customer := range customerData {
		if customer.ID == id {
			return &customer, nil
		}
	}
	return nil, errors.New("Customer not found")
}

// GetCustomerByName fetches a customer by their Name
func GetCustomerByName(name string) (*Customer, error) {
	mu.Lock()
	defer mu.Unlock()

	for _, customer := range customerData {
		if customer.Name == name {
			return &customer, nil
		}
	}
	return nil, errors.New("Customer not found")
}

// AddCustomer adds a new customer to the data
func AddCustomer(customer Customer) {
	mu.Lock()
	defer mu.Unlock()
	customerData = append(customerData, customer)
}

// UpdateCustomer updates an existing customer by ID or Customer name
func UpdateCustomer(identifier string, updated Customer) error {
	mu.Lock()
	defer mu.Unlock()

	for i, customer := range customerData {
		if customer.ID == identifier || customer.Name == identifier {
			customerData[i] = updated
			return nil
		}
	}
	return errors.New("Customer not found")
}

// DeleteCustomer removes a customer by ID or Customer name
func DeleteCustomer(identifier string) error {
	mu.Lock()
	defer mu.Unlock()

	for i, customer := range customerData {
		if customer.ID == identifier || customer.Name == identifier {
			customerData = append(customerData[:i], customerData[i+1:]...)
			return nil
		}
	}
	return errors.New("Customer not found")
}
