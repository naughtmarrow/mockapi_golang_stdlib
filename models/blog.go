package models

import (
	"errors"
	"fmt"

	"apitest.com/api/controllers"
)

type Blog struct {
	Id          int    `json:"id"`
	Title       string `json:"title" faker:"sentence"`
	Description string `json:"description" faker:"paragraph"`
	Upload_date string `json:"upload_date" faker:"date"`
	Update_date string `json:"update_date" faker:"date"`
	Link_to_md  string `json:"link_to_md" faker:"url"`
	Link_to_jsx string `json:"link_to_jsx" faker:"url"`
}

func (b Blog) Create() (Blog, error) {
	res, err := controllers.DB.Exec(`
        INSERT INTO blogs (title, description, upload_date, update_date, link_to_md, link_to_jsx)
        VALUES (?, ?, ?, ?, ?, ?);
    `, b.Title, b.Description, b.Upload_date, b.Update_date, b.Link_to_md, b.Link_to_jsx)
	if err != nil {
		return b, errors.New(fmt.Sprintf("Error inserting blog with id %d into database with error: %s", b.Id, err))
	}

	insertedId, err := res.LastInsertId()
	return ReadBlogById(int(insertedId))
}

func ReadBlogById(id int) (Blog, error) {
	var blog Blog

	err := controllers.DB.QueryRow(`SELECT id, title, description, upload_date, update_date, link_to_md, link_to_jsx FROM blogs WHERE id = ?;`, id).Scan(&blog.Id, &blog.Title, &blog.Description, &blog.Upload_date, &blog.Update_date, &blog.Link_to_md, &blog.Link_to_jsx)
	if err != nil {
		return blog, errors.New(fmt.Sprintf("Error querying blog with id %d with error: %s", id, err))
	}

	return blog, nil
}

func ReadBlogByName(title string) (Blog, error) {
	var blog Blog

	err := controllers.DB.QueryRow(`SELECT id, title, description, upload_date, update_date, link_to_md, link_to_jsx FROM blogs WHERE title LIKE ?;`, title).Scan(&blog.Id, &blog.Title, &blog.Description, &blog.Upload_date, &blog.Update_date, &blog.Link_to_md, &blog.Link_to_jsx)
	if err != nil {
		return blog, errors.New(fmt.Sprintf("Error querying blog with title %s with error: %s", title, err))
	}

	return blog, nil
}

func ReadAllBlogs() ([]Blog, error) {
	var blogs []Blog
	rows, err := controllers.DB.Query(`SELECT id, title, description, upload_date, update_date, link_to_md, link_to_jsx FROM blogs;`)
	if err != nil {
		return blogs, errors.New(fmt.Sprintf("Error querying all blogs with error: %s", err))
	}
	defer rows.Close()

	for rows.Next() {
		var b Blog
		err = rows.Scan(&b.Id, &b.Title, &b.Description, &b.Upload_date, &b.Update_date, &b.Link_to_md, &b.Link_to_jsx)
		if err != nil {
			return blogs, errors.New(fmt.Sprintf("Error scanning row during read all blogs query with error: %s", err))
		}
		blogs = append(blogs, b)
	}

	err = rows.Err()
	if err != nil {
		return blogs, errors.New(fmt.Sprintf("Error in rows during read all blogs query with error: %s", err))
	}

	return blogs, nil
}

func (blog Blog) Delete() error {
	_, err := controllers.DB.Exec(`DELETE FROM blogs WHERE id = ?;`, blog.Id)
	if err != nil {
		return errors.New(fmt.Sprintf("Error deleting blog with id %d with error: %s", blog.Id, err))
	}
	return nil
}

/*
ENTRIES:
id INT AUTO_INCREMENT NOT NULL,
title VARCHAR(255) NOT NULL,
description TEXT,
upload_date DATE NOT NULL,
update_date DATE,
link_to_md VARCHAR(255) NOT NULL,
link_to_jsx VARCHAR(255) NOT NULL,
*/

// TODO: MAKE PARTIAL UPDATES
