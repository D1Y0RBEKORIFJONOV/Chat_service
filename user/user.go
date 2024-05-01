package user

import (
	"errors"
	"fmt"
	"github.com/D1Y0RBEKORIFJONOV/Chat_service/chats"
	"github.com/D1Y0RBEKORIFJONOV/Chat_service/chats/message"
	"github.com/D1Y0RBEKORIFJONOV/Chat_service/postgres"
	"time"
)

type User struct {
	ID         int
	UserName   string
	Password   string
	Created_at time.Time
	Contact    []chats.Chat
	Group      []chats.Chat
	db         postgres.DB
}

func NewUser(username, password string) (*User, error) {
	db := postgres.DB{}
	var err error
	user := &User{}
	user.UserName = username
	user.Password = password

	user.ID, user.Created_at, err = db.UserInsert(username, password)
	if err != nil {
		return nil, err
	}
	user.db = db
	return user, nil
}

func (u *User) AddContact(contact_id int) error {
	err := u.db.Connect()
	if err != nil {
		return err
	}
	defer u.db.Close()
	query := fmt.Sprintf("SELECT chat_id FROM contacts WHERE user_id = %d AND contact_id = %d ", u.ID, contact_id)
	var d = 0
	err = u.db.DB.QueryRow(query).Scan(&d)
	if d != 0 {
		return errors.New("contact already exists")
	}

	ch, err := chats.NewChat("Persons", "Private")
	if err != nil {
		return err
	}
	err = u.db.ContactInsert(u.ID, contact_id, ch.Chat_id)
	if err != nil {

		return err
	}
	err = u.db.ContactInsert(contact_id, u.ID, ch.Chat_id)
	if err != nil {

		return err
	}
	mes, err := message.NewMessage(ch.Chat_id, u.ID, contact_id, "Hi,how are you :)", ch.Chat_type)
	ch.SendMessage(mes)
	u.Contact = append(u.Contact, *ch)

	return nil
}

func (u *User) SendMessageToReceiver(receiver_id int, message_text string) error {
	err := u.db.Connect()
	if err != nil {
		return err
	}
	defer u.db.Close()
	query := fmt.Sprintf("SELECT chat_id FROM contacts WHERE user_id = %d AND contact_id = %d ", u.ID, receiver_id)
	var chat_id = 0
	err = u.db.DB.QueryRow(query).Scan(&chat_id)
	if err != nil {
		return err
	}
	mss, err := message.NewMessage(chat_id, u.ID, receiver_id, message_text, "Private")
	if err != nil {
		return err
	}
	for _, ch := range u.Contact {
		if ch.Chat_id == chat_id {
			ch.SendMessage(mss)
			break
		}
	}

	return nil
}

func (u *User) CreateGroup(group_name string) error {
	chat, err := chats.NewChat(group_name, "Public")
	if err != nil {
		return err
	}
	mess := fmt.Sprintf("%s:Created the group", u.UserName)
	m, err := message.NewMessage(chat.Chat_id, u.ID, 0, mess, "Public")
	if err != nil {
		return err
	}
	chat.SendMessage(m)
	u.Group = append(u.Group, *chat)

	err = u.db.GroupInsert(u.ID, chat.Chat_id)
	if err != nil {
		return err
	}
	return nil
}

func (u *User) AddToGroup(chat_id, adding_user_id int) error {
	err := u.db.Connect()
	if err != nil {
		return err
	}
	defer u.db.Close()
	query := fmt.Sprintf("SELECT user_id FROM contacts WHERE chat_id = %d AND user_id = %d", chat_id, adding_user_id)
	var d = 0
	err = u.db.DB.QueryRow(query).Scan(&d)
	if d != 0 {
		return errors.New("User already added to group")
	}

	err = u.db.GroupInsert(adding_user_id, chat_id)
	if err != nil {
		return err
	}
	err = u.db.Connect()
	if err != nil {
		return err
	}
	defer u.db.Close()
	var name string
	err = u.db.DB.QueryRow("SELECT username FROM users WHERE user_id = $1", adding_user_id).Scan(&name)
	if err != nil {
		return err
	}
	fmt.Println(name, adding_user_id, err)
	str := fmt.Sprintf("%s:Added to group %s:", u.UserName, name)
	mess, err := message.NewMessage(chat_id, u.ID, 0, str, "Public")

	for _, ch := range u.Group {
		if ch.Chat_id == chat_id {

			ch.SendMessage(mess)
			break
		}
	}

	return nil
}

func (u *User) SendMessageToGroup(chat_id int, message_text string) error {
	mess, err := message.NewMessage(chat_id, u.ID, 0, message_text, "Public")
	if err != nil {
		return err
	}
	for _, ch := range u.Group {
		if ch.Chat_id == chat_id {
			ch.SendMessage(mess)
			break
		}
	}

	return nil
}

func ReadUser(username, password string) (User, error) {
	var ch *chats.Chat
	db := postgres.DB{}
	err := db.Connect()
	if err != nil {
		return User{}, err
	}
	defer db.Close()
	var user User
	err = db.DB.QueryRow("SELECT * FROM users WHERE username = $1 and password = $2 ", username, password).
		Scan(&user.ID, &user.UserName, &user.Password, &user.Created_at)

	rows, err := db.DB.Query("SELECT chat_id FROM contacts WHERE user_id  = $1; ", user.ID)
	if err != nil {
		return User{}, err
	}
	defer rows.Close()
	for rows.Next() {
		chat_id := 0
		err = rows.Scan(&chat_id)
		ch, err = chats.ReadChat(chat_id)
		if err != nil {
			return User{}, err
		}
		if ch.Chat_type == "Public" {
			user.Group = append(user.Group, *ch)
		} else {
			user.Contact = append(user.Contact, *ch)
		}
	}
	return user, nil
}
