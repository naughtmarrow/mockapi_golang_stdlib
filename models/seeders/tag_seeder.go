package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"apitest.com/api/models"
    "github.com/go-faker/faker/v4"
    "apitest.com/api/controllers"
	"github.com/joho/godotenv"
)

func __main()  {  
    godotenv.Load(".env")

	dbpath := fmt.Sprintf("%s:%s@(%s:%s)/%s",
		os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOSTNAME"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))

	err := controllers.DBConnect(dbpath)
	if err != nil {
		fmt.Println(err)
	}
	defer controllers.DBClose()

    amount, err := strconv.Atoi(os.Args[1])
    if err != nil {
        log.Fatal(err)
    }

    for i := 0; i < amount; i++{
        t := models.Tag{}

        err := faker.FakeData(&t)
        if err != nil {
            log.Fatal(err)
        }

        _, err = t.Create()
        if err != nil {
            log.Fatal(err)
        }
    }

    fmt.Printf("Database seeded with %d tags\n", amount)
}
