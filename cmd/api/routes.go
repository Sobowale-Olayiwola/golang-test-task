package main

import "github.com/gin-gonic/gin"

func (app *application) routes() *gin.Engine {
	router := gin.Default()
	router.POST("/message", app.createMessage)

	return router
}
