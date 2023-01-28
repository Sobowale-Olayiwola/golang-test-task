package main

import "github.com/gin-gonic/gin"

func (app *application) routes() *gin.Engine {
	router := gin.Default()

	router.GET("/message/list/:sender/:receiver", app.getMessages)
	return router
}
