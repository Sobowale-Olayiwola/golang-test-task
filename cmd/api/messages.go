package main

import (
	"net/http"
	"strings"
	"time"
	"twitch_chat_analysis/internal/data"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func inputValidator(input interface{}) error {
	validate = validator.New()
	err := validate.Struct(input)
	if err != nil {
		return err
	}
	return nil
}

func (app *application) createMessage(c *gin.Context) {
	var input struct {
		Sender   string `json:"sender" validate:"required"`
		Receiver string `json:"receiver" validate:"required"`
		Message  string `json:"message" validate:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, "Bad Request")
		return
	}
	inputErr := inputValidator(input)
	if inputErr != nil {
		c.JSON(http.StatusBadRequest, inputErr.Error())
		return
	}

	message := data.Message{
		Sender:    strings.ToLower(input.Sender),
		Receiver:  strings.ToLower(input.Receiver),
		Message:   input.Message,
		CreatedAt: time.Now(),
	}

	if err := app.models.Message.CreateMessage(message); err != nil {
		c.JSON(http.StatusBadRequest, "Bad Request")
		return
	}
	c.JSON(http.StatusOK, "OK")
}
