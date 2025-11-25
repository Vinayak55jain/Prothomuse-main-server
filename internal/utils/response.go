package utils
import(
	"net/http"
	"encoding/json"
)
type Response struct {
	
Success bool        `json:"success"`
Message string      `json:"message,omitempty"`
Data   interface{} `json:"data,omitempty"`
Error string      `json:"error,omitempty"`
}

func JSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(
		Response{
			Success: statusCode  <400,
			Data:    data,
		})
}

func ErrorResponse(w http.ResponseWriter, statusCode int, errMsg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(
		Response{
			Success: false,
			Error:   errMsg,
		})
}

func SuccessResponse(w http.ResponseWriter, statusCode int, message string,data interface{}	) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(
		Response{
			Success: true,
			Message: message,
			Data:    data,
		})
}