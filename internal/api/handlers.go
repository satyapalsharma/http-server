package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// User represents a user entity in our API.
type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// ErrorResponse represents a standardized error response structure for the API.
type ErrorResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// respondJSON is a helper function to write JSON responses.
// It sets the Content-Type header, writes the status code, and marshals the payload to JSON.
func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		// Log the error but don't expose internal marshalling errors to the client directly.
		log.Printf("Error marshalling JSON response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(response)
	if err != nil {
		log.Printf("Error writing JSON response: %v", err)
	}
}

// respondError is a helper function to write standardized JSON error responses.
func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, ErrorResponse{Message: message, Code: status})
}

// GetHealthHandler provides a simple health check endpoint for the API.
// It returns a 200 OK with a status message.
//
// Path: /api/health
// Method: GET
func GetHealthHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received health check request from %s", r.RemoteAddr)
	respondJSON(w, http.StatusOK, map[string]string{"status": "healthy", "service": "api"})
}

// GetUsersHandler retrieves a list of dummy users.
// In a real application, this would fetch data from a database or service.
//
// Path: /api/users
// Method: GET
func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received request to get all users from %s", r.RemoteAddr)

	// Simulate fetching users from a data store
	users := []User{
		{ID: "usr-001", Name: "Alice Smith", Email: "alice@example.com"},
		{ID: "usr-002", Name: "Bob Johnson", Email: "bob@example.com"},
		{ID: "usr-003", Name: "Charlie Brown", Email: "charlie@example.com"},
	}

	respondJSON(w, http.StatusOK, users)
}

// CreateUserHandler handles the creation of a new user.
// It expects a JSON payload in the request body, validates it, and "creates" the user.
//
// Path: /api/users
// Method: POST
func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received request to create user from %s", r.RemoteAddr)

	// Ensure the request method is POST
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed. Only POST is supported.")
		return
	}

	var newUser User
	// Decode the JSON request body into the newUser struct
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		log.Printf("Error decoding request body: %v", err)
		respondError(w, http.StatusBadRequest, fmt.Sprintf("Invalid request payload: %v", err))
		return
	}

	// Basic server-side validation
	if newUser.Name == "" {
		respondError(w, http.StatusBadRequest, "User name is required.")
		return
	}
	if newUser.Email == "" {
		respondError(w, http.StatusBadRequest, "User email is required.")
		return
	}
	// A more robust validation would include email format, uniqueness, etc.

	// In a real application:
	// 1. Generate a unique ID (e.g., UUID)
	// 2. Save the user to a database
	// 3. Handle potential database errors (e.g., duplicate email)
	// For this example, we'll just assign a dummy ID and log it.
	newUser.ID = fmt.Sprintf("usr-%d", len(newUser.Name)+len(newUser.Email)) // Simple dummy ID generation

	log.Printf("Successfully 'created' new user: %+v", newUser)

	// Respond with the created user and a 201 Created status
	respondJSON(w, http.StatusCreated, newUser)
}

// NotFoundHandler provides a generic 404 handler for API routes.
// This can be used by the router for any API path that doesn't match an explicit handler.
//
// Path: Any unmatched /api/* path
// Method: Any
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("API Not Found: %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
	respondError(w, http.StatusNotFound, fmt.Sprintf("API endpoint not found: %s %s", r.Method, r.URL.Path))
}<ctrl63>