package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
	
)

type Category struct {
	Id   int
	Name string
}


func getCategories(db *sql.DB) []Category {
	rows, err := db.Query("SELECT id, name FROM category")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		var category Category
		err := rows.Scan(&category.Id, &category.Name)
		if err != nil {
			log.Fatal(err)
		}
		categories = append(categories, category)
	}
	return categories
}
