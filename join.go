package main

import (
	"log"
	"net/http"
)

func joinRoomHandler(w http.ResponseWriter, r *http.Request) {

	roomName := r.FormValue("newroomname")
	if roomName == "" {
		log.Printf("Room name nil error")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	room := newRoom(roomName)

	go room.run()

	return
}
