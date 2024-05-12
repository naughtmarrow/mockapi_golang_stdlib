package controllers

import (
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func DBConnect(dbpath string) error {
    var err error
    DB, err = sql.Open("mysql", dbpath)
    if err != nil {
        return err
    }
    return nil
}

func DBClose() error {
    return DB.Close()
}
