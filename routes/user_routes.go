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
		d.DeleteUserById(w, r)
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

	newUsr, err := u.Create()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        fmt.Println("Error while creating user: ", err)
		return
    }

    w.Header().Set("Content-Type", "application/json")

    jsonUsr, err := json.Marshal(newUsr)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        fmt.Println("Error while marshalling user to json: ", err)
		return
    }
    
    w.Write(jsonUsr)
}

func (d *UsersRoute) GetUsers(w http.ResponseWriter, r *http.Request) {
	userSlice, err := models.ReadAllUsers()
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println("Error while getting user list", err)
		return
    }

    w.Header().Set("Content-Type", "application/json")

    jsonData, err := json.Marshal(userSlice)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        fmt.Println("Error while marshalling user slice to json: ", err)
		return
    }
    
    w.Write(jsonData)	
}

func (d *UsersRoute) GetUsersById(w http.ResponseWriter, r *http.Request) {
	id_path := UserByIdRegex.FindStringSubmatch(r.URL.Path)
	uid, err := strconv.Atoi(id_path[1])
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println("Error while getting id from path", err)
		return
    }

	u, err := models.ReadUserById(uid)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("Error while reading user by id", err)
		return
    }

    w.Header().Set("Content-Type", "application/json")
	
	jsonData, err := json.Marshal(u)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        fmt.Println("Error while marshalling user to json while getting user by id: ", err)
		return
    }
    
    w.Write(jsonData)
}

func (d *UsersRoute) GetUsersByName(w http.ResponseWriter, r *http.Request) {
	id_path := UserByNameRegex.FindStringSubmatch(r.URL.Path)
	username := id_path[1]

	u, err := models.ReadUserByName(username)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println("Error while reading user by name", err)
		return
    }
	
    w.Header().Set("Content-Type", "application/json")

	jsonData, err := json.Marshal(u)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        fmt.Println("Error while marshalling user to json while getting user by name: ", err)
		return
    }
    
    w.Write(jsonData)
}

func (d *UsersRoute) UpdateUserKeyValuePair(w http.ResponseWriter, r *http.Request) {
	id_path := UserByIdRegex.FindStringSubmatch(r.URL.Path)
	uid, err := strconv.Atoi(id_path[1])
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println("Error while getting id from path during update", err)
		return
    }

	u, err := models.ReadUserById(uid)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println("Error while reading user by id during update", err)
		return
    }
	

	type keyRes struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}
	var kvp keyRes

	err = json.NewDecoder(r.Body).Decode(&kvp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("Error while decoding user post from json", err)
		return
	}

	u, err = u.Update(kvp.Key, kvp.Value)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("Error while updating user", err)
		return
    }

    w.Header().Set("Content-Type", "application/json")
	    
    jsonData, err := json.Marshal(u)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        fmt.Println("Error while marshalling user to json after update: ", err)
		return
    }
    
    w.Write(jsonData)
}

func (d *UsersRoute) DeleteUserById(w http.ResponseWriter, r *http.Request) {
	id_path := UserByIdRegex.FindStringSubmatch(r.URL.Path)
	uid, err := strconv.Atoi(id_path[1])
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println("Error while getting id from path during deletion", err)
		return
    }

	u, err := models.ReadUserById(uid)
	if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println("Error while reading user by id during deletion", err)
		return
    }

    err = u.Delete()
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println("Error while deleting user", err)
        return
    }

    w.Write([]byte(fmt.Sprintf("User with id %d was deleted succesfully", uid)))
}

/*
    TODO: MAKE HTMX PAGE FOR USER ADMINISTRATION
*/
