package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"apitest.com/api/controllers"
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

	db := controllers.DB
	// dbsetup end

	// table setups
	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS admins (
            id INT AUTO_INCREMENT NOT NULL,
            username VARCHAR(255) NOT NULL,
            password VARCHAR(255) NOT NULL,
            UNIQUE(id),
            PRIMARY KEY (id)
        );
    `)
	if err != nil {
		fmt.Println("Error creating admins table", err)
	}

	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS blogs (
            id INT AUTO_INCREMENT NOT NULL,
            title VARCHAR(255) NOT NULL,
            description TEXT,
            upload_date DATE NOT NULL,
            update_date DATE,
            link_to_md VARCHAR(255) NOT NULL,
            link_to_jsx VARCHAR(255) NOT NULL,
            UNIQUE(id),
            PRIMARY KEY (id)
        );
    `)
	if err != nil {
		fmt.Println("Error creating blogs table", err)
	}

	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS tags (
            id INT AUTO_INCREMENT NOT NULL,
            title VARCHAR(255) NOT NULL,
            color VARCHAR(255) NOT NULL,
            link_to_svg VARCHAR(255) NOT NULL,
            UNIQUE(id),
            PRIMARY KEY (id)
        );
    `)
	if err != nil {
		fmt.Println("Error creating tags table", err)
	}

	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS blog_tag (
            blog_id INT NOT NULL,
            tag_id INT NOT NULL,
            FOREIGN KEY (blog_id) REFERENCES blogs(id),
            FOREIGN KEY (tag_id) REFERENCES tags(id),
            CONSTRAINT BT_Union PRIMARY KEY (blog_id, tag_id)
        );
    `)
	if err != nil {
		fmt.Println("Error creating blog-tag relational table", err)
	}
    // table setups end
}
