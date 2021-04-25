package dto

// SuccessMessage dto
type HealthCheckResponse struct {
	Message string `json:"message"`
}

// ErrorMessage dto
type ErrorMessage struct {
	Message     string `json:"message"`
	Description string `json:"description"`
}

// SuccessMessage dto
type SuccessMessage struct {
	Message string      `json:"message"`
	ID      interface{} `json:"id"`
}

// User dto
type User struct {
	ID        interface{} `json:"_id"`
	UserName  string      `json:"username"`
	FirstName string      `json:"firstname"`
	Lastname  string      `json:"lastname"`
}

// RequestUserUpdate dto
type RequestUserUpdate struct {
	ID          interface{} `json:"_id"`
	FirstName   string      `json:"firstname" binding:"required"`
	LastName    string      `json:"lastname" binding:"required"`
	OldPassword string      `json:"oldpassword" binding:"required"`
	NewPassword string      `json:"newpassword" binding:"required"`
}

// Message dto
type Message struct {
	UserID     string `json:"user_id,omitempty" bson:"user_id,omitempty"`
	Body       string `json:"body,omitempty" bson:"body,omitempty"`
}
