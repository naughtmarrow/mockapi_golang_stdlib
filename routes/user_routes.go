package routes

import "net/http"

type UsersRoute struct{}

func (d *UsersRoute) ServeHTTP(w http.ResponseWriter, r *http.Request){
    w.Write([]byte("Test response"))
}
