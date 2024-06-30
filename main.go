package main

import (
	"fmt"
	"net/http"
	"os"

	"apitest.com/api/controllers"
	"apitest.com/api/routes"
	"github.com/joho/godotenv"
)

func main() {
    // dbsetup
	godotenv.Load(".env")

	dbpath := fmt.Sprintf("%s:%s@(%s:%s)/%s",
		os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOSTNAME"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))

	err := controllers.DBConnect(dbpath)
	if err != nil {
		fmt.Println(err)
	}
	defer controllers.DBClose()
    // dbsetup end

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("./static")))
	mux.Handle("/users", &routes.UsersRoute{})
    mux.Handle("/users/", &routes.UsersRoute{})
    mux.Handle("/admin", http.StripPrefix("/admin", http.FileServer(http.Dir("./static"))))
    mux.Handle("/admin/", http.StripPrefix("/admin/", http.FileServer(http.Dir("./static"))))

    err = http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("SERVER_PORT")), mux)
    fmt.Printf("Error: %s", err);
}
