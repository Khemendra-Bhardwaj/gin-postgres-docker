package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	connStr := "user=postgres password=postgres123 dbname=mydatabase host=postgres sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS people (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100),
		age INT
	)`)

	if err != nil {
		log.Fatal(err)
	}

	router := gin.Default()

	router.GET("/people", func(c *gin.Context) {
		rows, err := db.Query("SELECT id, name, age FROM people")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var people []map[string]interface{}
		for rows.Next() {
			var id int
			var name string
			var age int
			if err := rows.Scan(&id, &name, &age); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			people = append(people, map[string]interface{}{
				"id":   id,
				"name": name,
				"age":  age,
			})
		}
		c.JSON(http.StatusOK, people)
	})

	router.POST("/people", func(c *gin.Context) {
		var person struct {
			Name string `json:"name"`
			Age  int    `json:"age"`
		}
		if err := c.ShouldBindJSON(&person); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		_, err := db.Exec("INSERT INTO people (name, age) VALUES ($1, $2)", person.Name, person.Age)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "success"})

	})

	router.Run(":8080")

}
