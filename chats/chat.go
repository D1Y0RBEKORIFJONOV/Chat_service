package chats

import (
	"github.com/D1Y0RBEKORIFJONOV/Chat_service/chats/message"
	"github.com/D1Y0RBEKORIFJONOV/Chat_service/postgres"
	"time"
)

type Chat struct {
	Chat_id   int
	Chat_name string
	Chat_type string
	Create_at time.Time
	Messages  []message.Message
	db        postgres.DB
}

func NewChat(chat_name, chat_type string) (*Chat, error) {
	db := postgres.DB{}
	var err error
	chat := Chat{}

	chat.Chat_id, chat.Create_at, err = db.ChatsInsert(chat_name, chat_type)
	if err != nil {
		return nil, err
	}
	chat.Chat_name = chat_name
	chat.Chat_type = chat_type
	chat.db = db
	return &chat, nil
}
func (c *Chat) SendMessage(message message.Message) {
	c.Messages = append(c.Messages, message)
}

func ReadChat(chat_id int) (*Chat, error) {
	chat := Chat{}
	db := postgres.DB{}
	err := db.Connect()
	if err != nil {
		return nil, err
	}
	err = db.DB.QueryRow("SELECT * FROM chats WHERE chat_id = $1", chat_id).
		Scan(&chat.Chat_id, &chat.Chat_name, &chat.Chat_type, &chat.Create_at)
	if err != nil {
		return nil, err
	}

	rows, err := db.DB.Query("SELECT * FROM messages WHERE chat_id = $1", chat.Chat_id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		mes := message.Message{}
		err = rows.Scan(&mes.Message_id, &mes.Chat_id, &mes.Sender_id, &mes.Receiver_id, &mes.Message_text, &mes.Sent_at)
		if err != nil {
			return nil, err
		}
		chat.Messages = append(chat.Messages, mes)
	}
	return &chat, nil
}
