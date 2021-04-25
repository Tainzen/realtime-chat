package controller

import (
	"encoding/json"
	"github.com/Tainzen/realtime-chat/src/controller/dto"
	"github.com/Tainzen/realtime-chat/src/model"
	"github.com/Tainzen/realtime-chat/src/repository"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

//RoomMap
var RoomMap = make(map[string]Room)

//Room to keep connections and broadcast message
type Room struct {
	Clients   map[*websocket.Conn]bool
	Broadcast chan dto.Message
}

//upgrader
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

//realTimeChatRepository to save into monogdb
var realTimeChatRepository = repository.RealTimeChatRepository{}

// RealTimeChatController - Structure
type RealTimeChatController struct{}

// HealthCheckPath - URL Path for health check
const HealthCheckPath = "/"

// HealthCheck Controller
// @Summary Health check API
// @Description Health check API
// @Produce json
// @Success 200 {object} dto.HealthCheckResponse "Success"
// @Failure 500 {object} dto.ErrorMessage "Internal Server Error"
// @Router / [get]
func (realTimeChatController *RealTimeChatController) HealthCheck(w http.ResponseWriter, r *http.Request) {

	//adding Content-type
	w.Header().Set("Content-Type", "application/json")

	response := dto.HealthCheckResponse{
		Message: "Realtime chat backend",
	}

	json.NewEncoder(w).Encode(response)
}

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
		errMessage.Message = "Error decoding request body"
		errMessage.Description = err.Error()
		json.NewEncoder(w).Encode(errMessage)
		return
	}

	//check if chat-room already exists
	count, err := realTimeChatRepository.CountChatRoomByChatName(chatRoom.Name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errMessage.Message = "Error counting chat-room"
		errMessage.Description = err.Error()
		json.NewEncoder(w).Encode(errMessage)
		return
	}

	//return if username already exists
	if count != 0 {
		w.WriteHeader(http.StatusBadRequest)
		errMessage.Message = "Chat-room already exists"
		errMessage.Description = "Chat-room already taken please try another name"
		json.NewEncoder(w).Encode(errMessage)
		return
	}

	// create chat room
	result, err := realTimeChatRepository.CreateChatRoom(chatRoom)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errMessage.Message = "Error creating chat-room"
		errMessage.Description = err.Error()
		json.NewEncoder(w).Encode(errMessage)
		return
	}

	// response message body
	response := dto.SuccessMessage{
		Message: "Chat room created successfully!",
		ID:      result,
	}

	json.NewEncoder(w).Encode(response)
}

// GetAllChatRoomsPath - URL Path to get all chat rooms
const GetAllChatRoomsPath = "/chat-rooms"

// GetAllChatRoom controller
// @Summary Get all chat rooms API
// @Description Get all chat room
// @Produce json
// @Success 200 {object} []model.ChatRoom "Success"
// @Failure 500 {object} dto.ErrorMessage "Internal Server Error"
// @Router /chat-rooms [get]
func (realTimeChatController *RealTimeChatController) GetAllChatRoom(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var errMessage dto.ErrorMessage

	// get chat-room by id
	result, err := realTimeChatRepository.FindAllChatRooms()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errMessage.Message = "Error creating chat-room"
		errMessage.Description = err.Error()
		json.NewEncoder(w).Encode(errMessage)
		return
	}

	json.NewEncoder(w).Encode(result)
}

// GetChatRoomPath - URL Path to get chat room by id
const GetChatRoomPath = "/chat-rooms/{room_id}"

// GetChatRoom controller
// @Summary Get chat room by id API
// @Description Get chat room by id
// @Param roomid path string true "room id"
// @Produce json
// @Success 200 {object} model.ChatRoom "Success"
// @Failure 500 {object} dto.ErrorMessage "Internal Server Error"
// @Router /chat-rooms [get]
func (realTimeChatController *RealTimeChatController) GetChatRoom(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var errMessage dto.ErrorMessage

	//ger paramaters
	param := mux.Vars(r)["room_id"]
	roomid, err := primitive.ObjectIDFromHex(param)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errMessage.Message = "Error while converting roomid to ObjectID"
		errMessage.Description = err.Error()
		json.NewEncoder(w).Encode(errMessage)
		return
	}

	// get chat room by id
	result, err := realTimeChatRepository.FindChatRoomByID(roomid)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errMessage.Message = "Error getting chat-room by id"
		errMessage.Description = err.Error()
		json.NewEncoder(w).Encode(errMessage)
		return
	}

	json.NewEncoder(w).Encode(result)
}

