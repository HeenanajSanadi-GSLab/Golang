package errorhandling
import (

	"../models"
	
)
type Response struct {
   StatusCode int `json:"statusCode"`              
   ErrorMessage string `json:"errorMessage"`               
   Payload []models.User `json:"payload"`
}