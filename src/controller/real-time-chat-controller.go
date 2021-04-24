package controller

import (
	"encoding/json"
	"github.com/Tainzen/realtime-chat/src/model"
	"github.com/Tainzen/realtime-chat/src/controller/dto"
	"github.com/Tainzen/realtime-chat/src/repository"
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
