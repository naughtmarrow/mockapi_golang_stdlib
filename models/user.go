package models

import (
	"errors"
	"fmt"

	"apitest.com/api/controllers"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id       int    `json:"id"`
	Username string `json:"username" faker:"username"`
	Password string `json:"password" faker:"password"`
}

func (usr User) Create() (User, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(usr.Password), 14)
	if err != nil {
		return usr, errors.New(fmt.Sprintf("Error hashing password for user %d with error: %s", usr.Id, err))
	}

	hashedpass := string(bytes)

	res, err := controllers.DB.Exec(`
        INSERT INTO admins (username, password)
        VALUES (?, ?);
    `, usr.Username, hashedpass)
	if err != nil {
		return usr, errors.New(fmt.Sprintf("Error inserting user with id %d into database with error: %s", usr.Id, err))
	}

	insertedId, err := res.LastInsertId()
	return ReadUserById(int(insertedId))
}

func ReadAllUsers() ([]User, error) {
	var users []User
	rows, err := controllers.DB.Query(`SELECT id, username, password FROM admins;`)
	if err != nil {
		return users, errors.New(fmt.Sprintf("Error querying all users with error: %s", err))
	}
	defer rows.Close()

	for rows.Next() {
		var u User
		err = rows.Scan(&u.Id, &u.Username, &u.Password)
		if err != nil {
			return users, errors.New(fmt.Sprintf("Error scanning row during read all query with error: %s", err))
		}
		users = append(users, u)
	}

	err = rows.Err()
	if err != nil {
		return users, errors.New(fmt.Sprintf("Error in rows during read all query with error: %s", err))
	}

	return users, nil
}

func ReadUserById(uid int) (User, error) {
	var resUsr User
	err := controllers.DB.QueryRow(`SELECT id, username, password FROM admins WHERE id = ?;`, uid).Scan(&resUsr.Id, &resUsr.Username, &resUsr.Password)
	if err != nil {
		return resUsr, errors.New(fmt.Sprintf("Error querying user with id %d with error: %s", uid, err))
	}

	return resUsr, nil
}

func ReadUserByName(username string) (User, error) {
	var resUsr User

	err := controllers.DB.QueryRow(`SELECT id, username, password FROM admins WHERE username LIKE ?;`, username).Scan(&resUsr.Id, &resUsr.Username, &resUsr.Password)
	if err != nil {
		return resUsr, errors.New(fmt.Sprintf("Error querying user with username %s with error: %s", username, err))
	}

	return resUsr, nil
}

func (usr User) Update(key string, value string) (User, error) {
	if key != "username" && key != "password" {
		fmt.Printf("Error updating user, key is invalid")
		return usr, errors.New("Error updating user, key is invalid")
	}

	if key == "password" {
		bytes, err := bcrypt.GenerateFromPassword([]byte(usr.Password), 14)
		if err != nil {
			return usr, errors.New(fmt.Sprintf("Error hashing password for user with id %d while updating with error: %s", usr.Id, err))
		}
		pass := string(bytes)
		value = pass
	}

	query := "UPDATE admins SET " + key + " = ? WHERE id = ?;"

	res, err := controllers.DB.Exec(query, value, usr.Id)
	if err != nil {
		return usr, errors.New(fmt.Sprintf("Error updating user with id %d with error: %s", usr.Id, err))
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return usr, errors.New(fmt.Sprintf("Error updating user with id %d rows affected not received with error: %s", usr.Id, err))
	} else {
		fmt.Println("Rows affected: ", rows)
	}

	return ReadUserById(usr.Id)
}

func (usr User) Delete() error {
	_, err := controllers.DB.Exec(`DELETE FROM admins WHERE id = ?;`, usr.Id)
	if err != nil {
		return errors.New(fmt.Sprintf("Error deleting user with id %d with error: %s", usr.Id, err))
	}
	return nil
}
