package dto

// ErrorMessage dto
type ErrorMessage struct {
	Message     string `json:"message"`
	Description string `json:"description"`
}

// SuccessMessage dto
type SuccessMessage struct {
	Message string `json:"message"`
	ID      interface{} `json:"id"`
}
