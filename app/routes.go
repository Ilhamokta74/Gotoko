package app

import "github.com/Ilhamokta74/Gotoko/app/controllers"

func (server *Server) InitializeRoutes() {
	server.Router.HandleFunc("/", controllers.Home).Methods("GET")
}