// UpdateChatRoomPath - URL Path to update chat room
const UpdateChatRoomPath = "/chat-rooms/{room_id}"

// UpdateChatRoom controller
// @Summary Update chat room API
// @Description Update new chat room and saves in mongo db
// @Param roomid path string true "room id"
// @Param ChatRoom body model.ChatRoom true "Request body Chat Room details"
// @Produce json
// @Success 200 {object} dto.SuccessMessage "Success"
// @Failure 500 {object} dto.ErrorMessage "Internal Server Error"
// @Router /chat-rooms [put]
func (realTimeChatController *RealTimeChatController) UpdateChatRoom(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var chatRoom model.ChatRoom
	var errMessage dto.ErrorMessage

	//get paramaters
	param := mux.Vars(r)["room_id"]
	roomid, err := primitive.ObjectIDFromHex(param)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errMessage.Message = "Error while converting roomid to ObjectID"
		errMessage.Description = err.Error()
		json.NewEncoder(w).Encode(errMessage)
		return
	}

	chatRoom.ID = roomid

	// storing chatRoom
	err = json.NewDecoder(r.Body).Decode(&chatRoom)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errMessage.Message = "Error decoding request body"
		errMessage.Description = err.Error()
		json.NewEncoder(w).Encode(errMessage)
		return
	}

	// update chat room
	_, err = realTimeChatRepository.UpdateChatRoom(chatRoom)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errMessage.Message = "Error creating chat-room"
		errMessage.Description = err.Error()
		json.NewEncoder(w).Encode(errMessage)
		return
	}

	response := dto.SuccessMessage{
		Message: "Chat room updated successfully!",
		ID:      roomid,
	}

	json.NewEncoder(w).Encode(response)

}

// DeleteChatRoomPath - URL Path to delete chat room by id
const DeleteChatRoomPath = "/chat-rooms/{room_id}"

// DeleteChatRoom controller
// @Summary Delete new chat room API
// @Description Delete chat room by id mongo db
// @Param roomid path string true "room id"
// @Produce json
// @Success 200 {object} dto.SuccessMessage "Success"
// @Failure 500 {object} dto.ErrorMessage "Internal Server Error"
// @Router /chat-rooms [delete]
func (realTimeChatController *RealTimeChatController) DeleteChatRoom(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var errMessage dto.ErrorMessage

	//get paramaters
	param := mux.Vars(r)["room_id"]
	roomid, err := primitive.ObjectIDFromHex(param)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errMessage.Message = "Error while converting roomid to ObjectID"
		errMessage.Description = err.Error()
		json.NewEncoder(w).Encode(errMessage)
		return
	}

	// update chat room
	_, err = realTimeChatRepository.DeleteChatRoom(roomid)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errMessage.Message = "Error deleteing chat-room"
		errMessage.Description = err.Error()
		json.NewEncoder(w).Encode(errMessage)
		return
	}

	response := dto.SuccessMessage{
		Message: "Chat room deleted successfully!",
		ID:      roomid,
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
// @Router /users [post]
func (realTimeChatController *RealTimeChatController) CreateUser(w http.ResponseWriter, r *http.Request) {

	//adding Content-type
	w.Header().Set("Content-Type", "application/json")

	var user model.User
	var errMessage dto.ErrorMessage

	//storing User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errMessage.Message = "Error decoding request body"
		errMessage.Description = err.Error()
		json.NewEncoder(w).Encode(errMessage)
		return
	}

	//check if username already exists
	count, err := realTimeChatRepository.CountUserByUsername(user.UserName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errMessage.Message = "Error counting user by username"
		errMessage.Description = err.Error()
		json.NewEncoder(w).Encode(errMessage)
		return
	}

	//return if username already exists
	if count != 0 {
		w.WriteHeader(http.StatusBadRequest)
		errMessage.Message = "Username already exists"
		errMessage.Description = "Username already taken please try another username"
		json.NewEncoder(w).Encode(errMessage)
		return
	}

	//generating hash to save password safely
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.MinCost)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errMessage.Message = "Error hashing password"
		errMessage.Description = err.Error()
		json.NewEncoder(w).Encode(errMessage)
		return
	}

	//storing as users password
	user.Password = string(hashBytes)

	//create user
	result, err := realTimeChatRepository.CreateUser(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errMessage.Message = "Error creating user"
		errMessage.Description = err.Error()
		json.NewEncoder(w).Encode(errMessage)
		return
	}

	//response message body
	response := dto.SuccessMessage{
		Message: "User created successfully!",
		ID:      result,
	}

	json.NewEncoder(w).Encode(response)
}

