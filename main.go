package main

import (
	//"os"
	
	"github.com/Tainzen/realtime-chat/src/controller"
	"github.com/gorilla/mux"
	
	"net/http"
)

// @title RealTime-Chat Microservice
// @version 1
// @description This microservice serves as Realtime chat backend

// @BasePath /realtime-chat/api/v1
func main() {
	route := mux.NewRouter()

	// Instantiate controllers
	realTimeChatController := controller.RealTimeChatController{}

	//basepath:=os.Getenv("/realtime-chat")
	api := route.PathPrefix("/realtime-chat/api/v1").Subrouter()
	api.HandleFunc(controller.CreateChatRoomPath, realTimeChatController.CreateChatRoom).Methods("POST")
	api.HandleFunc(controller.CreateUserPath, realTimeChatController.CreateUser).Methods("POST")
	api.HandleFunc(controller.CreateMessagePath, realTimeChatController.CreateMessage).Methods("POST")
	http.ListenAndServe(":8081", api)
}



