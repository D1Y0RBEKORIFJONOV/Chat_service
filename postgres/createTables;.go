package postgres

func UserCreateTables() error {
	db := DB{}
	err := db.Connect()
	if err != nil {
		return err
	}
	defer db.Close()
	query := `
	CREATE TABLE IF NOT EXISTS users (user_id SERIAL NOT NULL PRIMARY KEY,
	username  VARCHAR(255) NOT NULL,
	password VARCHAR(255) NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
`
	_, err = db.DB.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func ChatsCreateTable() error {
	db := DB{}
	err := db.Connect()
	if err != nil {
		return err
	}
	defer db.Close()
	query := `
	CREATE TABLE IF NOT EXISTS chats (
	    chat_id SERIAL NOT NULL PRIMARY KEY,
	    chat_name VARCHAR(255) NOT NULL,
	    chat_type VARCHAR(255) NOT NULL,
	    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
`
	_, err = db.DB.Exec(query)
	if err != nil {
		return err
	}
	return nil
}
func MessagesCreateTable() error {
	db := DB{}
	err := db.Connect()
	if err != nil {
		return err
	}
	defer db.Close()
	query := `
	CREATE TABLE IF NOT EXISTS messages (
	    message_id SERIAL NOT NULL PRIMARY KEY,
	    chat_id INT REFERENCES chats(chat_id) NOT NULL,
	    sender_id INT REFERENCES users(user_id) NOT NULL,
	    recipient_id INT REFERENCES users(user_id) NOT NULL,
	    message_text text NOT NULL,
	    sent_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
`
	_, err = db.DB.Exec(query)
	if err != nil {
		return err
	}
	return nil
}
func ContactCreateTable() error {
	db := DB{}
	err := db.Connect()
	if err != nil {
		return err
	}
	defer db.Close()
	query := `
	CREATE TABLE IF NOT EXISTS contacts (user_id INT REFERENCES users(user_id) NOT NULL,contact_id INT REFERENCES users(user_id) NOT NULL, 
	    chat_id INT REFERENCES chats(chat_id) NOT NULL );
`
	_, err = db.DB.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func InsertDefoultUser() error {
	db := DB{}
	err := db.Connect()
	if err != nil {
		return err
	}
	defer db.Close()
	var d = -1
	err = db.DB.QueryRow("SELECT user_id FROM users WHERE user_id = 0;").Scan(&d)
	if err != nil {
		return err
	}
	if d == -1 {
		_, err = db.DB.Exec("INSERT INTO users(user_id,username,password) VALUES (0,' ',' ');")
		if err != nil {
			return err
		}
	}
	return nil
}

func Migration() error {
	err := UserCreateTables()
	if err != nil {
		return err
	}
	err = ChatsCreateTable()
	if err != nil {
		return err
	}
	err = MessagesCreateTable()
	if err != nil {
		return err
	}
	err = ContactCreateTable()
	if err != nil {
		return err
	}
	err = InsertDefoultUser()
	if err != nil {
		return err
	}

	return nil

}
