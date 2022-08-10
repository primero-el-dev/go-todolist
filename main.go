package main

import (
	"database/sql"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
	_ "github.com/lib/pq"
	"log"
	"os"
)

const (
	host     = "localhost"
	user     = "appadmin"
	password = "password"
	dbname   = "go_todo_list"
)

func main() {
	// Connect to database
	connStr := fmt.Sprintf("user=%s password=%s host=%s dbname=%s sslmode=disable", user, password, host, dbname)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	// Init app
	engine := html.New("./views", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})
	app.Static("/", "./public")

	app.Get("/", func(c *fiber.Ctx) error {
		return indexHandler(c, db)
	})

	app.Post("/", func(c *fiber.Ctx) error {
		return postHandler(c, db)
	})

	app.Post("/:id/edit", func(c *fiber.Ctx) error {
		return putHandler(c, db)
	})

	app.Post("/:id/delete", func(c *fiber.Ctx) error {
		return deleteHandler(c, db)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	log.Fatalln(app.Listen(fmt.Sprintf(":%v", port)))
}

func indexHandler(c *fiber.Ctx, db *sql.DB) error {
	var task Task
	var tasks []Task

	rows, err := db.Query("SELECT * FROM task")
	defer rows.Close()
	if err != nil {
		log.Fatalln(err)
		c.JSON("An error occured")
	}

	for rows.Next() {
		rows.Scan(&task.Id, &task.Description, &task.Status, &task.CreatedAt)
		tasks = append(tasks, task)
	}

	return c.Render("index", fiber.Map{
		"Tasks": tasks,
	})
}

func postHandler(c *fiber.Ctx, db *sql.DB) error {
	task := Task{}
	if err := c.BodyParser(&task); err != nil {
		log.Printf("An error occured: %v", err)
		return c.SendString(err.Error())
	}

	if task.IsValid() {
		_, err := db.Exec("INSERT into task (Description, Status) VALUES ($1, $2)", task.Description, task.Status)
		if err != nil {
			log.Fatalf("An error occured while executing query: %v", err)
		}
	}

	return c.Redirect("/")
}

func putHandler(c *fiber.Ctx, db *sql.DB) error {
	id := c.Params("id")
	task := Task{}

	if err := c.BodyParser(&task); err != nil {
		log.Printf("An error occured: %v", err)

		return c.SendString(err.Error())
	}

	fmt.Println(task.Status)
	if task.IsValid() {
		db.Exec("UPDATE task SET description = $1, status = $2 WHERE id = $3", task.Description, task.Status, id)
	}

	return c.Redirect("/")
}

func deleteHandler(c *fiber.Ctx, db *sql.DB) error {
	id := c.Params("id")
	db.Exec("DELETE from task WHERE id = $1", id)

	return c.Redirect("/")
}