// GetUserPath - URL Path to get user by id
const GetUserPath = "/users/{uid}"

// GetUser controller
// @Summary Get user by id API
// @Description Get user by id
// @Param uid path string true "user id"
// @Produce json
// @Success 200 {object} dto.User "Success"
// @Failure 500 {object} dto.ErrorMessage "Internal Server Error"
// @Router /users [get]
func (realTimeChatController *RealTimeChatController) GetUser(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var errMessage dto.ErrorMessage

	//ger paramaters
	param := mux.Vars(r)["uid"]
	uid, err := primitive.ObjectIDFromHex(param)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errMessage.Message = "Error while converting uid to ObjectID"
		errMessage.Description = err.Error()
		json.NewEncoder(w).Encode(errMessage)
		return
	}

	// get user by id
	result, err := realTimeChatRepository.FindUserByID(uid)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errMessage.Message = "Error getting user by id"
		errMessage.Description = err.Error()
		json.NewEncoder(w).Encode(errMessage)
		return
	}

	u := dto.User{
		ID:        result.ID,
		UserName:  result.UserName,
		FirstName: result.FirstName,
		Lastname:  result.LastName,
	}

	json.NewEncoder(w).Encode(u)
}

// UpdateUserPath - URL Path to update user by id
const UpdateUserPath = "/users/{uid}"

// UpdateUser controller
// @Summary Update User API
// @Description Update user and saves in mongo db
// @Param userid path string true "user id"
// @Param User body dto.RequestUserUpdate true "Request body user details"
// @Produce json
// @Success 200 {object} dto.SuccessMessage "Success"
// @Failure 401 {object} dto.ErrorMessage "Wrong Password"
// @Failure 500 {object} dto.ErrorMessage "Internal Server Error"
// @Router /users/{uid} [put]
func (realTimeChatController *RealTimeChatController) UpdateUser(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var req dto.RequestUserUpdate
	var errMessage dto.ErrorMessage

	//get paramaters
	param := mux.Vars(r)["uid"]
	uid, err := primitive.ObjectIDFromHex(param)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errMessage.Message = "Error while converting uid to ObjectID"
		errMessage.Description = err.Error()
		json.NewEncoder(w).Encode(errMessage)
		return
	}

	// storing request user body
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errMessage.Message = "Error decoding request body"
		errMessage.Description = err.Error()
		json.NewEncoder(w).Encode(errMessage)
		return
	}

	// get user by id
	result, err := realTimeChatRepository.FindUserByID(uid)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errMessage.Message = "Error getting user by id"
		errMessage.Description = err.Error()
		json.NewEncoder(w).Encode(errMessage)
		return
	}

	//checking if old password is correct
	err = bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(req.OldPassword))
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		errMessage.Message = "Wrong old password"
		errMessage.Description = "Wrong old password, please insert the correct password"
		json.NewEncoder(w).Encode(errMessage)
		return
	}

	//generating hash to save password safely
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.MinCost)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errMessage.Message = "Error hashing password"
		errMessage.Description = err.Error()
		json.NewEncoder(w).Encode(errMessage)
		return
	}

	user := model.User{
		ID:        uid,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Password:  string(hashBytes),
	}

	// update user
	_, err = realTimeChatRepository.UpdateUser(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errMessage.Message = "Error creating chat-room"
		errMessage.Description = err.Error()
		json.NewEncoder(w).Encode(errMessage)
		return
	}

	response := dto.SuccessMessage{
		Message: "Chat room updated successfully!",
		ID:      uid,
	}

	json.NewEncoder(w).Encode(response)

}

