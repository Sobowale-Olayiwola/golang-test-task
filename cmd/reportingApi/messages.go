package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func (app *application) getMessages(c *gin.Context) {
	sender, _ := c.Params.Get("sender")
	receiver, _ := c.Params.Get("receiver")
	key := fmt.Sprintf("%s-%s", strings.ToLower(sender), strings.ToLower(receiver))
	reverseKey := fmt.Sprintf("%s-%s", strings.ToLower(receiver), strings.ToLower(sender))

	results, err := app.models.Message.GetMessages(key, reverseKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"results": results})
}
