package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	models "API_Demo/models"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file")
	}
}

const (
	GeoBaseURL     = "https://geoip.maxmind.com/geoip/v2.1/country"
	GeoLiteBaseURL = "https://geolite.info/geoip/v2.1/country"
)

// Login handles user login and verifies their location based on IP and customer info (using demo data from models)
// customer parameter expects a customer name or ID to verify against stored data
// example POST command: 127.0.0.1:8080/login
//
//	{
//	    "ip": "101.167.184.6",
//	    "customer": "27d8e79f-c3c2-4a3c-a07f-8a6a019725da"
//	}
func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var requestBody struct {
		IP       string `json:"ip"`
		Customer string `json:"customer"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	if requestBody.IP == "" {
		http.Error(w, "IP parameter is required", http.StatusBadRequest)
		return
	}

	if requestBody.Customer == "" {
		http.Error(w, "Customer parameter is required", http.StatusBadRequest)
		return
	}

	geoData, err := fetchGeoData(requestBody.IP)
	if err != nil {
		http.Error(w, "Failed to fetch geo data:", http.StatusInternalServerError)
		return
	}
	fmt.Println("GeoData:", geoData)

	c, err := models.GetCustomerByName(requestBody.Customer)
	if err != nil {
		c, err = models.GetCustomerByID(requestBody.Customer)
		if err != nil {
			http.Error(w, "Customer not found by name or ID", http.StatusNotFound)
			return
		}
	}
	fmt.Println("Customer Found:", c)

	if c.Restricted {
		if !contains(c.Restrictions, geoData.Country.IsoCode) {
			http.Error(w, "Access denied from this location: "+geoData.Country.IsoCode, http.StatusForbidden)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Verified User location is approved: " + geoData.Country.IsoCode))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("No restrictions. User location: " + geoData.Country.IsoCode))
}

// Login handles user login and verifies their location based off an IP and list of allowed countries (provided in the request body)
// example POST command: 127.0.0.1:8080/loginValidation
// example POST body:
//
//	{
//	    "ip": "101.167.184.6",
//	    "countries": ["US", "GB", "AU"]
//	}
func LoginValidation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var requestBody struct {
		IP        string   `json:"ip"`
		Countries []string `json:"countries"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	if requestBody.IP == "" {
		http.Error(w, "IP parameter is required", http.StatusBadRequest)
		return
	}

	if len(requestBody.Countries) == 0 {
		http.Error(w, "At least one country is required", http.StatusBadRequest)
		return
	}
	fmt.Println("Countries to validate:", requestBody.Countries)

	geoData, err := fetchGeoData(requestBody.IP)
	if err != nil {
		http.Error(w, "Failed to fetch geo data:", http.StatusInternalServerError)
		return
	}
	fmt.Println("GeoData:", geoData.Country.IsoCode)

	if !contains(requestBody.Countries, geoData.Country.IsoCode) {
		http.Error(w, "Access denied from this location: "+geoData.Country.IsoCode, http.StatusForbidden)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Verified User location is approved: " + geoData.Country.IsoCode))
	return
}

func contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

// fetchGeoData fetches geolocation data for a given IP address
func fetchGeoData(ip string) (models.GeoIPResponse, error) {
	GeoUsername := os.Getenv("GEO_USERNAME")
	GeoAPIKey := os.Getenv("GEO_API_KEY")

	url := fmt.Sprintf("%s/%s", GeoLiteBaseURL, ip)
	fmt.Printf("Fetching geo data from URL: %s\n", url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return models.GeoIPResponse{}, fmt.Errorf("failed to create request: %w", err)
	}

	fmt.Println("Using GeoIP Username:", GeoUsername)
	req.SetBasicAuth(GeoUsername, GeoAPIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return models.GeoIPResponse{}, fmt.Errorf("failed to connect to geoip service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body) // Read the body for debugging/logging purposes
		fmt.Printf("ERROR: Geoip response: %s\n", string(body))
		return models.GeoIPResponse{}, fmt.Errorf("geoip service returned status: %s", resp.Status)
	}

	var geoResp models.GeoIPResponse
	if err := json.NewDecoder(resp.Body).Decode(&geoResp); err != nil {
		return models.GeoIPResponse{}, fmt.Errorf("failed to decode geoip response: %w", err)
	}

	return geoResp, nil
}
