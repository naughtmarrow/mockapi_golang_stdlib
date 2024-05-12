package routes

import "net/http"

type DefaultRoute struct{}

func (d *DefaultRoute) ServeHTTP(w http.ResponseWriter, r *http.Request){
    w.Write([]byte("Test response"))
}
