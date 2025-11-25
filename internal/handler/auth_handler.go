package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"prothomuse-server/internal/model"
	"prothomuse-server/internal/services"
	"prothomuse-server/internal/utils"
)

type AuthHandler struct {
	authService *services.AuthService
}

// NewAuthHandler creates a new instance of AuthHandler
func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// RegisterUser handles user registration
func (h *AuthHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	// Only allow POST requests
	if r.Method != http.MethodPost {
		sendErrorResponse(w, http.StatusMethodNotAllowed, "only POST method is allowed")
		return
	}

	var req model.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("error decoding register request: %v", err)
		sendErrorResponse(w, http.StatusBadRequest, "invalid request body")
		return
	}
	defer r.Body.Close()

	// Call the auth service to register the user
	user, err := h.authService.RegisterUser(req)
	if err != nil {
		log.Printf("error registering user: %v", err)
		sendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"data": map[string]interface{}{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"apiKey":   user.APIKey,
			"isActive": user.IsActive,
		},
		"message": "user registered successfully",
	})
}

// Login handles user login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	// Only allow POST requests
	if r.Method != http.MethodPost {
		sendErrorResponse(w, http.StatusMethodNotAllowed, "only POST method is allowed")
		return
	}

	var req model.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("error decoding login request: %v", err)
		sendErrorResponse(w, http.StatusBadRequest, "invalid request body")
		return
	}
	defer r.Body.Close()

	// Call the auth service to login the user
	response, err := h.authService.Login(req)
	if err != nil {
		log.Printf("error logging in user: %v", err)
		sendErrorResponse(w, http.StatusUnauthorized, err.Error())
		return
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"data":    response,
		"message": "user logged in successfully",
	})
}

// UpdateUser updates the authenticated user's profile.
// Requires Authorization: Bearer <token>
func (h *AuthHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut && r.Method != http.MethodPatch {
		sendErrorResponse(w, http.StatusMethodNotAllowed, "only PUT or PATCH method is allowed")
		return
	}

	token := extractJWTFromHeader(r)
	if token == "" {
		sendErrorResponse(w, http.StatusUnauthorized, "JWT token is required")
		return
	}
	claims, err := utils.ValidateJWT(token)
	if err != nil {
		log.Printf("error validating JWT token for update: %v", err)
		sendErrorResponse(w, http.StatusUnauthorized, "invalid or expired token")
		return
	}

	var req model.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("error decoding update user request: %v", err)
		sendErrorResponse(w, http.StatusBadRequest, "invalid request body")
		return
	}
	defer r.Body.Close()

	// ensure the update ID is the same as token subject
	req.ID = claims.UserID

	updated, err := h.authService.UpdateUser(req)
	if err != nil {
		log.Printf("error updating user: %v", err)
		sendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"data": map[string]interface{}{
			"id":       updated.ID,
			"username": updated.Username,
			"email":    updated.Email,
			"isActive": updated.IsActive,
		},
		"message": "user updated successfully",
	})
}

// ValidateAPIKey validates the API key from the request header
func (h *AuthHandler) ValidateAPIKey(w http.ResponseWriter, r *http.Request) {
	// Get API key from Authorization header
	apiKey := extractAPIKeyFromHeader(r)
	if apiKey == "" {
		sendErrorResponse(w, http.StatusUnauthorized, "API key is required")
		return
	}

	// Get user by API key
	user, err := h.authService.GetUserByAPIKey(apiKey)
	if err != nil {
		log.Printf("error validating API key: %v", err)
		sendErrorResponse(w, http.StatusUnauthorized, err.Error())
		return
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"data": map[string]interface{}{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"isActive": user.IsActive,
		},
		"message": "API key is valid",
	})
}

// ValidateJWT validates the JWT token from the request header
func (h *AuthHandler) ValidateJWT(w http.ResponseWriter, r *http.Request) {
	// Get JWT token from Authorization header
	token := extractJWTFromHeader(r)
	if token == "" {
		sendErrorResponse(w, http.StatusUnauthorized, "JWT token is required")
		return
	}

	// Validate the JWT token
	claims, err := utils.ValidateJWT(token)
	if err != nil {
		log.Printf("error validating JWT token: %v", err)
		sendErrorResponse(w, http.StatusUnauthorized, "invalid or expired token")
		return
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"data":    claims,
		"message": "JWT token is valid",
	})
}

// sendErrorResponse sends a JSON error response
func sendErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "error",
		"message": message,
	})
}

// extractAPIKeyFromHeader extracts the API key from the Authorization header
// Expected format: "ApiKey <api_key>"
func extractAPIKeyFromHeader(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) == 2 && parts[0] == "ApiKey" {
		return parts[1]
	}

	return ""
}

// extractJWTFromHeader extracts the JWT token from the Authorization header
// Expected format: "Bearer <token>"
func extractJWTFromHeader(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) == 2 && parts[0] == "Bearer" {
		return parts[1]
	}

	return ""
}
