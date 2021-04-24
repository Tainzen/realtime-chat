package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ChatRoom model
type ChatRoom struct {
	ID   primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name string             `json:"name,omitempty" bson:"name,omitempty"`
}

// User model
type User struct {
	ID   primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserName string             `json:"username,omitempty" bson:"username,omitempty"`
	Password string `json:"password,omitempty" bson:"password,omitempty"`
}

// Message model
type Message struct{
	ID   primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	ChatRoomID string `json:"chatroom_id,omitempty" bson:"chatroom_id,omitempty"`
	UserID string `json:"user_id,omitempty" bson:"user_id,omitempty"`
	Body string `json:"body,omitempty" bson:"body,omitempty"`
}
