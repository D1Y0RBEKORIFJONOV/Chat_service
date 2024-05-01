package message

import (
	"github.com/D1Y0RBEKORIFJONOV/Chat_service/postgres"
	"time"
)

type Message struct {
	Message_id   int
	Chat_id      int
	Sender_id    int
	Receiver_id  int
	Message_text string
	Sent_at      time.Time
	db           postgres.DB
}

func NewMessage(chat_id, sender_id, receiver_id int, message_text, chat_type string) (Message, error) {
	db := postgres.DB{}
	var err error
	message := Message{}
	message.Chat_id = chat_id
	message.Sender_id = sender_id
	message.Receiver_id = receiver_id
	message.Message_text = message_text
	message.Message_id, message.Sent_at, err = db.MessageInsert(chat_id, sender_id, receiver_id, message_text, chat_type)
	if err != nil {
		return Message{}, err
	}
	message.db = db
	return message, nil
}
