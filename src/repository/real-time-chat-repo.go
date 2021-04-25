package repository

import (
	"context"
	"github.com/Tainzen/realtime-chat/src/model"
	"github.com/Tainzen/realtime-chat/utils/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// chat room collection
var chatRoomCollection = database.Db().Database("realtime_chat").Collection("chat_rooms")

// user collection
var userCollection = database.Db().Database("realtime_chat").Collection("users")

// message collection
var messageCollection = database.Db().Database("realtime_chat").Collection("messages")

// realTimeChatRepository - Structure
type RealTimeChatRepository struct{}

// CreateChatRoom - Inserts chat room into db
func (realTimeChat *RealTimeChatRepository) CreateChatRoom(chatRoom model.ChatRoom) (interface{}, error) {
	//insert into mongodb
	result, err := chatRoomCollection.InsertOne(context.TODO(), chatRoom)
	if err != nil {
		return nil, err
	}

	return result.InsertedID, nil
}

// FindAllChatRooms - Find all chat-rooms
func (realTimeChat *RealTimeChatRepository) FindAllChatRooms() ([]model.ChatRoom, error) {

	var chatRooms []model.ChatRoom
	//find all chat-rooms
	cur, err := chatRoomCollection.Find(context.TODO(), bson.D{{}})
	if err != nil {
		return nil, err
	}

	for cur.Next(context.TODO()) {

		var chatRoom model.ChatRoom
		err := cur.Decode(&chatRoom)
		if err != nil {
			return nil, err
		}

		//appending chatRooms
		chatRooms = append(chatRooms, chatRoom)
	}

	return chatRooms, nil
}

// FindChatRoomByID - Find chat room by id
func (realTimeChat *RealTimeChatRepository) FindChatRoomByID(id primitive.ObjectID) (model.ChatRoom, error) {

	var chatRoom model.ChatRoom
	//find chat-room with id
	err := chatRoomCollection.FindOne(context.TODO(), bson.D{{"_id", id}}).Decode(&chatRoom)
	if err != nil {
		return chatRoom, err
	}

	return chatRoom, nil
}

// UpdateChatRoom - Updates chat room into db
func (realTimeChat *RealTimeChatRepository) UpdateChatRoom(chatRoom model.ChatRoom) (model.ChatRoom, error) {

	var room model.ChatRoom
	//filter
	filter := bson.D{{"_id", chatRoom.ID}}

	//to return updated document
	after := options.After

	returnOpt := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
	}

	update := bson.D{{"$set", bson.D{{"name", chatRoom.Name}}}}

	err := chatRoomCollection.FindOneAndUpdate(context.TODO(), filter, update, &returnOpt).Decode(&room)
	if err != nil {
		return room, err
	}

	return room, nil
}

// DeleteChatRoom - Deletes chat room by id from db
func (realTimeChat *RealTimeChatRepository) DeleteChatRoom(id primitive.ObjectID) (interface{}, error) {

	//options
	opts := options.Delete().SetCollation(&options.Collation{})
	//filter by id
	filter := bson.D{{"_id", id}}
	res, err := chatRoomCollection.DeleteOne(context.TODO(), filter, opts)
	if err != nil {
		return res, err
	}

	return res, nil
}

// CountChatByChatName - Counts chat-room by username into db
func (realTimeChat *RealTimeChatRepository) CountChatRoomByChatName(name string) (int64, error) {

	//filter by name
	filter := bson.D{{"name", name}}
	count, err := userCollection.CountDocuments(context.TODO(), filter)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// CountChatByID - Counts chat-room by id into db
func (realTimeChat *RealTimeChatRepository) CountChatRoomByID(id primitive.ObjectID) (int64, error) {

	//filter by id
	filter := bson.D{{"_id", id}}
	count, err := userCollection.CountDocuments(context.TODO(), filter)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// CreateUser - Inserts user into db
func (realTimeChat *RealTimeChatRepository) CreateUser(user model.User) (interface{}, error) {
	//insert into mongodb
	result, err := userCollection.InsertOne(context.TODO(), user)
	if err != nil {
		return nil, err
	}

	return result.InsertedID, nil
}

// FindUserByID - Find user by id
func (realTimeChat *RealTimeChatRepository) FindUserByID(id primitive.ObjectID) (model.User, error) {

	var user model.User
	//find user with id
	err := userCollection.FindOne(context.TODO(), bson.D{{"_id", id}}).Decode(&user)
	if err != nil {
		return user, err
	}

	return user, nil
}

// UpdateUser - Updates user into db
func (realTimeChat *RealTimeChatRepository) UpdateUser(u model.User) (model.User, error) {

	var user model.User
	//filter
	filter := bson.M{"_id": u.ID}

	//to return updated document
	after := options.After

	returnOpt := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
	}

	update := bson.M{"$set": bson.M{"firstname": u.FirstName, "lastname": u.LastName, "password": u.Password}}

	err := userCollection.FindOneAndUpdate(context.TODO(), filter, update, &returnOpt).Decode(&user)
	if err != nil {
		return user, err
	}
	return user, nil
}

// FindUser - Finds user by username into db
func (realTimeChat *RealTimeChatRepository) FindUserByUsername(username string) (model.User, error) {
	var result model.User

	//filter by username
	filter := bson.D{{"username", username}}
	err := userCollection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		return result, err
	}

	return result, nil
}

// CountUserByUsername - Counts user by username into db
func (realTimeChat *RealTimeChatRepository) CountUserByUsername(username string) (int64, error) {

	//filter by username
	filter := bson.D{{"username", username}}
	count, err := userCollection.CountDocuments(context.TODO(), filter)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// CreateMessage - Inserts message into db
func (realTimeChat *RealTimeChatRepository) CreateMessage(message model.Message) (interface{}, error) {
	//insert into mongodb
	result, err := messageCollection.InsertOne(context.TODO(), message)
	if err != nil {
		return nil, err
	}

	return result.InsertedID, nil
}
