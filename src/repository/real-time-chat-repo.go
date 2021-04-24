package repository

import (
	"go.mongodb.org/mongo-driver/bson"
	"context"
	"github.com/Tainzen/realtime-chat/src/model"
	"github.com/Tainzen/realtime-chat/utils/database"
)

// chat room collection
var chatRoomCollection = database.Db().Database("realtime_chat").Collection("chat_rooms")

// user collection
var userCollection = database.Db().Database("realtime_chat").Collection("users")

// message collection
var  messageCollection = database.Db().Database("realtime_chat").Collection("messages")

// realTimeChatRepository - Structure
type RealTimeChatRepository struct{}

// CreateChatRoom - Inserts chat room into db
func (realTimeChat *RealTimeChatRepository) CreateChatRoom(chatRoom model.ChatRoom) (interface{}, error) {
	//insert into mongodb
	result, err := chatRoomCollection.InsertOne(context.TODO(), chatRoom)
	if err != nil {
		return nil,err
	}

	return result.InsertedID, nil
}

// CreateUser - Inserts user into db
func (realTimeChat *RealTimeChatRepository) CreateUser(user model.User) (interface{}, error) {
	//insert into mongodb
	result, err := userCollection.InsertOne(context.TODO(), user)
	if err != nil {
		return nil,err
	}

	return result.InsertedID, nil
}

// FindUser - Finds user by username into db
func (realTimeChat *RealTimeChatRepository) FindUserByUsername(username string) (model.User, error) {
	var result model.User

	//filter by username
	filter:=bson.D{{"username", username}}
	err := userCollection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		return result,err
	}

	return result, nil
}

// CreateMessage - Inserts message into db
func (realTimeChat *RealTimeChatRepository) CreateMessage(message model.Message) (interface{}, error) {
	//insert into mongodb
	result, err := messageCollection.InsertOne(context.TODO(), message)
	if err != nil {
		return nil,err
	}

	return result.InsertedID, nil
}
