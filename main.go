package main

import (
	"github.com/WhereNext-co/WhereNext-Backend.git/initializers"
	"github.com/WhereNext-co/WhereNext-Backend.git/server"
)

// Binding from JSON
type RegisterBody struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Fullname string `json:"fullname" binding:"required"`
	Avatar   string `json:"avatar" binding:"required"`
}
// @title WhereNext API
// @version 1.0
// @description This is the API for the WhereNext application using Gin Framework.
// @host localhost:3000
// @BasePath /

func main() {
	initializers.LoadEnvVariables()
	server.InitServer()
}
