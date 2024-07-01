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

/*
ENTRIES:
id INT AUTO_INCREMENT NOT NULL,
title VARCHAR(255) NOT NULL,
color VARCHAR(255) NOT NULL,
link_to_svg VARCHAR(255) NOT NULL,
*/

/*
TODO: CREATE PARTIAL UPDATE METHODS
*/
