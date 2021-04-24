package controller

import (
	"encoding/json"
	"github.com/Tainzen/realtime-chat/src/model"
	"github.com/Tainzen/realtime-chat/src/controller/dto"
	"github.com/Tainzen/realtime-chat/src/repository"
	//"go.mongodb.org/mongo-driver/bson"
	// "go.mongodb.org/mongo-driver/bson/primitive"
	// "go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
	"github.com/gorilla/mux"
	"net/http"
)

//realTimeChatRepository to save into monogdb
var realTimeChatRepository = repository.RealTimeChatRepository{}

// RealTimeChatController - Structure
type RealTimeChatController struct{}

// CreateChatRoomPath - URL Path to create chat room
const CreateChatRoomPath = "/chat-rooms"

// CreateChatRoom controller
// @Summary Create new chat room API
// @Description Create new chat room and saves in mongo db
// @Param ChatRoom body model.ChatRoom true "Request body Chat Room details"
// @Produce json
// @Success 200 {object} dto.SuccessMessage "Success"
// @Failure 500 {object} dto.ErrorMessage "Internal Server Error"
// @Router /chat-rooms [post]
func (realTimeChatController *RealTimeChatController) CreateChatRoom(w http.ResponseWriter, r *http.Request) {

	//adding Content-type
	w.Header().Set("Content-Type", "application/json")

	var chatRoom model.ChatRoom
	var errMessage dto.ErrorMessage

	// storing chatRoom
	err := json.NewDecoder(r.Body).Decode(&chatRoom)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errMessage.Message="Error decoding request body"
		errMessage.Description=err.Error()
		json.NewEncoder(w).Encode(errMessage)
	}

	// create chat room
	result, err := realTimeChatRepository.CreateChatRoom(chatRoom)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errMessage.Message="Error creating chat-room"
		errMessage.Description=err.Error()
		json.NewEncoder(w).Encode(errMessage)
	}

	// response message body
	response:=dto.SuccessMessage{
		Message:"Chat room created successfully!",
		ID:result,
	}

	json.NewEncoder(w).Encode(response)
}

// CreateUserPath - URL Path to create user
const CreateUserPath = "/users"

// CreateUser controller
// @Summary Create new user API
// @Description Create new user and saves in mongo db
// @Param User body model.User true "Request body has user details"
// @Produce json
// @Success 200 {object} dto.SuccessMessage "Success"
// @Failure 400 {object} dto.ErrorMessage "Bad Request"
// @Failure 500 {object} dto.ErrorMessage "Internal Server Error"
// @Router /chat-rooms [post]
func (realTimeChatController *RealTimeChatController) CreateUser(w http.ResponseWriter, r *http.Request) {

	//adding Content-type
	w.Header().Set("Content-Type", "application/json")

	var user model.User
	var errMessage dto.ErrorMessage

	//storing User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errMessage.Message="Error decoding request body"
		errMessage.Description=err.Error()
		json.NewEncoder(w).Encode(errMessage)
	}

	//check if username already exists
	resp,err:=realTimeChatRepository.FindUserByUsername(user.UserName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errMessage.Message="Error finding user by username"
		errMessage.Description=err.Error()
		json.NewEncoder(w).Encode(errMessage)
	}

	//return if username already exists
	if resp.UserName==user.UserName{
		w.WriteHeader(http.StatusBadRequest)
		errMessage.Message="Username already exists"
		errMessage.Description="Username already taken please try another username"
		json.NewEncoder(w).Encode(errMessage)
		return
	}

	//generating hash to save password safely
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.MinCost)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errMessage.Message="Error hashing password"
		errMessage.Description=err.Error()
		json.NewEncoder(w).Encode(errMessage)
	}

	//storing as users password
	user.Password= string(hashBytes) 

	//create user
	result, err := realTimeChatRepository.CreateUser(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errMessage.Message="Error creating user"
		errMessage.Description=err.Error()
		json.NewEncoder(w).Encode(errMessage)
	}

	//response message body
	response:=dto.SuccessMessage{
		Message:"User created successfully!",
		ID:result,
	}

	json.NewEncoder(w).Encode(response)
}

// CreateMessagePath - URL Path to create message
const CreateMessagePath = "/users/{uid}/chat-rooms/{room_id}/messages"

// CreateMessage controller
// @Summary Create new message API
// @Description Create new message and saves in mongo db
// @Param Message body model.Message true "Request body has message details"
// @Produce json
// @Success 200 {object} dto.SuccessMessage "Success"
// @Failure 500 {object} dto.ErrorMessage "Internal Server Error"
// @Router /messages [post]
func (realTimeChatController *RealTimeChatController) CreateMessage(w http.ResponseWriter, r *http.Request) {

	//set header
	w.Header().Set("Content-Type", "application/json")
	
	//ger paramaters
	roomid := mux.Vars(r)["room_id"]
	uid := mux.Vars(r)["uid"]

	// //convert params to mongodb Hex ID
	// _id, err := primitive.ObjectIDFromHex(params) 
	// if err != nil {
	// 	fmt.Printf(err.Error())
	// }
	//opts := options.Delete().SetCollation(&options.Collation{})

	var message model.Message
	var errMessage dto.ErrorMessage

	//storing Message
	err := json.NewDecoder(r.Body).Decode(&message)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errMessage.Message="Error decoding request body"
		errMessage.Description=err.Error()
		json.NewEncoder(w).Encode(errMessage)
	}

	message.ChatRoomID=roomid
	message.UserID=uid

	//create message
	result, err := realTimeChatRepository.CreateMessage(message)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errMessage.Message="Error creating message"
		errMessage.Description=err.Error()
		json.NewEncoder(w).Encode(errMessage)
	}

	//response message body
	response:=dto.SuccessMessage{
		Message:"Message saved successfully!",
		ID:result,
	}

	json.NewEncoder(w).Encode(response)
}
