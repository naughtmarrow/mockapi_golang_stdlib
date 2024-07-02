package models

import (
	"errors"
	"fmt"

	"apitest.com/api/controllers"
)

type Tag struct {
	Id          int    `json:"id"`
	Title       string `json:"title" faker:"sentence"`
	Color       string `json:"color" faker:"oneof: red, green, yellow, blue, purple"`
	Link_to_svg string `json:"link_to_svg" faker:"url"`
	Blog_list   []Blog `json:"blog_list"`
}

func (tag Tag) Create() (Tag, error) {
	res, err := controllers.DB.Exec(`
        INSERT INTO tags (title, color, link_to_svg)
        VALUES (?, ?, ?);
    `, tag.Title, tag.Color, tag.Link_to_svg)
	if err != nil {
		return tag, errors.New(fmt.Sprintf("Error inserting tag with id %d into database with error: %s", tag.Id, err))
	}

	insertedId, err := res.LastInsertId()
	return ReadTagById(int(insertedId))
}

func ReadTagById(id int) (Tag, error) {
	var tag Tag
	err := controllers.DB.QueryRow(`SELECT id, title, color, link_to_svg FROM tags WHERE id = ?;`, id).Scan(&tag.Id, &tag.Title, &tag.Color, &tag.Link_to_svg)
	if err != nil {
		return tag, errors.New(fmt.Sprintf("Error querying tag with id %d with error: %s", id, err))
	}

	return tag, nil
}

func ReadTagByName(title string) (Tag, error) {
	var t Tag

	err := controllers.DB.QueryRow(`SELECT id, title, color, link_to_svg FROM tags WHERE title LIKE ?;`, title).Scan(&t.Id, &t.Title, &t.Color, &t.Link_to_svg)
	if err != nil {
		return t, errors.New(fmt.Sprintf("Error querying tag with title %s with error: %s", title, err))
	}

	return t, nil
}

func ReadAllTags() ([]Tag, error) {
	var tags []Tag
	rows, err := controllers.DB.Query(`SELECT id, title, color, link_to_svg FROM tags;`)
	if err != nil {
		return tags, errors.New(fmt.Sprintf("Error querying all tags with error: %s", err))
	}
	defer rows.Close()

	for rows.Next() {
		var t Tag
		err = rows.Scan(&t.Id, &t.Title, &t.Color, &t.Link_to_svg)
		if err != nil {
			return tags, errors.New(fmt.Sprintf("Error scanning row during read all query with error: %s", err))
		}
		tags = append(tags, t)
	}

	err = rows.Err()
	if err != nil {
		return tags, errors.New(fmt.Sprintf("Error in rows during read all tags query with error: %s", err))
	}

	return tags, nil
}

func (tag Tag) Delete() error {
	_, err := controllers.DB.Exec(`DELETE FROM tags WHERE id = ?;`, tag.Id)
	if err != nil {
		return errors.New(fmt.Sprintf("Error deleting tag with id %d with error: %s", tag.Id, err))
	}
	return nil
}

// Methods used to work with the blogs related to a specific tag
func (tag Tag) GetBlogList() ([]Blog, error) {
	var blogs []Blog

	rows, err := controllers.DB.Query(
		`SELECT 
        id, title, description, upload_date, update_date, link_to_md, link_to_jsx 
        FROM blogs
        LEFT JOIN blog_tag
        ON blog_tag.blog_id = blogs.id
        WHERE blog_tag.tag_id = ?;
    `, tag.Id)

	if err != nil {
		return blogs, errors.New(fmt.Sprintf("Error querying blogs with given tag with error: %s", err))
	}
	defer rows.Close()

	for rows.Next() {
		var b Blog
		err = rows.Scan(&b.Id, &b.Title, &b.Description, &b.Upload_date, &b.Update_date, &b.Link_to_md, &b.Link_to_jsx)
		if err != nil {
			return blogs, errors.New(fmt.Sprintf("Error scanning row during read all query in blogs by tag endpoint with error: %s", err))
		}
		blogs = append(blogs, b)
	}

	err = rows.Err()
	if err != nil {
		return blogs, errors.New(fmt.Sprintf("Error in rows during read all tags query with error: %s", err))
	}

	return blogs, nil
}

func (tag Tag) AddBlog(bid int) error {
    _, err := controllers.DB.Exec(`
        INSERT INTO blog_tag (blog_id, tag_id)
        VALUES (?, ?);
    `, bid, tag.Id)
	if err != nil {
		return errors.New(fmt.Sprintf("Error inserting blog with the id %d to the tag with id %d into database with error: %s", bid, tag.Id, err))
	}

    return nil
}

func (tag Tag) DeleteBlog(bid int) error {
    _, err := controllers.DB.Exec(`DELETE FROM blog_tag WHERE blog_id = ? AND tag_id = ?;`, bid, tag.Id)
	if err != nil {
		return errors.New(fmt.Sprintf("Error deleting blog with id %d from tag list on tag with id %d with error: %s", bid, tag.Id, err))
	}
	return nil
}

/*
ENTRIES:
id INT AUTO_INCREMENT NOT NULL,
title VARCHAR(255) NOT NULL,
color VARCHAR(255) NOT NULL,
link_to_svg VARCHAR(255) NOT NULL,
*/

/*
TODO: CREATE PARTIAL UPDATE METHODS
CREATE BLOG LIST METHOD
*/

/*
CREATE TABLE IF NOT EXISTS blog_tag (
    blog_id INT NOT NULL,
    tag_id INT NOT NULL,
    FOREIGN KEY (blog_id) REFERENCES blogs(id),
    FOREIGN KEY (tag_id) REFERENCES tags(id),
    CONSTRAINT BT_Union PRIMARY KEY (blog_id, tag_id)
);
*/
