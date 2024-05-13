package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"apitest.com/api/models"
)

var (
	UserBaseRegex   = regexp.MustCompile(`^/users/*$`)
	UserByIdRegex   = regexp.MustCompile(`^/users/([0-9]+)`)
	UserByNameRegex = regexp.MustCompile(`^/users/name/([a-zA-Z0-9_.-]*)$`)
)

type UsersRoute struct{}

// main router function
func (d *UsersRoute) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodPost && UserBaseRegex.MatchString(r.URL.Path):
		d.CreateUser(w, r)
		return
	case r.Method == http.MethodGet && UserBaseRegex.MatchString(r.URL.Path):
		d.GetUsers(w, r)
		return
	case r.Method == http.MethodGet && UserByIdRegex.MatchString(r.URL.Path):
		d.GetUsersById(w, r)
		return
	case r.Method == http.MethodGet && UserByNameRegex.MatchString(r.URL.Path):
		d.GetUsersByName(w, r)
		return
	case r.Method == http.MethodPut && UserByIdRegex.MatchString(r.URL.Path):
		d.UpdateUserKeyValuePair(w, r)
		return
    case r.Method == http.MethodDelete && UserByIdRegex.MatchString(r.URL.Path):
        d.DeleteUserById(w,r)
        return
	default:
		w.Write([]byte("Test response from users"))
		return
	}
}

// routes
func (d *UsersRoute) CreateUser(w http.ResponseWriter, r *http.Request) {
	var u models.User
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println("Error while decoding user post from json", err)
		return
	}
	newUsr := u.Create()
	fmt.Fprintln(w, newUsr)
}

func (d *UsersRoute) GetUsers(w http.ResponseWriter, r *http.Request) {
	userSlice := models.ReadAllUsers()
	for i := range userSlice {
	    fmt.Fprintln(w, userSlice[i])
	}
}

func (d *UsersRoute) GetUsersById(w http.ResponseWriter, r *http.Request) {
	id_path := UserByIdRegex.FindStringSubmatch(r.URL.Path)
	uid, _ := strconv.Atoi(id_path[1])
	u := models.ReadUserById(uid)
	if u.Id == 0 {
		fmt.Println("Error while reading user from id")
		return
	}
	fmt.Fprintln(w, u)
}

func (d *UsersRoute) GetUsersByName(w http.ResponseWriter, r *http.Request) {
	id_path := UserByNameRegex.FindStringSubmatch(r.URL.Path)
	username := id_path[1]
	u := models.ReadUserByName(username)
	if u.Id == 0 {
		fmt.Println("Error while reading user from id")
		return
	}
	fmt.Fprintln(w, u)
}

func (d *UsersRoute) UpdateUserKeyValuePair(w http.ResponseWriter, r *http.Request) {
    id_path := UserByIdRegex.FindStringSubmatch(r.URL.Path)
    uid, _ := strconv.Atoi(id_path[1])
    u := models.ReadUserById(uid)
    if u.Id == 0 {
		fmt.Println("Error while reading user from id in update method")
		return
	}
    
    type keyRes struct {
        Key string `json:"key"`
        Value string `json:"value"`
    }
    var kvp keyRes

    err := json.NewDecoder(r.Body).Decode(&kvp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println("Error while decoding user post from json", err)
		return
	}

    u.Update(kvp.Key, kvp.Value)
    fmt.Fprintln(w, "Before: ")
    fmt.Fprintln(w, u)
    fmt.Fprintln(w, "_________________________")
    fmt.Fprintln(w, "After: ")
    fmt.Fprintln(w, models.ReadUserById(uid))
}

func (d *UsersRoute) DeleteUserById(w http.ResponseWriter, r *http.Request) {
    id_path := UserByIdRegex.FindStringSubmatch(r.URL.Path)
    uid, _ := strconv.Atoi(id_path[1])
    u := models.ReadUserById(uid)
    if u.Id == 0 {
		fmt.Println("Error while reading user from id in update method")
		return
	}

    u.Delete()
    fmt.Fprintln(w, "User has been deleted")
}

/*
    TODO : REFACTOR USER METHODS TO INCLUDE ERROR AS RETURN TYPE FOR BETTER ERROR HANDLING
    TODO : REFACTOR RESPONSES FROM API TO SEND JSON AS RESPONSE TO BE INTERPRETED BY FRONTEND INSTEAD OF WRITING TO OUTPUT STREAM
*/
