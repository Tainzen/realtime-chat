package controller

import (
	"encoding/json"
	"fmt"
	"github.com/Tainzen/realtime-chat/src/model"
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
// @Success 200 {object} dto.ResponseLogin "Success"
// @Failure 500 {object} dto.ErrorDTO "Internal Server Error"
// @Router /chat-rooms [post]
func (realTimeChatController *RealTimeChatController) CreateChatRoom(w http.ResponseWriter, r *http.Request) {

	//adding Content-type
	w.Header().Set("Content-Type", "application/json")

	var chatRoom model.ChatRoom

	// storing chatRoom
	err := json.NewDecoder(r.Body).Decode(&chatRoom)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message":"` + err.Error() + `"}`))
	}

	//create chat room
	result, err := realTimeChatRepository.CreateChatRoom(chatRoom)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message":"` + err.Error() + `"}`))
	}

	fmt.Println("Inserted a single document: ", result)
	// return the mongodb ID of generated document
	json.NewEncoder(w).Encode(result)

}
