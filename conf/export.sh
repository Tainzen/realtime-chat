#!/bin/sh

export SVR_BASEPATH=/realtime-chat/api/v1
export SVR_PORT=:8081
export DB_DRIVER=mongodb
export DB_NAME=realtime_chat
export DB_USER=admin
export DB_PASSWORD=password 
export DB_HOST=localhost
export DB_PORT=27017

bin/server