// handleConnections handles all the connection on websockets
func handleConnections(ws *websocket.Conn, roomid string) (Room, error) {

	var room Room

	value, _ := RoomMap[roomid]

	room = value

	// Register our new client
	room.Clients[ws] = true

	for {
		var msg dto.Message

		// Read in a new message as JSON and map it to a Message object
		err := ws.ReadJSON(&msg)
		if err != nil {
			delete(room.Clients, ws)
			break
		}
		//save message in db
		rid, err := primitive.ObjectIDFromHex(roomid)
		if err != nil {
			return room, err
		}

		uid, err := primitive.ObjectIDFromHex(msg.UserID)
		if err != nil {
			return room, err
		}

		m := model.Message{
			UserID:     uid,
			ChatRoomID: rid,
			Body:       msg.Body,
		}

		// get user by id
		_, err = realTimeChatRepository.FindUserByID(uid)
		if err != nil {
			return room, err
		}

		//create message
		_, err = realTimeChatRepository.CreateMessage(m)
		if err != nil {
			return room, err
		}

		// Send the newly received message to the broadcast channel
		room.Broadcast <- msg

	}

	return room, nil

}

//handleMessages function that writes into connection
func handleMessages(room Room) {

	for {
		// Grab the next message from the broadcast channel
		msg := <-room.Broadcast
		// Send it out to every client that is currently connected
		for client := range room.Clients {
			err := client.WriteJSON(msg)
			if err != nil {
				client.Close()
				delete(room.Clients, client)
			}
		}
	}

}

const ChatRoomWebsocket = "/ws/chat-room/{room_id}"

// WebSocketHandler controller
// @Summary Websocket handler API
// @Description Websocket handler api to initiate websockets
// @Param roomid path string true "room id"
// @Param User body dto.Message true "Request body user id and message body"
// @Produce json
// @Success 200 {object} dto.SuccessMessage "Success"
// @Failure 500 {object} dto.ErrorMessage "Internal Server Error"
// @Router /ws/chat-room/{room_id} [get]
func (realTimeChatController *RealTimeChatController) WebSocketHandler(w http.ResponseWriter, r *http.Request) {

	var errMessage dto.ErrorMessage

	//get paramaters
	rid := mux.Vars(r)["room_id"]
	roomid, err := primitive.ObjectIDFromHex(rid)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errMessage.Message = "Error while converting roomid to ObjectID"
		errMessage.Description = err.Error()
		json.NewEncoder(w).Encode(errMessage)
		return
	}

	//get chat room by id to check if room id is present or not
	_, err = realTimeChatRepository.FindChatRoomByID(roomid)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errMessage.Message = "Error getting chat-room"
		errMessage.Description = err.Error()
		json.NewEncoder(w).Encode(errMessage)
		return
	}

	//resolve origin
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errMessage.Message = "Error while upgrading http connection to websockets"
		errMessage.Description = err.Error()
		json.NewEncoder(w).Encode(errMessage)
		return
	}

	defer conn.Close()

	//initailizing room
	var room Room
	_, ok := RoomMap[rid]

	room.Clients = map[*websocket.Conn]bool{}
	room.Broadcast = nil

	//register connection
	if ok == false {
		RoomMap[rid] = room
	}

	//handle connection
	cr, err := handleConnections(conn, rid)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errMessage.Message = "Error broadcasting message"
		errMessage.Description = err.Error()
		json.NewEncoder(w).Encode(errMessage)
		return
	}

	RoomMap[rid] = cr

	go handleMessages(cr)

}
