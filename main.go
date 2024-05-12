package main

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

/*
Create a Customer struct

Each customer includes:
- ID
- Name
- Role
- Email
- Phone
- Contacted (i.e., indication of whether or not the customer has been contacted)
Data is mapped to logical, appropriate types (e.g., Name should not be a bool).
*/
type Customer struct {
	Id          uint64 `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Role        string `json:"role,omitempty"`
	Email       string `json:"email,omitempty"`
	Phone       string `json:"phone,omitempty"`
	IsContacted bool   `json:"contacted,omitempty"`
}

/*
Create a mock "database" to store customer data
Customers are stored appropriately in a basic data structure (e.g., slice, map, etc.) that represents a "database."

Seed the database with initial customer data
The "database" data structure is non-empty. That is, prior to any CRUD operations performed by the user (e.g., adding a customer),
the database includes at least three existing (i.e., "hard-coded") customers.

Assign unique IDs to customers in the database
Customers in the database have unique ID values (i.e., no two customers have the same ID value).
*/
var customers = map[uint64]Customer{
	1: {
		Id:          1,
		Name:        "John Doe",
		Role:        "Admin",
		Email:       "john.doe@gmail.com",
		Phone:       "1234567890",
		IsContacted: false,
	},
	2: {
		Id:          2,
		Name:        "Jane Doe",
		Role:        "User",
		Email:       "jane.doe@gmail.com",
		Phone:       "0987654321",
		IsContacted: false,
	},
	3: {
		Id:          3,
		Name:        "John Smith",
		Role:        "User",
		Email:       "john.smith@gmail.com",
		Phone:       "1234567890",
		IsContacted: false,
	},
}

// Function to read customer data from the request body
func readCustomerFromRequestBody(r *http.Request) (Customer, error) {
	// Read request data
	// The application leverages the io/ioutil package to read I/O (e.g., request) data.
	// Read the request body into a byte slice
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return Customer{}, err
	}

	// Parse JSON data
	// The applications leverages the encoding/json package to parse JSON data.
	// Unmarshal the byte slice into a Customer struct
	var customer Customer
	err = json.Unmarshal(body, &customer)
	if err != nil {
		return Customer{}, err
	}

	return customer, nil
}

// Create a function to get a slice of all customers
// The application includes a function that returns a slice of all customers in the "database."
func getCustomerSlices() []Customer {
	result := make([]Customer, 0, len(customers))

	for _, value := range customers {
		result = append(result, value)
	}

	return result
}

// Create and assign handlers for requests

// Getting all customers through a the /customers path
// The application returns all customers in the "database" when a GET request is made to the /customers path.
func getCustomers(w http.ResponseWriter, r *http.Request) {
	// Set headers to indicate the proper media type
	// An appropriate Content-Type header is sent in server responses.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(getCustomerSlices())
}

// Getting a single customer through a /customers/{id} path
// The application returns a single customer when a GET request is made to the /customers/{id} path.
func getCustomer(w http.ResponseWriter, r *http.Request) {
	// Set headers to indicate the proper media type
	// An appropriate Content-Type header is sent in server responses.
	w.Header().Set("Content-Type", "application/json")

	urlPathVars := mux.Vars(r)
	if givenId, exist := urlPathVars["id"]; exist {
		id, err := strconv.ParseUint(givenId, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode("Failed to parse the given id: " + givenId + ". " + err.Error())
			return
		}

		// Includes basic error handling for non-existent customers
		if customer, ok := customers[id]; ok {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(customer)
		} else {
			// If the user queries for a customer that doesn't exist (i.e., when getting a customer, updating a customer, or deleting a customer), the server response includes:
			// A 404 status code in the header
			w.WriteHeader(http.StatusNotFound)
			// null or an empty JSON object literal or an error message
			json.NewEncoder(w).Encode("Customer with id: " + givenId + " not found.")
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

// Creating a customer through a /customers path
// The application adds a new customer to the "database" when a POST request is made to the /customers path.
func addCustomer(w http.ResponseWriter, r *http.Request) {
	// Set headers to indicate the proper media type
	// An appropriate Content-Type header is sent in server responses.
	w.Header().Set("Content-Type", "application/json")

	// Decode the request body into a Golang value (a Customer struct)
	customer, err := readCustomerFromRequestBody(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Failed to read customer data from the request body. " + err.Error())
		return
	}

	if customer.Id == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Invalid customer id 0. Customer id should be greater than 0.")
		return
	}

	if _, exist := customers[customer.Id]; exist {
		// If the key already exists in the "database", update the HTTP status with a "407 Conflict" message
		// In such a case, the original "database" is not updated at all
		w.WriteHeader(http.StatusConflict)
	} else {
		// If the key doesn't exist, add it to the "database" and return a "201 Created" in the header
		customers[customer.Id] = customer
		w.WriteHeader(http.StatusCreated)
	}

	// Regardless of successful resource creation or not, return the current state of the "database"
	json.NewEncoder(w).Encode(getCustomerSlices())
}

// Updating a customer through a /customers/{id} path
// The application updates an existing customer in the "database" when a PUT request is made to the /customers/{id} path.
func updateCustomer(w http.ResponseWriter, r *http.Request) {
	// Set headers to indicate the proper media type
	// An appropriate Content-Type header is sent in server responses.
	w.Header().Set("Content-Type", "application/json")

	urlPathVars := mux.Vars(r)
	if givenId, exist := urlPathVars["id"]; exist {
		id, err := strconv.ParseUint(givenId, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode("Failed to parse the given id: " + givenId + ". " + err.Error())
			return
		}

		customer, err := readCustomerFromRequestBody(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode("Failed to read customer data from the request body. " + err.Error())
			return
		}

		if customer.Id != id {
			w.WriteHeader(http.StatusBadRequest)
			// An empty JSON object in the response body
			json.NewEncoder(w).Encode("Customer id in the request body does not match the id in the URL path.")
			return
		}

		// Includes basic error handling for non-existent customers
		if _, ok := customers[id]; ok {
			customers[id] = customer
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(getCustomerSlices())
		} else {
			// If the user queries for a customer that doesn't exist (i.e., when getting a customer, updating a customer, or deleting a customer), the server response includes:
			// A 404 status code in the header
			w.WriteHeader(http.StatusNotFound)
			// null or an empty JSON object literal or an error message
			json.NewEncoder(w).Encode("Customer with id: " + givenId + " not found.")
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

// Deleting a customer through a /customers/{id} path
// The application deletes an existing customer from the "database" when a DELETE request is made to the /customers/{id} path.
func deleteCustomer(w http.ResponseWriter, r *http.Request) {
	// Set headers to indicate the proper media type
	// An appropriate Content-Type header is sent in server responses.
	w.Header().Set("Content-Type", "application/json")

	urlPathVars := mux.Vars(r)
	if givenId, exist := urlPathVars["id"]; exist {
		id, err := strconv.ParseUint(givenId, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode("Failed to parse the given id: " + givenId + ". " + err.Error())
			return
		}

		// Includes basic error handling for non-existent customers
		if _, ok := customers[id]; ok {
			delete(customers, id)
			w.WriteHeader(http.StatusOK)
			// return the current state of the "database"
			json.NewEncoder(w).Encode(getCustomerSlices())
		} else {
			// If the user queries for a customer that doesn't exist (i.e., when getting a customer, updating a customer, or deleting a customer), the server response includes:
			// A 404 status code in the header
			w.WriteHeader(http.StatusNotFound)
			// null or an empty JSON object literal or an error message
			json.NewEncoder(w).Encode("Customer with id: " + givenId + " not found.")
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func main() {

	// Set up and configure a router
	// The application uses a router (e.g., gorilla/mux, http.ServeMux, etc.) that supports HTTP method-based routing and variables in URL paths.
	router := mux.NewRouter()

	// Create RESTful server endpoints for CRUD operations
	// The application handles the following 5 operations for customers in the "database":
	// Each RESTful route is associated with the correct HTTP verb.

	// Getting a single customer through a /customers/{id} path
	router.HandleFunc("/customers/{id}", getCustomer).Methods("GET")

	// Getting all customers through a the /customers path
	router.HandleFunc("/customers", getCustomers).Methods("GET")

	// Creating a customer through a /customers path
	router.HandleFunc("/customers", addCustomer).Methods("POST")

	// Updating a customer through a /customers/{id} path
	router.HandleFunc("/customers/{id}", updateCustomer).Methods("PUT")

	// Deleting a customer through a /customers/{id} path
	router.HandleFunc("/customers/{id}", deleteCustomer).Methods("DELETE")

	// Serve static HTML at the home ("/") route
	// The home route is a client API endpoint, and includes a brief overview of the API (e.g., available endpoints). Note: This is the only route that does not return a JSON response.
	fileServer := http.FileServer(http.Dir("./static"))
	router.Handle("/", fileServer)

	// Serve the API locally
	// The API can be accessed via localhost.
	http.ListenAndServe(":3000", router)
}
