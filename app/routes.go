package app

import (
	"github.com/Ilhamokta74/Gotoko/app/controllers"
	"github.com/gorilla/mux"
)

func (server *Server) InitializeRoutes() {
	server.Router = mux.NewRouter()
	server.Router.HandleFunc("/", controllers.Home).Methods("GET")
}
