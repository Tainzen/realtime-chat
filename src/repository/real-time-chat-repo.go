package repository

import (
	"context"
	"github.com/Tainzen/realtime-chat/src/model"
	"github.com/Tainzen/realtime-chat/utils/database"
)

// get chat rooms collection
var chatRoomsCollection = database.Db().Database("realtime_chat").Collection("chat_rooms")

// realTimeChatRepository - Structure
type RealTimeChatRepository struct{}

//CreateChatRoom - Inserts chat room into db
func (realTimeChat *RealTimeChatRepository) CreateChatRoom(chatRoom model.ChatRoom) (interface{}, error) {
	//insert into mongodb
	result, err := chatRoomsCollection.InsertOne(context.TODO(), chatRoom)
	if err != nil {
		return nil,err
	}

	return result.InsertedID, nil

}
