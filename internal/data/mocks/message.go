package mocks

import (
	"errors"
	"twitch_chat_analysis/internal/data"
)

var mockMessage = data.Message{
	Sender:   "john",
	Receiver: "doe",
}

type MessageModel struct{}

func (m MessageModel) CreateMessage(message data.Message) error {
	if message.Sender == "error" {
		return errors.New("an error occured")
	}
	return nil
}

func (m MessageModel) GetMessages(key, reverseKey string) ([]data.Message, error) {
	if key == "not-found" {
		return nil, errors.New("redis: nil")
	}
	return []data.Message{mockMessage}, nil
}
