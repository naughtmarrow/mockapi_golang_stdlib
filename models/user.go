package models

import (
	"fmt"

	"apitest.com/api/controllers"
	"golang.org/x/crypto/bcrypt"
)

type User struct{
    Id int `json:"id"`
    Username string `json:"username"`
    Password string `json:"password"`
}

func (usr User) Create() User {
    bytes, err := bcrypt.GenerateFromPassword([]byte(usr.Password), 14) 
    if err != nil {
        fmt.Printf("Error hashing password for user %d\n", usr.Id)
        fmt.Println(err)
    }
    hashedpass := string(bytes)
    res, err := controllers.DB.Exec(`
        INSERT INTO admins (username, password)
        VALUES (?, ?);
    `, usr.Username, hashedpass)
    if err != nil {
        fmt.Printf("Error inserting user with id %d into database", usr.Id)
        fmt.Println(err)
    }
    
    insertedId, err := res.LastInsertId()
    return ReadUserById(int(insertedId))
}

func ReadAllUsers() []User {
    var users []User
    rows, err := controllers.DB.Query(`SELECT id, username, password FROM admins;`)
    if err != nil {
        fmt.Println("Error querying all users", err)
    }
    defer rows.Close()

    for rows.Next(){
        var u User
        err = rows.Scan(&u.Id, &u.Username, &u.Password) 
        if err != nil {
            fmt.Println("Error scanning row during read all query", err)
        }
        users = append(users, u)
    }

    err = rows.Err()
    if err != nil {
        fmt.Println("Error in rows during read all query", err)
    }

    return users
}

func ReadUserById(uid int) User {
    var resUsr User
    err := controllers.DB.QueryRow(`SELECT id, username, password FROM admins WHERE id = ?;`, uid).Scan(&resUsr.Id, &resUsr.Username, &resUsr.Password) 
    if err != nil {
        fmt.Printf("Error querying user with id %d\n", uid)
        fmt.Println(err)
    }

    return resUsr
}

func ReadUserByName(username string) User {
    var resUsr User

    err := controllers.DB.QueryRow(`SELECT id, username, password FROM admins WHERE username LIKE ?;`, username).Scan(&resUsr.Id, &resUsr.Username, &resUsr.Password)
    if err != nil {
        fmt.Printf("Error querying user with username %s\n", username)
        fmt.Println(err)
    }
     
    return resUsr
}

func (usr User) Update(key string, value string){
    if key != "username" && key != "password" {
        fmt.Printf("Error updating user, key is invalid")
        return
    }
    if key == "password" {
        bytes, err := bcrypt.GenerateFromPassword([]byte(usr.Password), 14)
        if err != nil {
            fmt.Printf("Error hashing password for user %d while updating\n", usr.Id)
            fmt.Println(err)
        }
        pass := string(bytes)   
        value = pass
    }

    query := "UPDATE admins SET " + key + " = ? WHERE id = ?;"

    res, err := controllers.DB.Exec(query, value, usr.Id)
    if err != nil {
        fmt.Printf("Error updating user with id %d\n", usr.Id)
        fmt.Println(err)
    }
    
    rows, err := res.RowsAffected()
    if err != nil {
        fmt.Printf("Error updating user with id %d rows affected not received\n", usr.Id)
        fmt.Println(err)
    } else {
        fmt.Println("Rows affected: ", rows)
    }
}

func (usr User) Delete(){
    _, err := controllers.DB.Exec(`DELETE FROM admins WHERE id = ?;`, usr.Id)
    if err != nil {
        fmt.Printf("Error deleting user with id %d", usr.Id)
        fmt.Println(err)
    }
}
