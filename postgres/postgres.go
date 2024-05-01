package postgres

import (
	"fmt"
	"time"
)

func IsHaveInDataBase(table, soursObj, like string) (bool, error) {
	var date string
	db := DB{}
	err := db.Connect()
	if err != nil {
		return false, err
	}

	query := fmt.Sprintf("SELECT  %s FROM %s WHERE %s =  '%s';", soursObj, table, soursObj, like)
	defer db.Close()
	row := db.DB.QueryRow(query)
	err = row.Scan(&date)
	if err != nil {
		return false, nil
	}
	if date == "" {
		return false, nil
	}
	return date == like, nil
}

func (db *DB) UserInsert(userName string, password string) (int, time.Time, error) {
	ok, err := IsHaveInDataBase("users", "username", userName)
	if err != nil {
		return 0, time.Time{}, err
	}
	if ok {
		return 0, time.Time{}, fmt.Errorf("User %s does already exist", userName)
	}
	var (
		id int
		t  time.Time
	)

	err = db.Connect()
	if err != nil {
		return 0, time.Time{}, err
	}
	defer db.Close()
	query := `
	INSERT INTO users(username, password) VALUES ($1, $2) 
	RETURNING user_id, created_at ;
`
	err = db.DB.QueryRow(query, userName, password).Scan(&id, &t)
	if err != nil {
		return 0, time.Time{}, err
	}
	return id, t, err
}

func (db *DB) ChatsInsert(chat_name, chat_type string) (int, time.Time, error) {
	err := db.Connect()
	if err != nil {
		return 0, time.Time{}, err
	}
	defer db.Close()

	queyry := `
	INSERT INTO chats(chat_name, chat_type) VALUES ($1, $2)
	RETURNING chat_id, created_at ;
`
	var chat_id int
	var t time.Time
	err = db.DB.QueryRow(queyry, chat_name, chat_type).Scan(&chat_id, &t)

	if err != nil {
		fmt.Println(err)
		return 0, time.Time{}, err
	}
	return chat_id, t, nil
}

func (db *DB) MessageInsert(chat_id, sender_id, recipient_id int, message_text, chat_type string) (int, time.Time, error) {
	err := db.Connect()
	if err != nil {
		return 0, time.Time{}, err
	}
	defer db.Close()
	var message_id = 0
	var t time.Time
	query := ""
	if chat_type == "Public" {
		query = `
		INSERT INTO messages(chat_id, sender_id,recipient_id, message_text) VALUES ($1, $2,0, $3) 
		RETURNING chat_id, sent_at ;
		;
`
	} else if chat_type == "Private" {
		query = `
	INSERT INTO messages(chat_id, sender_id,recipient_id, message_text ) VALUES ($1, $2, $3, $4)
	RETURNING chat_id, sent_at ;
	;
`
	}

	if chat_type == "Public" {
		err = db.DB.QueryRow(query, chat_id, sender_id, message_text).Scan(&message_id, &t)
	} else if chat_type == "Private" {
		err = db.DB.QueryRow(query, chat_id, sender_id, recipient_id, message_text).Scan(&message_id, &t)
	}
	if err != nil {
		fmt.Println(err)
	}
	return message_id, t, err
}

func (db *DB) ContactInsert(user_id, contact_id, chat_id int) error {
	err := db.Connect()
	if err != nil {
		return err
	}

	defer db.Close()

	query := `
	INSERT INTO contacts(user_id, contact_id, chat_id) VALUES ($1, $2, $3) ;
	`
	_, err = db.DB.Exec(query, user_id, contact_id, chat_id)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) GroupInsert(user_id, chat_id int) error {
	err := db.Connect()
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()
	query := `
	INSERT INTO contacts(user_id,contact_id, chat_id) VALUES ($1,0,$2) ;
	`
	_, err = db.DB.Exec(query, user_id, chat_id)
	if err != nil {
		return err
	}

	return nil
}
